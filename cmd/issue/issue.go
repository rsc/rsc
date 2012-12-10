package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: issue [-p project] query

If query is a single number, prints the full history for the issue.
Otherwise, prints a table of matching results.
The special query 'go1' is shorthand for 'Priority-Go1'.
`)
	os.Exit(2)
}

type Feed struct {
	Entry Entries `xml:"entry"`
}

type Entry struct {
	ID        string    `xml:"id"`
	Title     string    `xml:"title"`
	Published time.Time `xml:"published"`
	Content   string    `xml:"content"`
	Updates   []Update  `xml:"updates"`
	Author    struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Owner  string   `xml:"owner>username"`
	Status string   `xml:"status"`
	Label  []string `xml:"label"`

	Dir    string
	Number int
}

func (e Entry) IsStatus(s string) bool { return e.Status == s }

func (e Entry) Labeled(name string) string {
	for _, l := range e.Label {
		if l == name {
			return name
		}
		if strings.HasSuffix(name, "-") && strings.HasPrefix(l, name) {
			return l[len(name):]
		}
	}
	return ""
}

type Update struct {
	Summary string `xml:"summary"`
	Owner   string `xml:"ownerUpdate"`
	Label   string `xml:"label"`
	Status  string `xml:"status"`
}

type Entries []Entry

func (e Entries) Len() int           { return len(e) }
func (e Entries) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e Entries) Less(i, j int) bool { return e[i].Title < e[j].Title }

var project = flag.String("p", "go", "code.google.com project identifier")
var v = flag.Bool("v", false, "verbose")
var go11flag = flag.Bool("go11", false, "go11 web page update")
var quick = flag.Bool("quick", false, "use cached xml from disk")

func main() {
	flag.Usage = usage
	flag.Parse()

	if *go11flag {
		go11()
		return
	}

	if flag.NArg() != 1 {
		usage()
	}

	full := false
	q := flag.Arg(0)
	n, _ := strconv.Atoi(q)
	if n != 0 {
		q = "id:" + q
		full = true
	}
	if q == "go1" {
		q = "label:Priority-Go1"
	}

	log.SetFlags(0)

	query := url.Values{
		"q":           {q},
		"max-results": {"600"},
	}
	if !full {
		query["can"] = []string{"open"}
	}
	u := "https://code.google.com/feeds/issues/p/" + *project + "/issues/full?" + query.Encode()
	if *v {
		log.Print(u)
	}
	r, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}

	var feed Feed
	if err := xml.NewDecoder(r.Body).Decode(&feed); err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	sort.Sort(feed.Entry)
	for _, e := range feed.Entry {
		id := e.ID
		if i := strings.Index(id, "id="); i >= 0 {
			id = id[:i+len("id=")]
		}
		fmt.Printf("%s\t%s\n", id, e.Title)
		if full {
			fmt.Printf("Reported by %s (%s)\n", e.Author.Name, e.Published.Format("2006-01-02 15:04:05"))
			if e.Owner != "" {
				fmt.Printf("\tOwner: %s\n", e.Owner)
			}
			if e.Status != "" {
				fmt.Printf("\tStatus: %s\n", e.Status)
			}
			for _, l := range e.Label {
				fmt.Printf("\tLabel: %s\n", l)
			}
			if e.Content != "" {
				fmt.Printf("\n\t%s\n", wrap(html.UnescapeString(e.Content), "\t"))
			}
			u := "https://code.google.com/feeds/issues/p/" + *project + "/issues/" + id + "/comments/full"
			if *v {
				log.Print(u)
			}
			r, err := http.Get(u)
			if err != nil {
				log.Fatal(err)
			}

			var feed Feed
			if err := xml.NewDecoder(r.Body).Decode(&feed); err != nil {
				log.Fatal(err)
			}
			r.Body.Close()

			for _, e := range feed.Entry {
				fmt.Printf("\n%s (%s)\n", e.Title, e.Published.Format("2006-01-02 15:04:05"))
				for _, up := range e.Updates {
					if up.Summary != "" {
						fmt.Printf("\tSummary: %s\n", up.Summary)
					}
					if up.Owner != "" {
						fmt.Printf("\tOwner: %s\n", up.Owner)
					}
					if up.Status != "" {
						fmt.Printf("\tStatus: %s\n", up.Status)
					}
					if up.Label != "" {
						fmt.Printf("\tLabel: %s\n", up.Label)
					}
				}
				if e.Content != "" {
					fmt.Printf("\n\t%s\n", wrap(html.UnescapeString(e.Content), "\t"))
				}
			}
		}
	}
}

func wrap(t string, prefix string) string {
	out := ""
	t = strings.Replace(t, "\r\n", "\n", -1)
	lines := strings.Split(t, "\n")
	for i, line := range lines {
		if i > 0 {
			out += "\n" + prefix
		}
		s := line
		for len(s) > 70 {
			i := strings.LastIndex(s[:70], " ")
			if i < 0 {
				i = 69
			}
			i++
			out += s[:i] + "\n" + prefix
			s = s[i:]
		}
		out += s
	}
	return out
}

// go11 web page update

type ByDir []Entry

func (x ByDir) Len() int      { return len(x) }
func (x ByDir) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x ByDir) Less(i, j int) bool {
	if x[i].Dir != x[j].Dir {
		return x[i].Dir < x[j].Dir
	}
	return x[i].Number < x[j].Number
}

type Point struct {
	Time      time.Time
	Go11      int
	Go11Maybe int
}

func (p Point) JDate() template.JS {
	yy, mm, dd := p.Time.Date()
	h, m, s := p.Time.Clock()
	return template.JS(fmt.Sprintf("new Date(%d, %d, %d, %d, %d, %d)", yy, mm-1, dd, h, m, s))
}

func go11() {
	data, err := ioutil.ReadFile("go11.graph")
	if err != nil {
		log.Fatal(err) // in wrong directory
	}
	var graph []Point
	if err := json.Unmarshal(data, &graph); err != nil {
		log.Fatal(err)
	}

	var data1, data2 []byte
	if *quick {
		data1, err = ioutil.ReadFile("go11.xml")
		if err != nil {
			log.Fatal(err)
		}
		data2, err = ioutil.ReadFile("go11maybe.xml")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		r, err := http.Get("https://code.google.com/feeds/issues/p/go/issues/full?can=open&q=label:Go1.1&max-results=600")
		if err != nil {
			log.Fatal(err)
		}
		data1, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
		r, err = http.Get("https://code.google.com/feeds/issues/p/go/issues/full?can=open&q=label:Go1.1Maybe&max-results=600")
		if err != nil {
			log.Fatal(err)
		}
		data2, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
		ioutil.WriteFile("go11.xml", data1, 0666)
		ioutil.WriteFile("go11maybe.xml", data2, 0666)
	}

	var feed Feed
	if err := xml.Unmarshal(data1, &feed); err != nil {
		log.Fatal(err)
	}
	n1 := len(feed.Entry)
	if err := xml.Unmarshal(data2, &feed); err != nil {
		log.Fatal(err)
	}
	n2 := len(feed.Entry) - n1

	if !*quick {
		graph = append(graph, Point{time.Now(), n1, n2})
		buf, err := json.Marshal(graph)
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile("go11.graph", buf, 0666); err != nil {
			log.Fatal(err)
		}
	}

	for i := range feed.Entry {
		e := &feed.Entry[i]
		e.Number, _ = strconv.Atoi(e.ID)
		dir := e.Title
		if i := strings.Index(dir, ":"); i >= 0 {
			dir = dir[:i]
		}
		if i := strings.Index(dir, ","); i >= 0 {
			dir = dir[:i]
		}
		e.Dir = dir
	}

	sort.Sort(ByDir(feed.Entry))

	var groups [][]Entry
	dir := ""
	for _, e := range feed.Entry {
		if e.Dir != dir {
			dir = e.Dir
			groups = append(groups, nil)
		}
		n := len(groups) - 1
		groups[n] = append(groups[n], e)
	}

	var small []Point
	now := time.Now()
	day := -1
	for _, p := range graph {
		if p.Go11Maybe == 0 {
			continue
		}
		d := p.Time.Day()
		if d != day || now.Sub(p.Time) < 3*24*time.Hour {
			day = d
			small = append(small, p)
		}
	}

	tmpl := template.Must(
		template.New("main").
			Funcs(template.FuncMap{"labeledGroup": labeledGroup}).
			Parse(go11template),
	)
	var all = struct {
		Graph []Point
		Issue [][]Entry
	}{
		Graph: small,
		Issue: groups,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, &all); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("go11.html", buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}

func labeledGroup(g []Entry, s string) string {
	for _, e := range g {
		if t := e.Labeled(s); t != "" {
			return t
		}
	}
	return ""
}

var go11template = `<html>
  <head>
    <script type="text/javascript" src="https://www.google.com/jsapi"></script>
    <script type="text/javascript">
      google.load("visualization", "1", {packages:["corechart"]});
      google.setOnLoadCallback(drawCharts);
      function drawCharts() {
        var data = new google.visualization.DataTable();
        data.addColumn('datetime', 'Date');
        data.addColumn('number', 'Go 1.1');
        data.addColumn('number', 'Go 1.1 + Maybe');
        var one = 1;
        data.addRows([
{{range .Graph}}          [{{.JDate}}, {{.Go11}}, {{.Go11}}+{{.Go11Maybe}}],
{{end}}
        ])
        var options = {
          width: 800, height: 400,
          title: 'Go 1.1 Issues',
          strictFirstColumnType: true,
          vAxis: {minValue: 0, maxValue: 299},
          vAxes: {0: {title: 'Open Issues'}}
        };
        var chart = new google.visualization.AreaChart(document.getElementById('open_div'));
        chart.draw(data, options);
      }
    </script>
    <script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/1.8.2/jquery.min.js"></script>
    <script>
      var onlySuggest = false;
      var hideMaybe = true;
      function rehide() {
        $("tr").show();
        if(onlySuggest) {
          $("tr.nosuggest").hide();
        }
        if(hideMaybe) {
          $("tr.maybe").hide();
          $("tr.maybex").hide();
        }
        
        if(onlySuggest && !hideMaybe) {
          $("#suggest").html("suggested issues only");
        } else {
          $("#suggest").html("<a href='javascript:dosuggest()'>show suggested issues only</a>");
        }
        if(!onlySuggest && hideMaybe) {
          $("#go11").html("Go1.1 issues only");
        } else {
          $("#go11").html("<a href='javascript:dogo11()'>show Go 1.1 issues only</a>");
        }
        if(!onlySuggest && !hideMaybe) {
          $("#go11maybe").html("all issues");
        } else {
          $("#go11maybe").html("<a href='javascript:dogo11maybe()'>show all issues</a>");
        }
      }
      function dosuggest() {
        onlySuggest = true;
        hideMaybe = false;
        window.location.hash = "s";
        rehide();
      }
      function dogo11() {
        onlySuggest = false;
        hideMaybe = true;
        window.location.hash = "";
        rehide();
      }
      function dogo11maybe() {
        onlySuggest = false;
        hideMaybe = false;
        window.location.hash = "m";
        rehide();
      }
      function start() {
        if (window.location.hash == "s" || window.location.hash == "#s") {
          dosuggest();
        } else if (window.location.hash == "m" || window.location.hash == "#m") {
          dogo11maybe();
        } else {
          dogo11();
        }
      }
    </script>
    
    <style>
      td.dir {font-weight: bold;}
      td.suggest {padding-left: 1em;}
      .size {font-family: sans-serif; font-size: 70%; text-align: center;}
      tr.maybe {color: #aaa;}
      tr.suggest {}
      h1 {font-size: 120%;}
      a {color: #000;}
      tr.maybe a {color: #aaa;}
      .key, .key td {font-family: sans-serif; font-size: 90%;}
    </style>
  </head>

  <body onload="start()">
    <h1>Go 1.1: Open Issues</h1>

    <div id="open_div"></div>
    
    <div class="key">
    Key:
    <table>
    <tr><td class="suggest"><td class="size">S</td><td>small change: less than 30 minutes (e.g. doc fix)
    <tr><td class="suggest"><td class="size">M</td><td>medium change: less than 2 hours (e.g. small feature/fix + tests)
    <tr><td class="suggest"><td class="size">L</td><td>large change: less than 8 hours
    <tr><td class="suggest"><td class="size">XL</td><td>extra large change: more than one day
    <tr><td class="suggest"><td>&#x261e;</td><td>suggested for people looking for work
    </table>
    </div>
    <br><br>
    
    <span id="suggest"></span> | <span id="go11"></span> | <span id="go11maybe"></span>

    <br><br>
    <table>
    {{range .Issue}}
      <tr class="{{if labeledGroup . "Go1.1Maybe"}}maybex{{else}}nomaybe{{end}} {{if labeledGroup . "Suggested"}}suggest{{else}}nosuggest{{end}}"><td class="dir" colspan="4">{{(index . 0).Dir}}
      {{range .}}
        <tr class="{{if .Labeled "Go1.1Maybe"}}maybe{{else}}nomaybe{{end}} {{if .Labeled "Suggested"}}suggest{{else}}nosuggest{{end}}">
          <td class="suggest">{{if .Labeled "Suggested"}}&#x261e;{{end}}
          <td class="size">{{.Labeled "Size-"}}
          <td class="num">{{.Number}}
          <td class="title"><a href="http://golang.org/issue/{{.Number}}">{{.Title}}</a>
            {{if .Labeled "Go1.1Maybe"}}[maybe]{{end}}
            {{if .IsStatus "Started"}}[<i>started by {{.Owner}}</i>]{{end}}
      {{end}}
    {{end}}
    </table>
  </body>
</html>
`
