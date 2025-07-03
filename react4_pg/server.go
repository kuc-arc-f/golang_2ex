package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
    "html/template"
	"log"
	"net/http"
    "path"
	"time"

	_ "github.com/lib/pq"
)
type HomeData struct {
    Title  string
}

// Todo defines the structure for a todo item
type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Completed   bool   `json:"completed"`
	ContentType string `json:"content_type"`
	IsPublic    bool   `json:"is_public"`
	FoodOrange  bool   `json:"food_orange"`
	FoodApple   bool   `json:"food_apple"`
	FoodBanana  bool   `json:"food_banana"`
	FoodMelon   bool   `json:"food_melon"`
	FoodGrape   bool   `json:"food_grape"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTodoRequest defines the structure for creating a new todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Completed   bool   `json:"completed"`
	ContentType string `json:"content_type"`
	IsPublic    bool   `json:"is_public"`
	FoodOrange  bool   `json:"food_orange"`
	FoodApple   bool   `json:"food_apple"`
	FoodBanana  bool   `json:"food_banana"`
	FoodMelon   bool   `json:"food_melon"`
	FoodGrape   bool   `json:"food_grape"`
}

// DeleteTodoRequest defines the structure for deleting a todo
type DeleteTodoRequest struct {
	ID int `json:"id"`
}

// UpdateTodoRequest defines the structure for updating a todo
type UpdateTodoRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Completed   bool   `json:"completed"`
	ContentType string `json:"content_type"`
	IsPublic    bool   `json:"is_public"`
	FoodOrange  bool   `json:"food_orange"`
	FoodApple   bool   `json:"food_apple"`
	FoodBanana  bool   `json:"food_banana"`
	FoodMelon   bool   `json:"food_melon"`
	FoodGrape   bool   `json:"food_grape"`
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
	rows, err := db.Query("SELECT id, title, content, completed, content_type, is_public, food_orange, food_apple, food_banana, food_melon, food_grape, created_at, updated_at FROM todos ORDER BY id DESC")
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &todo.ContentType,
			&todo.IsPublic, &todo.FoodOrange, &todo.FoodApple, &todo.FoodBanana,
			&todo.FoodMelon, &todo.FoodGrape, &createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning todo: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todo.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		todo.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
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

	sqlStmt := `
    INSERT INTO todos (title, content, completed, content_type, is_public, food_orange, food_apple, food_banana, food_melon, food_grape)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id, created_at, updated_at`
	var id int
	var createdAt, updatedAt time.Time
	err = db.QueryRow(sqlStmt, req.Title, req.Content, req.Completed, req.ContentType, req.IsPublic, req.FoodOrange, req.FoodApple, req.FoodBanana, req.FoodMelon, req.FoodGrape).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("New record ID is:", id)

	todo := Todo{
		ID:          id,
		Title:       req.Title,
		Content:     req.Content,
		Completed:   req.Completed,
		ContentType: req.ContentType,
		IsPublic:    req.IsPublic,
		FoodOrange:  req.FoodOrange,
		FoodApple:   req.FoodApple,
		FoodBanana:  req.FoodBanana,
		FoodMelon:   req.FoodMelon,
		FoodGrape:   req.FoodGrape,
		CreatedAt:   createdAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   updatedAt.Format("2006-01-02 15:04:05"),
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

	fmt.Printf("delete.targetId:  %d\n", req.ID)
	sqlStatement := `
    DELETE FROM todos
    WHERE id = $1;`
	res, err := db.Exec(sqlStatement, req.ID)
	if err != nil {
		log.Printf("Error deleting todo: %v", err)
		http.Error(w, "Error deleting todo", http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Error deleting todo", http.StatusInternalServerError)
		return
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
	now := time.Now()
	sqlStatement := `
    UPDATE todos
    SET title = $2, content = $3, completed = $4, content_type = $5, is_public = $6, food_orange = $7, food_apple = $8, food_banana = $9, food_melon = $10, food_grape = $11, updated_at = $12
    WHERE id = $1;`
	res, err := db.Exec(sqlStatement, req.ID, req.Title, req.Content, req.Completed, req.ContentType, req.IsPublic, req.FoodOrange, req.FoodApple, req.FoodBanana, req.FoodMelon, req.FoodGrape, now)
	if err != nil {
		log.Printf("Error updating todo: %v", err)
		http.Error(w, "error , db.Exec", http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "error , res.RowsAffected", http.StatusInternalServerError)
		return
	}
	fmt.Printf("更新された行数: %d\n", count)

	var updatedTodo Todo
	var createdAt, updatedAt time.Time
	err = db.QueryRow("SELECT id, title, content, completed, content_type, is_public, food_orange, food_apple, food_banana, food_melon, food_grape, created_at, updated_at FROM todos WHERE id = $1", req.ID).Scan(
		&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.Content, &updatedTodo.Completed, &updatedTodo.ContentType,
		&updatedTodo.IsPublic, &updatedTodo.FoodOrange, &updatedTodo.FoodApple, &updatedTodo.FoodBanana,
		&updatedTodo.FoodMelon, &updatedTodo.FoodGrape, &createdAt, &updatedAt,
	)
	if err != nil {
		log.Printf("Error fetching updated todo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updatedTodo.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	updatedTodo.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    var err error

    homeData := HomeData{"hello"}
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fp := path.Join("templates", "index.html")
        if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl, err := template.ParseFiles(fp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := tmpl.Execute(w, homeData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    fs := http.FileServer(http.Dir("./public"))
    // "/public/" プレフィックスを外して、中身を返すようにする
    http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/api/list", listTodosHandler)
	http.HandleFunc("/api/create", createTodoHandler)
	http.HandleFunc("/api/delete", deleteTodoHandler)
	http.HandleFunc("/api/update", updateTodoHandler)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}