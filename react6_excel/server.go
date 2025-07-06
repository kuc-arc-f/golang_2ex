package main

import (
    "bytes"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"
    "path"

    "github.com/xuri/excelize/v2"    
)
type HomeData struct {
    Title  string
}

func editXLSX(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    //var xls_name = "./data/test.xlsx"
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
    f.SetCellValue(sheet_name, "A4", "hoge7")
    f.SetCellValue(sheet_name, "B4", 1014)

    // 元テンプレートは、保存しない。
    //if err := f.SaveAs(xls_name); err != nil {
    //    log.Fatal(err)
    //}

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

func serveXLSX(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // 1. Excel ファイルを生成
    f := excelize.NewFile()
    var sheet_name = "Sheet1"
    index, _ := f.NewSheet(sheet_name)
    f.SetCellValue(sheet_name, "A1", "こんにちは")
    f.SetCellValue(sheet_name, "A2", "hoge2")
    f.SetActiveSheet(index)

    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        http.Error(w, "Excel 書き込みに失敗しました", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", `attachment; filename="report.xlsx"`)
    w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

    // 4. レスポンスとして書き出し
    if _, err := w.Write(buf.Bytes()); err != nil {
        // クライアント切断などで Write に失敗する可能性あり
        // 必要に応じてログ出力など
    }
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

func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Goodbye!")
}

func main() {
    fs := http.FileServer(http.Dir("./public"))
    // "/public/" プレフィックスを外して、中身を返すようにする
    http.Handle("/public/", http.StripPrefix("/public/", fs))

    http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/goodbye", goodbyeHandler)
    http.HandleFunc("/download", serveXLSX)
    http.HandleFunc("/edit_excel", editXLSX)

    fmt.Println("Server running on :8080")
    log.Println("Listening on :8080...")
    http.ListenAndServe(":8080", nil)
}
