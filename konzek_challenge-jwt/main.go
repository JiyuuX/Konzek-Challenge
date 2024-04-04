/*
KONZEK-Backend Developer Assignment

with JWT authorization, hashing-crypting, register, loging etc.

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

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// database const
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "tasks"
)

var db *sql.DB
var mu sync.Mutex

// jwt private key
var jwtKey = []byte("gizli_bir_anahtar")

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable options='-c client_encoding=UTF8'",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//endpoints
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", authenticate(handler))

	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// handler function for GET, POST, UPDATE, DELETE
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

// Task Struct
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// User Struct
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// writeResponse, to writing HTTP response
func writeResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create task
func createTask(title, description, status string) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := db.Exec("INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)", title, description, status)
	return err
}

// Get ALL Tasks
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

// UPDATE Task by ID
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

// DELETE Task by ID
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

// Register
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Kullanıcı adı, şifre ve e-posta adresi gereklidir.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Şifre hashlenemedi.", http.StatusInternalServerError)
		return
	}

	if err := createUser(user.Username, hashedPassword, user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, err := generateToken(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(tokenString))
}

// LOGIN
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedPassword, err := getUserPassword(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := comparePasswords(storedPassword, user.Password); err != nil {
		http.Error(w, "Kullanıcı adı veya şifre hatalı.", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateToken(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

// Authenticate with JWT
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Yetkisiz erişim.", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Geçersiz token")
			}
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Geçersiz token.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Hashing User Password for Postgres
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Checking user's hashed password whether is valid or not
func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Create User unique by username
func createUser(username, hashedPassword, email string) error {
	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, hashedPassword, email)
	return err
}

// Get User's Hashed Password from DB
func getUserPassword(username string) (string, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

// Generate JWT token with sha-256
func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	return token.SignedString(jwtKey)
}
