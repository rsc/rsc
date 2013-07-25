// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/rsc/oauthprompt"
)

var auth struct {
	APIClientID     string
	APIClientSecret string
}

const Version = "Go1.2"

var aflag = flag.Bool("a", false, "run in acme mode")
var project = flag.String("p", "go", "code.google.com project identifier")
var v = flag.Bool("v", false, "verbose")
var version = flag.String("dash", "", "Label for dashboard, like \"Go1.2\"")
var quick = flag.Bool("quick", false, "use cached xml from disk")
var xmlflag = flag.Bool("xml", false, "dump xml")

func usage() {
	fmt.Fprintf(os.Stderr, `usage: issue [-a] [-dash Go1.2] [-p project] [query]

If query is a single number, prints the full history for the issue.
Otherwise, prints a table of matching results.
The special query 'go1' is shorthand for 'Priority-Go1'.

The -a flag runs as an Acme window, making the query optional.

The -dash mode generates an issue dashboard for a Go release.
It maintains a collection of output files named for the release:
go12.html, go12.graph, and so on.
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
	Owner      string   `xml:"owner>username"`
	Status     string   `xml:"status"`
	Label      []string `xml:"label"`
	MergedInto string   `xml:"mergedInto"`
	CC         []string `xml:"cc>username"`

	Dir      string
	Number   int
	Comments []Entry
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

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	if *version != "" {
		dash(*version)
		return
	}

	if *aflag {
		if err := login(); err != nil {
			log.Fatal(err)
		}
		acmeMode()
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

	data, err := fetch(q, full)
	if err != nil {
		log.Fatal(err)
	}

	if full {
		printFull(os.Stdout, data)
	} else {
		printList(os.Stdout, data)
	}
}

type Change struct {
	Summary string
	Status  string
	Owner   string
	Label   []string
	CC      []string
	Comment string
}

func write(id int, ch *Change) error {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version='1.0' encoding='UTF-8'?>
<entry xmlns='http://www.w3.org/2005/Atom' xmlns:issues='http://schemas.google.com/projecthosting/issues/2009'>
  <content type='html'>`)
	xml.Escape(&buf, []byte(ch.Comment))
	buf.WriteString(`</content>
  <author>
    <name>ignored</name>
  </author>
  <issues:updates>
`)
	tag := func(t, data string) {
		buf.WriteString(`    ` + t)
		xml.Escape(&buf, []byte(data))
		buf.WriteString(`</` + t[1:])
	}

	if ch.Summary != "" {
		tag("<issues:summary>", ch.Summary)
	}
	if ch.Status != "" {
		status := ch.Status
		merge := ""
		if strings.HasPrefix(status, "Duplicate ") {
			merge = strings.TrimPrefix(status, "Duplicate ")
			status = "Duplicate"
		}
		tag("<issues:status>", status)
		if merge != "" {
			tag("<issues:mergedInto>", merge)
		}
	}
	if ch.Owner != "" {
		tag("<issues:ownerUpdate>", ch.Owner)
	}
	for _, l := range ch.Label {
		tag("<issues:label>", l)
	}
	for _, cc := range ch.CC {
		tag("<issues:ccUpdate>", cc)
	}
	buf.WriteString(`
  </issues:updates>
</entry>
`)

	// Done with XML!

	u := "https://code.google.com/feeds/issues/p/" + *project + "/issues/" + fmt.Sprint(id) + "/comments/full"
	req, err := http.NewRequest("POST", u, &buf)
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}
	req.Header.Set("Content-Type", "application/atom+xml")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		buf.Reset()
		io.Copy(&buf, resp.Body)
		return fmt.Errorf("write: %v\n%s", resp.Status, buf.String())
	}
	return nil
}

