package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    _ "github.com/lib/pq"
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

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello!")
}
var db *sql.DB

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
    db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    rows, err := db.Query("SELECT id, title FROM todos ORDER BY id DESC")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        //fmt.Printf("id: %d , 名前: %s \n", &todo.ID, &todo.Title)
        //err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.CreatedAt, &todo.UpdatedAt)
        err := rows.Scan(&todo.ID, &todo.Title)
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
    db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CreateTodoRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if req.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }

    now := time.Now().Format("2006-01-02 15:04:05")
    
    sqlStmt := `
    INSERT INTO todos (title, content)
    VALUES ($1, $2)
    RETURNING id`
    var id int
    err = db.QueryRow(sqlStmt, req.Title, req.Content).Scan(&id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println("New record ID is:", id)

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
    db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req DeleteTodoRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if req.ID <= 0 {
        http.Error(w, "Valid ID is required", http.StatusBadRequest)
        return
    }

    //var targetId = 1;
    fmt.Printf("delete.targetId:  %d\n", req.ID)
    sqlStatement := `
    DELETE FROM todos
    WHERE id = $1;`
    res, err := db.Exec(sqlStatement, req.ID)
    if err != nil {
        panic(err)
    }
    count, err := res.RowsAffected()
    if err != nil {
        panic(err)
    }
    fmt.Printf("削除された行数: %d\n", count)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
    db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req UpdateTodoRequest
    err = json.NewDecoder(r.Body).Decode(&req)
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
    var targetId = req.ID;
    sqlStatement := `
    UPDATE todos
    SET title = $2
    WHERE id = $1;`
    res, err := db.Exec(sqlStatement, targetId , req.Title)
    if err != nil {
        log.Fatal(err)
        http.Error(w, "error , db.Exec", http.StatusBadRequest)
        return
    }
    count, err := res.RowsAffected()
    if err != nil {
        log.Fatal(err)
        http.Error(w, "error , res.RowsAffected", http.StatusBadRequest)
        return
    }
    fmt.Printf("更新された行数: %d\n", count)

    var respTodo UpdateTodoRequest
    respTodo.ID = req.ID
    respTodo.Title = req.Title
    respTodo.Content = req.Content

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(respTodo)
}

func main() {
    db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/list", listTodosHandler)
    http.HandleFunc("/create", createTodoHandler)
    http.HandleFunc("/delete", deleteTodoHandler)
    http.HandleFunc("/update", updateTodoHandler)

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

