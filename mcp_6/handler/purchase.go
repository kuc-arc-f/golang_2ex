package handler

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"os"

	"example.com/go-remote-mcp-server6/models"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

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

/**
*
* @param
*
* @return
*/
func PurchaseHnadler(name string, price int64) {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type PurchaseItem struct {
		Name string `json:"name"`
		Price  int64    `json:"price"`
	}	
	item := PurchaseItem{
		Name: name,
		Price:  price,
	}

	jsonBytes, err := json.Marshal(item)
	if err != nil {
		log.Fatalf("JSONへの変換に失敗しました: %v", err)
	}
	jsonString := string(jsonBytes)
	log.Printf(" jsonString %s", jsonString)

	sql := "INSERT INTO item_price (data) VALUES (?)"
	log.Printf("sql= %s", sql)

	result, err := db.Exec(sql, jsonString)
	if err != nil {
			log.Printf("erro , failed to insert data")
	}	
	log.Printf("%v", result)	
}


/**
*
* @param
*
* @return
*/
func PurchaseListHnadler() string {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := "SELECT id, data, created_at, updated_at FROM item_price ORDER BY created_at DESC LIMIT 5"
	log.Printf("sql= %s", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		return "error, InternalServerError"
	}
	defer rows.Close()

	var todos []models.Item
	for rows.Next() {
		var todo models.Item
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&todo.ID, &todo.Data, 
			&createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning todo: %v", err)
			return "error, InternalServerError"
		}
		todo.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		todo.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		todos = append(todos, todo)
	}

	jsonBytes, err := json.Marshal(todos)
	if err != nil {
			fmt.Println("JSON 変換エラー:", err)
			return "error, json convert"
	}
	jsonString := string(jsonBytes)
	return jsonString
}

/**
*
* @param
*
* @return
*/
func sendResponse(writer *bufio.Writer, id interface{}, result interface{}) {
	resp := models.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	data, _ := json.Marshal(resp)
	writer.Write(data)
	writer.WriteByte('\n')
	writer.Flush()
}

func sendError(writer *bufio.Writer, id interface{}, code int, message string) {
	resp := models.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &models.RPCError{
			Code:    code,
			Message: message,
		},
	}
	
	data, _ := json.Marshal(resp)
	writer.Write(data)
	writer.WriteByte('\n')
	writer.Flush()
}
