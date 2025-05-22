package main

import (
	"html/template"
	"log"
	"net/http"
)

type todoItem struct {
	Task string
	Done bool
}

type todoList struct {
	Items []todoItem
}

var todos = todoList{
	Items: []todoItem{
		{Task: "Learn Go", Done: true},
		{Task: "Build a web app", Done: false},
		{Task: "Deploy to production", Done: false},
	},
}

func main() {
	http.HandleFunc("/", indexHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rendertemplate(w, "index.html")
}

func rendertemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, todos)

}
