package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	add := flag.String("add", "", "Add a new todo")
	list := flag.Bool("list", false, "List todos")
	done := flag.Int("done", 0, "Mark todo as done by ID")
	del := flag.Int("delete", 0, "Delete todo by ID")

	flag.Parse()

	todos, err := loadTodos()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch {
	case *add != "":
		id := len(todos) + 1
		todos = append(todos, Todo{ID: id, Title: *add})
		saveTodos(todos)
		fmt.Println("Todo added")

	case *list:
		for _, t := range todos {
			status := " "
			if t.Completed {
				status = "✔"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
		}

	case *done > 0:
		for i := range todos {
			if todos[i].ID == *done {
				todos[i].Completed = true
				saveTodos(todos)
				fmt.Println("Todo marked as done")
				return
			}
		}
		fmt.Println("Todo not found")

	case *del > 0:
		for i, t := range todos {
			if t.ID == *del {
				todos = append(todos[:i], todos[i+1:]...)
				saveTodos(todos)
				fmt.Println("Todo deleted")
				return
			}
		}
		fmt.Println("Todo not found")

	default:
		flag.Usage()
		os.Exit(1)
	}
}
