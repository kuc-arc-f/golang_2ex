package handler

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	//"os"

	"example.com/go-mcp-server-7/config"
	"example.com/go-mcp-server-7/models"
	"github.com/xuri/excelize/v2"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func connectDB() (*sql.DB, error) {
  dbURL := config.TURSO_DATABASE_URL
  authToken := config.TURSO_AUTH_TOKEN

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
func PurchaseHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
	var params models.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}
	var args models.PurchaseParams
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

	sql := "INSERT INTO item_price (data) VALUES (?)"
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
				Text: fmt.Sprintf("購入情報\n品名: %s\n価格: %d円", args.Name, args.Price),
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
func PurchaseListHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
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

/**
*
* @param
*
* @return
*/
func PurchaseListExcelHnadler(writer *bufio.Writer, req models.JSONRPCRequest) {
  type ArgmentsData struct {
		Template string `json:"template_purchase"`
		OutDir  string  `json:"xls_out_dir"`
  }

  type DataItem struct {
		Name string `json:"name"`
		Price  int  `json:"price"`
  }
	var params models.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(writer, req.ID, -32602, "Invalid params")
		return
	}

	var args ArgmentsData
	if err := json.Unmarshal(params.Arguments, &args); err != nil {
		sendError(writer, req.ID, -32602, "Invalid arguments")
		return
	}
	//log.Printf("args: %v", args)
	log.Printf("args.Template: %v", args.Template)
	log.Printf("args.OutDir: %v", args.OutDir)

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var xls_name = args.Template
	f, err := excelize.OpenFile(xls_name)
	if err != nil {
			fmt.Println(err)
			return
	}	
	defer func() {
			// Close the spreadsheet.
			if err := f.Close(); err != nil {
					fmt.Println(err)
			}
	}()


	sql := "SELECT id, data, created_at, updated_at FROM item_price ORDER BY created_at DESC LIMIT 10"
	log.Printf("sql= %s", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		sendError(writer, req.ID, -32602, "error, InternalServerError")
		return
	}
	defer rows.Close()

	var todos []models.Item
	var sheet_name = "Sheet1"
	var count = 2
	var out_str = ""
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
		log.Printf("data=%s\n", string(todo.Data))
		var itemdata DataItem
		err = json.Unmarshal([]byte(todo.Data), &itemdata)
    if err != nil {
        fmt.Println("JSON Unmarshalエラー:", err)
        return
    }		
		fmt.Printf("name: %s, price: %d\n", itemdata.Name, itemdata.Price)
		out_str += fmt.Sprintf("* ID: %d , name: %s, price: %d\n", todo.ID, itemdata.Name, itemdata.Price)
		var a_col_str = fmt.Sprintf("A%d", count)
		var b_col_str = fmt.Sprintf("B%d", count)
		var c_col_str = fmt.Sprintf("C%d", count)
	  f.SetCellValue(sheet_name, a_col_str, todo.ID)
	  f.SetCellValue(sheet_name, b_col_str, itemdata.Name)
	  f.SetCellValue(sheet_name, c_col_str, itemdata.Price)
		count = count + 1
	}
	now := time.Now()
	milliseconds := now.UnixMilli()
	var out_filename = fmt.Sprintf("output_%d.xlsx", milliseconds)
	var out_file_path = fmt.Sprintf("%s/%s", args.OutDir, out_filename)
	if err := f.SaveAs(out_file_path); err != nil {
			log.Fatal(err)
	}
	out_str += "***\n* 下記リンクをおすと、ダウンロードできます。\n\n"
	out_str += fmt.Sprintf("http://localhost:3000/data/%s\n", out_filename)
  /*
	jsonBytes, err := json.Marshal(todos)
	if err != nil {
			fmt.Println("JSON 変換エラー:", err)
  		sendError(writer, req.ID, -32602, "error, json convert")
			return
	}
	jsonString := string(jsonBytes)
	*/

	toolResult := models.ToolResult{
		Content: []models.Content{
			{
				Type: "text",
				Text: out_str,
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
