// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"os"
	"time"
	"strings"
	"path"
	"sort"
)

func init() {
	os.Chdir(os.Getenv("HOME") + "/blog")
	http.HandleFunc("/", serve)
}

var funcMap = template.FuncMap{
	"now": time.Now,
	"date": timeFormat,
}

func timeFormat(fmt string, t time.Time) string {
	return t.Format(fmt)
}

type blogTime struct {
	time.Time
}

var timeFormats = []string{
	time.RFC3339,
	"Monday, January 2, 2006",
	"January 2, 2006 15:00 -0700",
}

func (t *blogTime) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	for _, f := range timeFormats {
		tt, err := time.Parse(`"` + f + `"`, str)
		if err == nil {
			t.Time = tt
			return nil
		}
	}
	return fmt.Errorf("did not recognize time: %s", str)
}		

type PostData struct {
	Title string
	Date blogTime
	Name string
	OldURL string
	Summary string
	Favorite bool

	PlusAuthor string  // Google+ ID of author
	PlusPage string  // Google+ Post ID for comment post
	PlusAPIKey string // Google+ API key
	PlusURL string
}

func (d *PostData) IsDraft() bool {
	return d.Date.IsZero() || d.Date.After(time.Now())
}

const plusRsc = "116810148281701144465"
const plusKey = "AIzaSyB_JO6hyAJAL659z0Dmu0RUVVvTx02ZPMM"

var replacer = strings.NewReplacer(
	"⁰", "<sup>0</sup>",
	"¹", "<sup>1</sup>",
	"²", "<sup>2</sup>",
	"³", "<sup>3</sup>",
	"⁴", "<sup>4</sup>",
	"⁵", "<sup>5</sup>",
	"⁶", "<sup>6</sup>",
	"⁷", "<sup>7</sup>",
	"⁸", "<sup>8</sup>",
	"⁹", "<sup>9</sup>",
)

func serve(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprint(err), 500)
		}
	}()
	
	p := path.Clean(req.URL.Path)
	if p != req.URL.Path {
		http.Redirect(w, req, p, http.StatusFound)
		return
	}

	if p == "" || p == "/" || p == "/draft" {
		toc(w, req, p=="/draft")
		return
	}

	draft := false
	if strings.HasPrefix(p, "/draft/") {
		draft = true
		p = p[len("/draft"):]
	}		
	
	if strings.Contains(p[1:], "/") {
		http.Error(w, "No such page, sorry.", 404)
		return
	}

	if strings.Contains(p, ".") {
		http.ServeFile(w, req, "img/"+p)
		return
	}

	t := mainTemplate()	
	meta, article := loadPost(p)
	if !draft && meta.IsDraft() {
		http.Error(w, "No such page, sorry.", 404)
		return
	}
	template.Must(t.New("article").Parse(article))

	var buf bytes.Buffer
	if err := t.Execute(&buf, meta); err != nil {
		panic(err)
	}
	w.Write(buf.Bytes())
}

func mainTemplate() *template.Template {
	t := template.New("main")
	t.Funcs(funcMap)

	main, err := ioutil.ReadFile("main.html")
	if err != nil {
		panic(err)
	}
	_, err = t.Parse(string(main))
	if err != nil {
		panic(err)
	}
	return t
}

func loadPost(name string) (meta *PostData, article string) {
	meta = &PostData{
		Name: name,
		Title: "TITLE HERE",
		PlusAuthor: plusRsc,
		PlusAPIKey: plusKey,
	}

	art, err := ioutil.ReadFile("post/" + name)
	if err != nil {
		panic(err)
	}
	if bytes.HasPrefix(art, []byte("{\n")) {
		i := bytes.Index(art, []byte("\n}\n"))
		if i < 0 {
			panic("cannot find end of json metadata")
		}
		hdr, rest := art[:i+3], art[i+3:]
		if err := json.Unmarshal(hdr, meta); err != nil {
			panic(fmt.Sprintf("loading %s: %s", name, err))
		}
		art = rest
	}

	return meta, replacer.Replace(string(art))
}

type byTime []*PostData

func (x byTime) Len() int { return len(x) }
func (x byTime) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byTime) Less(i, j int) bool { return x[i].Date.Time.After(x[j].Date.Time) }

type TocData struct {
	Draft bool
	Posts []*PostData
}

func toc(w http.ResponseWriter, req *http.Request, draft bool) {
	dir, err := ioutil.ReadDir("post/")
	if err != nil {
		panic(err)
	}
	
	var all []*PostData
	for _, d := range dir {
		meta, _ := loadPost(d.Name())
		if meta.IsDraft() != draft {
			continue
		}
		all = append(all, meta)
	}
	
	sort.Sort(byTime(all))
	
	var buf bytes.Buffer
	t := mainTemplate()
	if err := t.Lookup("toc").Execute(&buf, &TocData{draft, all}); err != nil {
		panic(err)
	}
	w.Write(buf.Bytes())
}
