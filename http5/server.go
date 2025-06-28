package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type CreateTodoRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}

type DeleteTodoRequest struct {
    ID int `json:"id"`
}

type UpdateTodoRequest struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./todos.db")
    if err != nil {
        log.Fatal(err)
    }

    createTableSQL := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT,
        created_at TEXT,
        updated_at TEXT
    );`

    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatal(err)
    }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello!")
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM todos")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.CreatedAt, &todo.UpdatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        todos = append(todos, todo)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CreateTodoRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if req.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }

    now := time.Now().Format("2006-01-02 15:04:05")
    
    result, err := db.Exec(
        "INSERT INTO todos (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)",
        req.Title, req.Content, now, now,
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    id, err := result.LastInsertId()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    todo := Todo{
        ID:        int(id),
        Title:     req.Title,
        Content:   req.Content,
        CreatedAt: now,
        UpdatedAt: now,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req DeleteTodoRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if req.ID <= 0 {
        http.Error(w, "Valid ID is required", http.StatusBadRequest)
        return
    }

    result, err := db.Exec("DELETE FROM todos WHERE id = ?", req.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req UpdateTodoRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if req.ID <= 0 {
        http.Error(w, "Valid ID is required", http.StatusBadRequest)
        return
    }

    if req.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }

    now := time.Now().Format("2006-01-02 15:04:05")
    
    result, err := db.Exec(
        "UPDATE todos SET title = ?, content = ?, updated_at = ? WHERE id = ?",
        req.Title, req.Content, now, req.ID,
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    var todo Todo
    err = db.QueryRow(
        "SELECT id, title, content, created_at, updated_at FROM todos WHERE id = ?",
        req.ID,
    ).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.CreatedAt, &todo.UpdatedAt)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func main() {
    initDB()
    defer db.Close()

    http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/list", listTodosHandler)
    http.HandleFunc("/create", createTodoHandler)
    http.HandleFunc("/delete", deleteTodoHandler)
    http.HandleFunc("/update", updateTodoHandler)

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

