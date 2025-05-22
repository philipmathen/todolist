package main

import "fmt"

func main() {
	var tasks = []string{
		"Syntax lernen",
		"Webservice erstellen",
		"Job finden",
	}

	fmt.Println("Initial task list:")
	printTasksToConsole(tasks)
	tasks = addTask(tasks, "nicht wichsen")
	printTasksToConsole(tasks)

}

func printTasksToConsole(tasks []string) {
	for i, task := range tasks {
		fmt.Printf("Task %d: %s\n", i+1, task)
	}
}

func addTask(tasks []string, task string) []string {
	newSlice := append(tasks, task)
	fmt.Printf("\nTask <%s> added to the list\n\n", task)
	return newSlice

}