func fetch(q string, full bool) ([]Entry, error) {
	query := url.Values{
		"q":           {q},
		"max-results": {"1000"},
	}
	if !full {
		query["can"] = []string{"open"}
	}
	u := "https://code.google.com/feeds/issues/p/" + *project + "/issues/full?" + query.Encode()
	if *v {
		log.Print(u)
	}
	r, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if *xmlflag {
		io.Copy(os.Stdout, r.Body)
		os.Exit(0)
	}

	var feed Feed
	if err := xml.NewDecoder(r.Body).Decode(&feed); err != nil {
		return nil, err
	}

	sort.Sort(feed.Entry)
	if full {
		for i := range feed.Entry {
			e := &feed.Entry[i]
			id := e.ID
			if i := strings.Index(id, "id="); i >= 0 {
				id = id[:i+len("id=")]
			}
			u := "https://code.google.com/feeds/issues/p/" + *project + "/issues/" + id + "/comments/full"
			if *v {
				log.Print(u)
			}
			r, err := http.Get(u)
			if err != nil {
				return nil, err
			}

			var feed Feed
			if err := xml.NewDecoder(r.Body).Decode(&feed); err != nil {
				return nil, err
			}
			r.Body.Close()
			e.Comments = feed.Entry
		}
	}

	return feed.Entry, nil
}

func printFull(w io.Writer, entries []Entry) {
	for _, e := range entries {
		id := e.ID
		if i := strings.Index(id, "id="); i >= 0 {
			id = id[:i+len("id=")]
		}
		fmt.Fprintf(w, "Summary: %s\n", e.Title)
		fmt.Fprintf(w, "Status: %s", e.Status)
		if e.Status == "Duplicate" {
			fmt.Fprintf(w, " %s", e.MergedInto)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "Owner: %s\n", e.Owner)
		fmt.Fprintf(w, "CC:")
		for _, cc := range e.CC {
			fmt.Fprintf(w, " %s", cc)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "Labels:")
		for _, l := range e.Label {
			fmt.Fprintf(w, " %s", l)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\n")

		fmt.Fprintf(w, "Reported by %s (%s)\n", e.Author.Name, e.Published.Format("2006-01-02 15:04:05"))
		if e.Content != "" {
			fmt.Fprintf(w, "\n\t%s\n", wrap(html.UnescapeString(e.Content), "\t"))
		}

		for _, e := range e.Comments {
			fmt.Fprintf(w, "\n%s (%s)\n", e.Title, e.Published.Format("2006-01-02 15:04:05"))
			for _, up := range e.Updates {
				if up.Summary != "" {
					fmt.Fprintf(w, "\tSummary: %s\n", up.Summary)
				}
				if up.Owner != "" {
					fmt.Fprintf(w, "\tOwner: %s\n", up.Owner)
				}
				if up.Status != "" {
					fmt.Fprintf(w, "\tStatus: %s\n", up.Status)
				}
				if up.Label != "" {
					fmt.Fprintf(w, "\tLabel: %s\n", up.Label)
				}
			}
			if e.Content != "" {
				fmt.Fprintf(w, "\n\t%s\n", wrap(html.UnescapeString(e.Content), "\t"))
			}
		}
	}
}

func printList(w io.Writer, entries []Entry) {
	for _, e := range entries {
		id := e.ID
		if i := strings.Index(id, "id="); i >= 0 {
			id = id[:i+len("id=")]
		}
		fmt.Fprintf(w, "%s\t%s\n", id, e.Title)
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

var client = http.DefaultClient

func login() error {
	if false {
		data, err := ioutil.ReadFile("../../authblob")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &auth); err != nil {
			return err
		}
	}

	auth.APIClientID = "993255737644.apps.googleusercontent.com"
	auth.APIClientSecret = "kjB02zudLVECBmJdKVMaZluI"

	tr, err := oauthprompt.GoogleToken(".token-code.google.com", auth.APIClientID, auth.APIClientSecret, "https://code.google.com/feeds/issues")
	if err != nil {
		return err
	}
	client = tr.Client()
	return nil
}
