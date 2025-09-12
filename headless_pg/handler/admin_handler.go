package handler

import(
	"fmt"
	"net/http"
  //"database/sql"
	"encoding/json"
	"log"
	//"os"
	"time"

	//"github.com/joho/godotenv"
  _ "github.com/lib/pq"
)

type AdminItem struct {
	ID          int     `json:"id"`
	Content     string  `json:"content"`
	Data        string  `json:"data"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func AdminDataList(w http.ResponseWriter, r *http.Request) {

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//query-string
	query := r.URL.Query()

	content := query.Get("content")
	order := query.Get("order")
	order_sql := "ORDER BY created_at ASC";
	if order == "desc" {
		order_sql = "ORDER BY created_at DESC";
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sql := fmt.Sprintf("SELECT id, data, content, created_at, updated_at FROM hcm_data WHERE content ='%s' %s", content, order_sql)

	fmt.Println("sql=%s\n", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []AdminItem
	for rows.Next() {
		var todo AdminItem
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


func AdminListHandler(w http.ResponseWriter, r *http.Request) {

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sql := fmt.Sprintf("SELECT distinct content FROM hcm_data;")
	fmt.Println("sql=%s\n", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []string
	for rows.Next() {
		var todo string
		err := rows.Scan(
			&todo,
		)
		if err != nil {
			log.Printf("Error scanning todo: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

