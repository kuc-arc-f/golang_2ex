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
	_ "github.com/tursodatabase/libsql-client-go/libsql"
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
  Content     string  `json:"content"`
	ID          int `json:"id"`
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
  err := godotenv.Load()
  if err != nil {
    log.Fatalf("Error loading .env file: %s", err)
  }	
  dbURL := os.Getenv("TURSO_DATABASE_URL")
  authToken := os.Getenv("TURSO_AUTH_TOKEN")

  if dbURL == "" || authToken == "" {
      log.Fatal("TURSO_DATABASE_URL または TURSO_AUTH_TOKEN が設定されていません")
  }

  // 接続文字列にトークンを付与
  fullURL := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
  // DB 接続
  db, err := sql.Open("libsql", fullURL)
  if err != nil {
      log.Fatalf("failed to open db: %v", err)
  }
  //defer db.Close()

	return db, nil
}


func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
    err := godotenv.Load()
    if err != nil {
      log.Fatalf("Error loading .env file: %s", err)
    }
    apiKey := os.Getenv("API_KEY")
    authHeader := r.Header.Get("Authorization")

    if apiKey != authHeader {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }
    //fmt.Fprintf(w, "Authorization: %s\n", authHeader)
    //fmt.Println("API_KEY:", apiKey)
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
    log.Printf("req.cont: %s", req.Content)
    //log.Printf("req.data: %v", req.Data)

    db, err := connectDB()
    if err != nil {
      log.Fatal(err)
    }
    defer db.Close()

    sql := fmt.Sprintf("INSERT INTO %s (data) VALUES (?)", req.Content)
    log.Printf("sql= %s", sql)

    result, err := db.Exec(sql, req.Data)
    if err != nil {
        log.Printf("failed to insert user: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    id, err := result.LastInsertId()
    if err != nil {
        log.Printf("failed to get last insert id: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Printf("Inserted user %s with ID %d", req.Content, id)  
    todo := Todo{
      ID:          0,
      Content:     req.Content,
      Data:       req.Data,
      CreatedAt:   "",
      UpdatedAt:   "",
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	apiKey := os.Getenv("API_KEY")
	authHeader := r.Header.Get("Authorization")

	fmt.Println("Authorization: %s\n", authHeader)
  if apiKey != authHeader {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//fmt.Println("API_KEY:", apiKey)

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
	sql := fmt.Sprintf("SELECT id, data, created_at, updated_at FROM %s %s", content, order_sql)

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
			&todo.ID, &todo.Data, 
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
    log.Printf("Error Unauthorized")
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
  sql := fmt.Sprintf("DELETE FROM %s WHERE id = %d;" , req.Content, req.ID)
  log.Printf("sql= %s", sql)

	fmt.Printf("delete.targetId:  %d\n", req.ID)
	res, err := db.Exec(sql)
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

  db, err := connectDB()
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

	if req.ID <= 0 {
		http.Error(w, "Valid ID is required", http.StatusBadRequest)
		return
	}

	if req.Data == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	//now := time.Now()
  sql := fmt.Sprintf("UPDATE %s SET data ='%s' WHERE id = %d;" , req.Content, req.Data, req.ID)
  log.Printf("sql= %s", sql)

	res, err := db.Exec(sql) 

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
	updatedTodo.ID = req.ID
	updatedTodo.Content = req.Content
	updatedTodo.Data = req.Data
	updatedTodo.CreatedAt = ""
	updatedTodo.UpdatedAt = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}


