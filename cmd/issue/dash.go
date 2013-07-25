// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Issue dashboard

package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// go12 web page update

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
	GoXX      int
	GoXXMaybe int
}

func (p Point) JDate() template.JS {
	yy, mm, dd := p.Time.Date()
	h, m, s := p.Time.Clock()
	return template.JS(fmt.Sprintf("new Date(%d, %d, %d, %d, %d, %d)", yy, mm-1, dd, h, m, s))
}

func dash(version string) {
	prefix := strings.Map(func(r rune) rune {
		if r == '.' {
			return -1
		}
		if 'A' <= r && r <= 'Z' {
			return r - 'A' + 'a'
		}
		return r
	}, version)

	data, err := ioutil.ReadFile(prefix + ".graph")
	if err != nil {
		log.Fatal(err) // in wrong directory
	}
	var graph []Point
	if err := json.Unmarshal(data, &graph); err != nil {
		log.Fatal(err)
	}

	var data1, data2 []byte
	if *quick {
		data1, err = ioutil.ReadFile(prefix + ".xml")
		if err != nil {
			log.Fatal(err)
		}
		data2, err = ioutil.ReadFile(prefix + "maybe.xml")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		r, err := http.Get("https://code.google.com/feeds/issues/p/go/issues/full?can=open&q=label:" + version + "&max-results=1000")
		if err != nil {
			log.Fatal(err)
		}
		data1, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
		r, err = http.Get("https://code.google.com/feeds/issues/p/go/issues/full?can=open&q=label:" + version + "Maybe&max-results=1000")
		if err != nil {
			log.Fatal(err)
		}
		data2, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
		ioutil.WriteFile(prefix+".xml", data1, 0666)
		ioutil.WriteFile(prefix+"maybe.xml", data2, 0666)
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
		if err := ioutil.WriteFile(prefix+".graph", buf, 0666); err != nil {
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
		if p.GoXXMaybe == 0 {
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
			Parse(dashTemplate),
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
	if err := ioutil.WriteFile(prefix+".html", buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}

func labeledGroup(g []Entry, s string) string {
	if strings.HasPrefix(s, "-") {
		for _, e := range g {
			if t := e.Labeled(s[1:]); t != "" {
				return ""
			}
		}
		return "ok"
	}

	for _, e := range g {
		if t := e.Labeled(s); t != "" {
			return t
		}
	}
	return ""
}

var dashTemplate = `<html>
  <head>
    <script type="text/javascript" src="https://www.google.com/jsapi"></script>
    <script type="text/javascript">
      google.load("visualization", "1", {packages:["corechart"]});
      google.setOnLoadCallback(drawCharts);
      function drawCharts() {
        var data = new google.visualization.DataTable();
        data.addColumn('datetime', 'Date');
        data.addColumn('number', '{{.Version}}');
        data.addColumn('number', '{{.Version}} + Maybe');
        var one = 1;
        data.addRows([
{{range .Graph}}          [{{.JDate}}, {{.GoXX}}, {{.GoXX}}+{{.GoXXMaybe}}],
{{end}}
        ])
        var options = {
          width: 800, height: 400,
          title: '{{.Version}} Issues',
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
          $("#govv").html("{{.Version}} issues only");
        } else {
          $("#govv").html("<a href='javascript:dogovv()'>show {{.Version}} issues only</a>");
        }
        if(!onlySuggest && !hideMaybe) {
          $("#govvmaybe").html("all issues");
        } else {
          $("#govvmaybe").html("<a href='javascript:dogovvmaybe()'>show all issues</a>");
        }
      }
      function dosuggest() {
        onlySuggest = true;
        hideMaybe = false;
        window.location.hash = "s";
        rehide();
      }
      function dogovv() {
        onlySuggest = false;
        hideMaybe = true;
        window.location.hash = "";
        rehide();
      }
      function dogovvmaybe() {
        onlySuggest = false;
        hideMaybe = false;
        window.location.hash = "m";
        rehide();
      }
      function start() {
        if (window.location.hash == "s" || window.location.hash == "#s") {
          dosuggest();
        } else if (window.location.hash == "m" || window.location.hash == "#m") {
          dogovvmaybe();
        } else {
          dogovv();
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
    <h1>{{.Version}}: Open Issues</h1>

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
    
    <span id="suggest"></span> | <span id="govv"></span> | <span id="govvmaybe"></span>

    <br><br>
    <table>
    {{range .Issue}}
      <tr class="{{if labeledGroup . "{{.Version}}"}}nomaybe{{else}}maybex{{end}} {{if labeledGroup . "Suggested"}}suggest{{else}}nosuggest{{end}}"><td class="dir" colspan="4">{{(index . 0).Dir}}
      {{range .}}
        <tr class="{{if .Labeled "{{.Version}}Maybe"}}maybe{{else}}nomaybe{{end}} {{if .Labeled "Suggested"}}suggest{{else}}nosuggest{{end}}">
          <td class="suggest">{{if .Labeled "Suggested"}}&#x261e;{{end}}
          <td class="size">{{.Labeled "Size-"}}
          <td class="num">{{.Number}}
          <td class="title"><a href="http://golang.org/issue/{{.Number}}">{{.Title}}</a>
            {{if .Labeled "{{.Version}}Maybe"}}[maybe]{{end}}
            {{if .IsStatus "Started"}}[<i>started by {{.Owner}}</i>]{{end}}
      {{end}}
    {{end}}
    </table>
  </body>
</html>
`
