/*
KONZEK-Backend Developer Assignment

Author/Applicant : Burak ŞEKER
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
)

// Database Connection Var.
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "tasks"
)

var db *sql.DB
var mu sync.Mutex

func main() {
	// Connect Postgres
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable options='-c client_encoding=UTF8'",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start HTTP server and push the DB variables to handler
	http.HandleFunc("/", handler)
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Handler for requests POST, GET, UPDATE, DELETE
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tasks, err := getAllTasks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeResponse(w, tasks)

	case "POST":
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := createTask(task.Title, task.Description, task.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Yeni görev oluşturuldu.")

	case "DELETE":
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, "Görev ID'si belirtilmedi", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(taskID)
		if err != nil {
			http.Error(w, "Geçersiz görev ID'si", http.StatusBadRequest)
			return
		}

		if err := deleteTask(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Görev başarıyla silindi.")

	case "PUT":
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, "Görev ID'si belirtilmedi", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(taskID)
		if err != nil {
			http.Error(w, "Geçersiz görev ID'si", http.StatusBadRequest)
			return
		}

		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := updateTask(id, task.Title, task.Description, task.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Görev başarıyla güncellendi.")

	default:
		http.Error(w, "Geçersiz istek methodu", http.StatusMethodNotAllowed)
	}
}

// Görev yapısı
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// writeResponse, HTTP yanıtını yazmak için
func writeResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create Task
func createTask(title, description, status string) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := db.Exec("INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)", title, description, status)
	return err
}

// Get All Tasks
func getAllTasks() ([]Task, error) {
	mu.Lock()
	defer mu.Unlock()
	rows, err := db.Query("SELECT id, title, description, status FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Update Task
func updateTask(id int, title, description, status string) error {
	mu.Lock()
	defer mu.Unlock()
	result, err := db.Exec("UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4", title, description, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Belirtilen ID'ye sahip görev bulunamadı")
	}

	return nil
}

// Delete Task
func deleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()
	result, err := db.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Belirtilen ID'ye sahip görev bulunamadı")
	}

	return nil
}
