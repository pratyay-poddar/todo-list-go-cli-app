package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pratyay-poddar/simple-to-do/todo"
	_ "modernc.org/sqlite"
)

func DbInit(filePath string) (*sql.DB, error) {

	// checking the validity of the directory
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("Unable to find the dir at provided File Path :: %v", err)
	}

	// establishment of the connection with the database
	db, err := sql.Open("sqlite", filePath)

	if err != nil {
		return nil, fmt.Errorf("Failed to establish connection with database ::%v", err)
	}

	// checking the communication protocol
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("communication protocol failed :: %v", err)
	}

	// creation of the tabular data structure
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        title TEXT PRIMARY KEY,
        status BOOLEAN NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL,
        completed_at DATETIME
    );`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to make the query table :: %v", err)
	}

	// no error all OK
	fmt.Println("all OK : 200")
	return db, nil
}

///

// SaveTask writes a newly added task from your Go slice into a permanent database row.
func SaveTask(db *sql.DB, item todo.Item) error {
	query := `INSERT INTO todos (title, status, created_at, completed_at) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, item.Title, item.Status, item.CreatedAt, item.CompletedAt)
	if err != nil {
		return fmt.Errorf("could not save task to database: %w", err)
	}
	return nil
}

///

// LoadTasks reads your SQLite file on boot and returns an populated Items collection.
func LoadTasks(db *sql.DB) (todo.Items, error) {
	query := `SELECT title, status, created_at, completed_at FROM todos`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items todo.Items

	for rows.Next() {
		var item todo.Item
		err := rows.Scan(&item.Title, &item.Status, &item.CreatedAt, &item.CompletedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

///

// UpdateTaskStatus changes a task's boolean status flag and logs completion time in SQLite.
func UpdateTaskStatus(db *sql.DB, title string, status bool, completedAt *time.Time) error {
	query := `UPDATE todos SET status = ?, completed_at = ? WHERE title = ?`
	_, err := db.Exec(query, status, completedAt, title)
	return err
}

// DeleteTaskFromDB completely wipes a row from your SQLite file by its tracking title.
func DeleteTaskFromDB(db *sql.DB, title string) error {
	query := `DELETE FROM todos WHERE title = ?`
	_, err := db.Exec(query, title)
	return err
}

///

func main() {
	dbPath := "./database/database.db"
	DB, err := DbInit(dbPath)
	if err != nil {
		log.Fatalf("Critical Database Initialisation Error :: %v", err)
	}
	defer DB.Close()

	// 1. Hydrate the application memory state from SQLite file on initialization
	var items todo.Items
	items, err = LoadTasks(DB)
	if err != nil {
		log.Fatalf("Failed to extract data records from storage module: %v", err)
	}

	// 2. Enforce minimal command arguments input block
	if len(os.Args) < 2 {
		fmt.Println("Usage options: go run cmd/main.go [add|list|complete|delete] [arguments]")
		return
	}

	// 3. Route matching terminal keywords to core business logic and sync to SQLite
	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing parameters! Syntax: go run cmd/main.go add \"Task Name\"")
			return
		}
		title := os.Args[2]

		// Execute memory safety array operations first
		if err := items.Add(title); err != nil {
			fmt.Printf("Business Logic Rejection: %v\n", err)
			return
		}

		// Retrieve latest index object structure and pass downstream to relational table
		freshItem := items[len(items)-1]
		if err := SaveTask(DB, freshItem); err != nil {
			fmt.Printf("Database layer operational error: %v\n", err)
			return
		}
		fmt.Printf("Success: Added task %q into database record mapping.\n", title)

	case "list":
		if len(items) == 0 {
			fmt.Println("The active task manager is empty!")
			return
		}
		fmt.Println("\n=== ACTIVE PROGRAM TO-DO ENTRIES ===")
		for idx, task := range items {
			marker := "[ ]"
			if task.Status {
				marker = "[✔]"
			}
			fmt.Printf("%d. %s %s (Created: %s)\n", idx+1, marker, task.Title, task.CreatedAt.Format("2006-01-02 15:04"))
		}
		fmt.Println()

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing target parameter! Syntax: go run cmd/main.go complete \"Task Name\"")
			return
		}
		title := os.Args[2]

		// Mutate memory reference structure state array
		items.Complete(title)

		// Search state representation matching parameter targets and fire updating raw SQL transaction
		for _, task := range items {
			if task.Title == title {
				if err := UpdateTaskStatus(DB, task.Title, task.Status, task.CompletedAt); err != nil {
					fmt.Printf("Failed storage layer query translation: %v\n", err)
					return
				}
				fmt.Printf("Success: Checked off task %q.\n", title)
				return
			}
		}
		fmt.Printf("Target identification error: Task named %q not tracked.\n", title)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing targeting parameters! Syntax: go run cmd/main.go delete \"Task Name\"")
			return
		}
		title := os.Args[2]

		// Unpack slice and truncate index values cleanly inside domain logic
		items.Delete(title)

		// Execute strict deletion command down to row identifier target
		if err := DeleteTaskFromDB(DB, title); err != nil {
			fmt.Printf("Failed record purge process execution: %v\n", err)
			return
		}
		fmt.Printf("Success: Erased task %q across storage modules.\n", title)

	default:
		fmt.Printf("Unsupported command reference argument: %q. Use add, list, complete, or delete.\n", command)
	}
}
