package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type todoItem struct {
	Number int
	Task   string
	Done   bool
}

type todoList struct {
	Items []todoItem
}

var MyTodoList = todoList{
	Items: []todoItem{
		{Number: 1, Task: "Learn Go", Done: true},
		{Number: 2, Task: "Build a web app", Done: false},
		{Number: 3, Task: "Deploy to production", Done: false},
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
	jsonBytes, err := json.MarshalIndent(MyTodoList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(jsonBytes))
	rendertemplate(w, "index.html")
}

// toggleHandler toggles the done status of a task
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	taskNumber, err := strconv.Atoi(r.FormValue("number"))
	if err != nil {
		log.Fatalf("Error converting task number: %v", err)
	}
	for index, item := range MyTodoList.Items {
		if item.Number == taskNumber {
			MyTodoList.Items[index].Done = !MyTodoList.Items[index].Done
			log.Printf("Task < %s > toggled\n", item.Task)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// addHandler adds a new task to the todo list
func addHandler(w http.ResponseWriter, r *http.Request) {
	reorderTasks()
	r.ParseForm()
	newTask := r.FormValue("task")
	newNumber := len(MyTodoList.Items) + 1
	MyTodoList.Items = append(MyTodoList.Items, todoItem{Number: newNumber, Task: newTask, Done: false})
	log.Printf("Task < %s > added to the list\n", newTask)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	taskNumber, err := strconv.Atoi(r.FormValue("number"))
	if err != nil {
		log.Fatalf("Error converting task number on delete: %v", err)
	}
	for index, item := range MyTodoList.Items {
		if item.Number == taskNumber {
			MyTodoList.Items = append(MyTodoList.Items[:index], MyTodoList.Items[index+1:]...)
			log.Printf("Task < %s > deleted from the list\n", item.Task)
		}
	}
	reorderTasks()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	taskNumber, err := strconv.Atoi(r.FormValue("number"))
	if err != nil {
		log.Fatalf("Error converting task number: %v", err)
	}
	newTask := r.FormValue("task")
	for index, item := range MyTodoList.Items {
		if item.Number == taskNumber {
			MyTodoList.Items[index].Task = newTask
			log.Printf("Task < %s > updated to < %s >\n", item.Task, newTask)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func reorderTasks() {
	for index, _ := range MyTodoList.Items {
		MyTodoList.Items[index].Number = index + 1
	}
}

func rendertemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, MyTodoList)

}
