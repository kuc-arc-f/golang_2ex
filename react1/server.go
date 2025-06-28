package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "path"
)
type HomeData struct {
    Title  string
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

    fmt.Println("Server running on :8080")
    log.Println("Listening on :8080...")
    http.ListenAndServe(":8080", nil)
}
