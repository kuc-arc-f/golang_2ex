package handler

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	//"os"

	"example.com/go-mcp-server-9/config"
	"example.com/go-mcp-server-9/models"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname,
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

/**
*
* @param
*
* @return
*/
func TestCreateHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
	type TestParams struct {
		Title string `json:"title"`
		Content string `json:"content"`
	}

	var params models.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}
	log.Printf("params %v", params)

	var args TestParams
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

	sqlStmt := `
    INSERT INTO test (title, content
		 )
    VALUES ($1, $2)
    RETURNING id`

	var id int
	err = db.QueryRow(sqlStmt, args.Title, args.Content ).Scan(&id)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		sendError(writer, req.ID, -32602, "Invalid arguments")
		return
	}

	toolResult := models.ToolResult{
		Content: []models.Content{
			{
				Type: "text",
				Text: "OK",
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
func TestListHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
	type TestInputParams struct {
		Content string `json:"content"`
	}
	type ListParam struct {
		ID          int     `json:"id"`
		Content     string  `json:"content"`
		Data        string  `json:"data"`
		CreatedAt   string  `json:"created_at"`
		UpdatedAt   string  `json:"updated_at"`
	}

	var params models.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}
	log.Printf("params %v", params)
	var args TestInputParams
	if err := json.Unmarshal(params.Arguments, &args); err != nil {
		sendError(writer, req.ID, -32602, "Invalid arguments")
		return
	}
	log.Printf("args %v", args)
	log.Printf("args.Content %s", args.Content)

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := fmt.Sprintf(`SELECT id, title, content, created_at, updated_at FROM test
	ORDER BY created_at ASC`)
	log.Printf("sql= %s", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		sendError(writer, req.ID, -32602, "error, InternalServerError")
		return
	}
	defer rows.Close()

	var todos []ListParam
	for rows.Next() {
		var todo ListParam
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&todo.ID, &todo.Data, &todo.Content, 
			&createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning data: %v", err)
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

/**
*
* @param
*
* @return
*/
func TestDeleteHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
	type TestParams struct {
		Id int     `json:"id"`
		Content string `json:"content"`
	}

	var params models.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}
	log.Printf("params %v", params)

	var args TestParams
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

	sqlStmt := `
    DELETE FROM test
    WHERE id = $1;`

	res, err := db.Exec(sqlStmt, args.Id)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		sendError(writer, req.ID, -32602, "Invalid arguments")
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		sendError(writer, req.ID, -32602, "NG, delete")
		return
	}
	fmt.Printf("削除された行数: %d\n", count)

	toolResult := models.ToolResult{
		Content: []models.Content{
			{
				Type: "text",
				Text: "OK",
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
func TestUpdateHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
	type TestUpdateParams struct {
		Id int `json:"id"`
		Title string `json:"title"`
		Content string `json:"content"`
	}

	var params models.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}
	log.Printf("params %v", params)
	log.Printf("params.Arguments %s", string(params.Arguments))

	var args TestUpdateParams
	if err := json.Unmarshal(params.Arguments, &args); err != nil {
		sendError(writer, req.ID, -32602, "Invalid arguments")
		return
	}
	log.Printf("args.id=%d", args.Id)
	
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

	sqlStatement := `
    UPDATE test
    SET title = $2, content = $3 
    WHERE id = $1;`
	res, err := db.Exec(sqlStatement, 
	args.Id, args.Title, args.Content)

	if err != nil {
		log.Printf("Error updating item: %v", err)
		sendError(writer, req.ID, -32602, "NG updating, Invalid arguments")
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		sendError(writer, req.ID, -32602, "NG RowsAffected")
		return
	}
	fmt.Printf("更新された行数: %d\n", count)	

	toolResult := models.ToolResult{
		Content: []models.Content{
			{
				Type: "text",
				Text: "OK",
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
