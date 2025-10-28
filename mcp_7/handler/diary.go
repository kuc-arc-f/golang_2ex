package handler

import (
	"bufio"
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	//"os"

	"example.com/go-mcp-server-7/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

/**
*
* @param
*
* @return
*/
func DiaryAddHandler(writer *bufio.Writer, req models.JSONRPCRequest) {
	type DiaryRequest struct {
			Text string `json:"text"`
	}	
	var params models.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}
	var args DiaryRequest
	if err := json.Unmarshal(params.Arguments, &args); err != nil {
		sendError(writer, req.ID, -32602, "Invalid arguments")
		return
	}
	
	jsonBytes, err := json.Marshal(args)
	if err != nil {
		log.Fatalf("JSONへの変換に失敗しました: %v", err)
	}
	jsonString := string(jsonBytes)
	log.Printf(" jsonString %s", jsonString)

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := "INSERT INTO mcp_diary (data) VALUES (?)"
	log.Printf("sql= %s", sql)

	result, err := db.Exec(sql, jsonString)
	if err != nil {
			log.Printf("failed to insert user: %v", err)
			sendError(writer, req.ID, -32602, "Invalid arguments")
			return
	}	
	log.Printf("%v", result)

	toolResult := models.ToolResult{
		Content: []models.Content{
			{
				Type: "text",
				Text: fmt.Sprintf("登録しました。"),
			},
		},
	}

	sendResponse(writer, req.ID, toolResult)
}


/**
*
* @param
*
* @return
*/
func DiaryListHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := "SELECT id, data, created_at, updated_at FROM mcp_diary ORDER BY created_at DESC LIMIT 5"
	log.Printf("sql= %s", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		sendError(writer, req.ID, -32602, "error, InternalServerError")
		return
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
  		sendError(writer, req.ID, -32602, "error, InternalServerError")
			return
		}
		todo.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		todo.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		todos = append(todos, todo)
	}

	jsonBytes, err := json.Marshal(todos)
	if err != nil {
			fmt.Println("JSON 変換エラー:", err)
  		sendError(writer, req.ID, -32602, "error, json convert")
			return
	}
	jsonString := string(jsonBytes)

	toolResult := models.ToolResult{
		Content: []models.Content{
			{
				Type: "text",
				Text: jsonString,
			},
		},
	}

	sendResponse(writer, req.ID, toolResult)
}