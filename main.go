package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type todoItem struct {
	Id   string
	Task string
	Done bool
}

type todoList struct {
	Items []todoItem
}

var MyTodoList = todoList{
	Items: []todoItem{
		{Id: "1", Task: "Learn Go", Done: true},
		{Id: "2", Task: "Build a web app", Done: false},
		{Id: "3", Task: "Deploy to production", Done: false},
	},
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/update", updateHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// indexHandler serves the index page
// It renders the template and passes the todo list to it
func indexHandler(w http.ResponseWriter, r *http.Request) {
	rendertemplate(w, "index.html")
}

// toggleHandler toggles the done status of a task
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	taskId := r.FormValue("id")
	for index, item := range MyTodoList.Items {
		if item.Id == taskId {
			MyTodoList.Items[index].Done = !MyTodoList.Items[index].Done
			log.Printf("Task < %s > toggled\n", item.Task)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// addHandler adds a new task to the todo list
func addHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	newTask := r.FormValue("task")
	newId := fmt.Sprintf("%d", len(MyTodoList.Items)+1)
	MyTodoList.Items = append(MyTodoList.Items, todoItem{Id: newId, Task: newTask, Done: false})
	log.Printf("Task <%s> added to the list\n", newTask)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	taskId := r.FormValue("id")
	for index, item := range MyTodoList.Items {
		if item.Id == taskId {
			MyTodoList.Items = append(MyTodoList.Items[:index], MyTodoList.Items[index+1:]...)
			log.Printf("Task <%s> deleted from the list\n", item.Task)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	taskId := r.FormValue("id")
	newTask := r.FormValue("task")
	for index, item := range MyTodoList.Items {
		if item.Id == taskId {
			MyTodoList.Items[index].Task = newTask
			log.Printf("Task <%s> updated to <%s>\n", item.Task, newTask)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func rendertemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, MyTodoList)

}
