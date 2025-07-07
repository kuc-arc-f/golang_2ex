package handler

import(
  "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	//"html/template"
	"log"
	"net/http"
  "strconv"
	"time"

  _ "github.com/lib/pq"
  "github.com/xuri/excelize/v2"    
)

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


// Todo defines the structure for a todo item
type Todo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Completed   bool    `json:"completed"`
	ContentType string  `json:"content_type"`
	IsPublic    bool    `json:"is_public"`
	FoodOrange  bool    `json:"food_orange"`
	FoodApple   bool    `json:"food_apple"`
	FoodBanana  bool    `json:"food_banana"`
	FoodMelon   bool    `json:"food_melon"`
	FoodGrape   bool    `json:"food_grape"`
	PubDate1    *string `json:"pub_date1,omitempty"`
	PubDate2    *string `json:"pub_date2,omitempty"`
	PubDate3    *string `json:"pub_date3,omitempty"`
	PubDate4    *string `json:"pub_date4,omitempty"`
	PubDate5    *string `json:"pub_date5,omitempty"`
	PubDate6    *string `json:"pub_date6,omitempty"`
	Qty1        string  `json:"qty1"`
	Qty2        string  `json:"qty2"`
	Qty3        string  `json:"qty3"`
	Qty4        string  `json:"qty4"`
	Qty5        string  `json:"qty5"`
	Qty6        string  `json:"qty6"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateTodoRequest defines the structure for creating a new todo
type CreateTodoRequest struct {
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Completed   bool    `json:"completed"`
	ContentType string  `json:"content_type"`
	IsPublic    bool    `json:"is_public"`
	FoodOrange  bool    `json:"food_orange"`
	FoodApple   bool    `json:"food_apple"`
	FoodBanana  bool    `json:"food_banana"`
	FoodMelon   bool    `json:"food_melon"`
	FoodGrape   bool    `json:"food_grape"`
	PubDate1    *string `json:"pub_date1"`
	PubDate2    *string `json:"pub_date2"`
	PubDate3    *string `json:"pub_date3"`
	PubDate4    *string `json:"pub_date4"`
	PubDate5    *string `json:"pub_date5"`
	PubDate6    *string `json:"pub_date6"`
	Qty1        string  `json:"qty1"`
	Qty2        string  `json:"qty2"`
	Qty3        string  `json:"qty3"`
	Qty4        string  `json:"qty4"`
	Qty5        string  `json:"qty5"`
	Qty6        string  `json:"qty6"`
}

// DeleteTodoRequest defines the structure for deleting a todo
type DeleteTodoRequest struct {
	ID int `json:"id"`
}

// UpdateTodoRequest defines the structure for updating a todo
type UpdateTodoRequest struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Completed   bool    `json:"completed"`
	ContentType string  `json:"content_type"`
	IsPublic    bool    `json:"is_public"`
	FoodOrange  bool    `json:"food_orange"`
	FoodApple   bool    `json:"food_apple"`
	FoodBanana  bool    `json:"food_banana"`
	FoodMelon   bool    `json:"food_melon"`
	FoodGrape   bool    `json:"food_grape"`
	PubDate1    *string `json:"pub_date1"`
	PubDate2    *string `json:"pub_date2"`
	PubDate3    *string `json:"pub_date3"`
	PubDate4    *string `json:"pub_date4"`
	PubDate5    *string `json:"pub_date5"`
	PubDate6    *string `json:"pub_date6"`
	Qty1        string  `json:"qty1"`
	Qty2        string  `json:"qty2"`
	Qty3        string  `json:"qty3"`
	Qty4        string  `json:"qty4"`
	Qty5        string  `json:"qty5"`
	Qty6        string  `json:"qty6"`
}

type XlsxlItems struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
}

var db *sql.DB

func ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT id, title, content, completed, content_type, is_public, food_orange, food_apple, food_banana, food_melon, food_grape, pub_date1, pub_date2, pub_date3, pub_date4, pub_date5, pub_date6, qty1, qty2, qty3, qty4, qty5, qty6, created_at, updated_at FROM todos ORDER BY id DESC")
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
		var pubDate1, pubDate2, pubDate3, pubDate4, pubDate5, pubDate6 sql.NullTime
		var qty1, qty2, qty3, qty4, qty5, qty6 sql.NullString
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &todo.ContentType,
			&todo.IsPublic, &todo.FoodOrange, &todo.FoodApple, &todo.FoodBanana,
			&todo.FoodMelon, &todo.FoodGrape,
			&pubDate1, &pubDate2, &pubDate3, &pubDate4, &pubDate5, &pubDate6,
			&qty1, &qty2, &qty3, &qty4, &qty5, &qty6,
			&createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning todo: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if pubDate1.Valid {
			d := pubDate1.Time.Format("2006-01-02")
			todo.PubDate1 = &d
		}
		if pubDate2.Valid {
			d := pubDate2.Time.Format("2006-01-02")
			todo.PubDate2 = &d
		}
		if pubDate3.Valid {
			d := pubDate3.Time.Format("2006-01-02")
			todo.PubDate3 = &d
		}
		if pubDate4.Valid {
			d := pubDate4.Time.Format("2006-01-02")
			todo.PubDate4 = &d
		}
		if pubDate5.Valid {
			d := pubDate5.Time.Format("2006-01-02")
			todo.PubDate5 = &d
		}
		if pubDate6.Valid {
			d := pubDate6.Time.Format("2006-01-02")
			todo.PubDate6 = &d
		}

		todo.Qty1 = qty1.String
		todo.Qty2 = qty2.String
		todo.Qty3 = qty3.String
		todo.Qty4 = qty4.String
		todo.Qty5 = qty5.String
		todo.Qty6 = qty6.String
		todo.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		todo.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		todos = append(todos, todo)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
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

	toSQL := func(s *string) interface{} {
		if s == nil || *s == "" {
			return nil
		}
		return *s
	}

	sqlStmt := `
    INSERT INTO todos (title, content, completed, content_type, is_public, food_orange, food_apple, food_banana, food_melon, food_grape, pub_date1, pub_date2, pub_date3, pub_date4, pub_date5, pub_date6, qty1, qty2, qty3, qty4, qty5, qty6)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
    RETURNING id, created_at, updated_at`
	var id int
	var createdAt, updatedAt time.Time
	err = db.QueryRow(sqlStmt, req.Title, req.Content, req.Completed, req.ContentType, req.IsPublic, req.FoodOrange, req.FoodApple, req.FoodBanana, req.FoodMelon, req.FoodGrape, toSQL(req.PubDate1), toSQL(req.PubDate2), toSQL(req.PubDate3), toSQL(req.PubDate4), toSQL(req.PubDate5), toSQL(req.PubDate6), req.Qty1, req.Qty2, req.Qty3, req.Qty4, req.Qty5, req.Qty6).Scan(&id, &createdAt, &updatedAt)
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
		PubDate1:    req.PubDate1,
		PubDate2:    req.PubDate2,
		PubDate3:    req.PubDate3,
		PubDate4:    req.PubDate4,
		PubDate5:    req.PubDate5,
		PubDate6:    req.PubDate6,
		Qty1:        req.Qty1,
		Qty2:        req.Qty2,
		Qty3:        req.Qty3,
		Qty4:        req.Qty4,
		Qty5:        req.Qty5,
		Qty6:        req.Qty6,
		CreatedAt:   createdAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   updatedAt.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
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

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
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

	toSQL := func(s *string) interface{} {
		if s == nil || *s == "" {
			return nil
		}
		return *s
	}

	sqlStatement := `
    UPDATE todos
    SET title = $2, content = $3, completed = $4, content_type = $5, is_public = $6, food_orange = $7, food_apple = $8, food_banana = $9, food_melon = $10, food_grape = $11, pub_date1 = $12, pub_date2 = $13, pub_date3 = $14, pub_date4 = $15, pub_date5 = $16, pub_date6 = $17, qty1 = $18, qty2 = $19, qty3 = $20, qty4 = $21, qty5 = $22, qty6 = $23, updated_at = $24
    WHERE id = $1;`
	res, err := db.Exec(sqlStatement, req.ID, req.Title, req.Content, req.Completed, req.ContentType, req.IsPublic, req.FoodOrange, req.FoodApple, req.FoodBanana, req.FoodMelon, req.FoodGrape, toSQL(req.PubDate1), toSQL(req.PubDate2), toSQL(req.PubDate3), toSQL(req.PubDate4), toSQL(req.PubDate5), toSQL(req.PubDate6), req.Qty1, req.Qty2, req.Qty3, req.Qty4, req.Qty5, req.Qty6, now)
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
	var pubDate1, pubDate2, pubDate3, pubDate4, pubDate5, pubDate6 sql.NullTime
	var qty1, qty2, qty3, qty4, qty5, qty6 sql.NullString
	err = db.QueryRow("SELECT id, title, content, completed, content_type, is_public, food_orange, food_apple, food_banana, food_melon, food_grape, pub_date1, pub_date2, pub_date3, pub_date4, pub_date5, pub_date6, qty1, qty2, qty3, qty4, qty5, qty6, created_at, updated_at FROM todos WHERE id = $1", req.ID).Scan(
		&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.Content, &updatedTodo.Completed, &updatedTodo.ContentType,
		&updatedTodo.IsPublic, &updatedTodo.FoodOrange, &updatedTodo.FoodApple, &updatedTodo.FoodBanana,
		&updatedTodo.FoodMelon, &updatedTodo.FoodGrape,
		&pubDate1, &pubDate2, &pubDate3, &pubDate4, &pubDate5, &pubDate6,
		&qty1, &qty2, &qty3, &qty4, &qty5, &qty6,
		&createdAt, &updatedAt,
	)
	if err != nil {
		log.Printf("Error fetching updated todo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if pubDate1.Valid {
		d := pubDate1.Time.Format("2006-01-02")
		updatedTodo.PubDate1 = &d
	}
	if pubDate2.Valid {
		d := pubDate2.Time.Format("2006-01-02")
		updatedTodo.PubDate2 = &d
	}
	if pubDate3.Valid {
		d := pubDate3.Time.Format("2006-01-02")
		updatedTodo.PubDate3 = &d
	}
	if pubDate4.Valid {
		d := pubDate4.Time.Format("2006-01-02")
		updatedTodo.PubDate4 = &d
	}
	if pubDate5.Valid {
		d := pubDate5.Time.Format("2006-01-02")
		updatedTodo.PubDate5 = &d
	}
	if pubDate6.Valid {
		d := pubDate6.Time.Format("2006-01-02")
		updatedTodo.PubDate6 = &d
	}

	updatedTodo.Qty1 = qty1.String
	updatedTodo.Qty2 = qty2.String
	updatedTodo.Qty3 = qty3.String
	updatedTodo.Qty4 = qty4.String
	updatedTodo.Qty5 = qty5.String
	updatedTodo.Qty6 = qty6.String
	updatedTodo.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	updatedTodo.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

//var xlsx_todo_items []XlsxlItems

/**
*
* @param
*
* @return
*/
func TodoDownloadHandler(w http.ResponseWriter, r *http.Request) {
  // 全要素をゼロ値（クリア）に
	//for i := range xlsx_todo_items {
	//	xlsx_todo_items[i] = XlsxlItems{}
	//}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}  
	var todoItem []XlsxlItems = getXlsxItems()
  //fmt.Printf(todoItem)
  //fmt.Printf("%+v\n", xlsx_todo_items)

  // xlsx edit
  var xls_name = "test.xlsx"
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
  // 1. Excel ファイルを生成
  var sheet_name = "Sheet1"

  for i, v := range todoItem {
    var row = v;
    var pos = i + 2;
    var targetId = strconv.Itoa(pos)
		//fmt.Printf("arr[%d] = %v\n", i, targetId)
		//fmt.Printf("arr[%d] = %d\n", i, row.ID)
    var cell_a = "A" + targetId
    var cell_b = "B" + targetId
    var cell_c = "C" + targetId
    f.SetCellValue(sheet_name, cell_a, row.ID)
    f.SetCellValue(sheet_name, cell_b, row.Title)
    f.SetCellValue(sheet_name, cell_c, row.Content)
	}

  var buf bytes.Buffer
  if err := f.Write(&buf); err != nil {
      http.Error(w, "Excel 書き込みに失敗しました", http.StatusInternalServerError)
      return
  }

  w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
  w.Header().Set("Content-Disposition", `attachment; filename="edit_report.xlsx"`)
  w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
  
  // 4. レスポンスとして書き出し
  if _, err := w.Write(buf.Bytes()); err != nil {
      // クライアント切断などで Write に失敗する可能性あり
      // 必要に応じてログ出力など
  }
}

/**
*
* @param
*
* @return
*/
func getXlsxItems() []XlsxlItems {
	var todoItem []XlsxlItems
  db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

  rows, err := db.Query("SELECT id, title, content FROM todos ORDER BY id DESC")
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		return todoItem
	}
	defer rows.Close()

	for rows.Next() {
		var todo XlsxlItems
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Content, 
		)
		if err != nil {
			log.Printf("Error scanning todo: %v", err)
			return todoItem
		}

		todoItem = append(todoItem, todo)
	}  
	return todoItem
}