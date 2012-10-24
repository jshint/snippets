package snippets

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Report struct {
	Code []byte
	Opts []byte
	Date time.Time
}

var templates = make(map[string]*template.Template)

func init() {
	for _, name := range []string{"500", "404", "report"} {
		tmpl := template.Must(template.New(name).ParseFiles("templates/" + name + ".html"))
		templates[name] = tmpl
	}

	http.HandleFunc("/reports/", show)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "http://jshint.com/", 302)
  })
}

func renderTemplate(w http.ResponseWriter, name string, report *Report) {
	err := templates[name].ExecuteTemplate(w, name+".html", report)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func show(w http.ResponseWriter, r *http.Request) {
	var report Report
	c := appengine.NewContext(r)

	intID, _ := strconv.ParseInt(r.URL.Path[9:], 10, 64)
	key := datastore.NewKey(c, "report", "", intID, nil)
	if err := datastore.Get(c, key, &report); err != nil {
		c.Criticalf(err.Error())
		renderTemplate(w, "404", nil)
		return
	}
	renderTemplate(w, "report", &report)
}
