package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/clementine/todo-list/store/csv"
)

func showTable(table [][]string) {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', tabwriter.TabIndent)
	for _, row := range table {
		s := strings.Join(row, "\t")
		fmt.Fprintln(w, s)
	}
	w.Flush()
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("todo-list [OPTION]")
		return
	}

	store := csv.NewDefaultCsvStore()
	defer store.Close()

	argsNoProg := os.Args[1:]
	switch option := argsNoProg[0]; option {
	case "add":
		// add todo
		store.Add(argsNoProg[1])
		showTable(store.List())
	case "delete":
		// delete todo
		store.Delete(argsNoProg[1])
		showTable(store.List())
	case "update":
		// update todo
		if len(argsNoProg) < 3 {
			fmt.Println("Not enough params.")
			fmt.Println("todo-list update {id} {task}")
			return
		}
		store.Update(argsNoProg[1], argsNoProg[2])
		showTable(store.List())
	case "list":
		// show todos
		records := store.List()
		showTable(records)
	case "clean":
		store.Clean()
		showTable(store.List())
	default:
		// unknown option
	}
}
