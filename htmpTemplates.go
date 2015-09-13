package main

import (
	"html/template"
	"net/http"
	"path"
)

type Book struct {
	Title string
	Author string
}

var tmpl *template.Template

func main() {
	parseHtml()
	http.HandleFunc("/", ShowBooks)
	http.ListenAndServe(":3000", nil)
}

func parseHtml() {
	fp := path.Join("templates", "index.html")
	temp, _ := template.ParseFiles(fp)
	tmpl = temp
}

func ShowBooks(w http.ResponseWriter, r *http.Request) {
	book := Book{"LA Weekly", "LAWEEKLY.COM"}

	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}