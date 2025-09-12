package handler

import(
	"fmt"
	"net/http"
  "database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
  _ "github.com/lib/pq"
)

// Todo defines the structure for a todo item
type Todo struct {
	ID          int     `json:"id"`
	Content     string  `json:"content"`
	Data        string  `json:"data"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateTodoRequest struct {
	Content     string  `json:"content"`
	Data        string  `json:"data"`
}

type DeleteTodoRequest struct {
	ID int `json:"id"`
}

type UpdateTodoRequest struct {
	ID          int     `json:"id"`
	Content     string  `json:"content"`
	Data        string  `json:"data"`
}

func Test() {
 fmt.Println("handler.TestHandler")
}
func Test2(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "handler.Test2.Hello!")
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

// connectDB はデータベースに接続し、*sql.DBオブジェクトを返します。
func connectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// データベースへの接続を確認
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}


func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	apiKey := os.Getenv("API_KEY")
	authHeader := r.Header.Get("Authorization")

	//fmt.Fprintf(w, "Authorization: %s\n", authHeader)
	//fmt.Println("API_KEY:", apiKey)
  if apiKey != authHeader {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	if req.Data == "" {
		http.Error(w, "Data is required", http.StatusBadRequest)
		return
	}

	sqlStmt := `
    INSERT INTO hcm_data (content, data
		 )
    VALUES ($1, $2)
    RETURNING id, created_at, updated_at`
	var id int
	var createdAt, updatedAt time.Time
	err = db.QueryRow(sqlStmt, req.Content, req.Data,  ).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("New record ID is:", id)

	todo := Todo{
		ID:          id,
		Content:     req.Content,
		Data:       req.Data,
		CreatedAt:   createdAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   updatedAt.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	// ファイルが見つからなくてもエラーにならないように設定することも可能
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	apiKey := os.Getenv("API_KEY")
	authHeader := r.Header.Get("Authorization")

	fmt.Println("Authorization: %s\n", authHeader)
	//fmt.Println("API_KEY:", apiKey)
  if apiKey != authHeader {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//query-string
	query := r.URL.Query()

	content := query.Get("content")
	order := query.Get("order")
	//fmt.Fprintf(w, "content=%s\n", content)
	//fmt.Fprintf(w, "order=%s\n", order)
	order_sql := "ORDER BY created_at ASC";
	if order == "desc" {
		order_sql = "ORDER BY created_at DESC";
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sql := fmt.Sprintf("SELECT id, data, content, created_at, updated_at FROM hcm_data WHERE content ='%s' %s", content, order_sql)

	//fmt.Fprintf(w, "sql=%s\n", sql)
	fmt.Println("sql=%s\n", sql)

	rows, err := db.Query(sql)
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
			&todo.ID, &todo.Data, &todo.Content, 
			&createdAt, &updatedAt,
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


func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	apiKey := os.Getenv("API_KEY")
	authHeader := r.Header.Get("Authorization")

	//fmt.Fprintf(w, "Authorization: %s\n", authHeader)
	//fmt.Println("API_KEY:", apiKey)
  if apiKey != authHeader {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
		
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
    DELETE FROM hcm_data
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


func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	apiKey := os.Getenv("API_KEY")
	authHeader := r.Header.Get("Authorization")

	//fmt.Fprintf(w, "Authorization: %s\n", authHeader)
	//fmt.Println("API_KEY:", apiKey)
  if apiKey != authHeader {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
		
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

	if req.Data == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	now := time.Now()

	sqlStatement := `
    UPDATE hcm_data
    SET data = $2, content = $3, 
		updated_at = $4
    WHERE id = $1;`
	res, err := db.Exec(sqlStatement, 
	req.ID, req.Data, req.Content, 
	now)
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
	err = db.QueryRow("SELECT id, content, data, created_at, updated_at FROM hcm_data WHERE id = $1", req.ID).Scan(
		&updatedTodo.ID, 
		&updatedTodo.Content, &updatedTodo.Data,
		&createdAt, &updatedAt,
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


