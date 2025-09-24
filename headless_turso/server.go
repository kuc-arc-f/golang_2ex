package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"example.com/headlessturso/handler"

	"github.com/joho/godotenv"
)
var sessionName = "gosession"

type HomeData struct {
	Title string
}

var db *sql.DB

func helloHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    if validAuthHandler(w, r) == false {
        return
    }

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

func loginPage(w http.ResponseWriter, r *http.Request) {
    html := `
        <html><body>
        <form method="POST" action="/api/login">
            Username: <input name="username"><br>
            Password: <input name="password" type="password"><br>
            <input type="submit" value="Login">
        </form>
        </body></html>`
    fmt.Fprint(w, html)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    cookie := http.Cookie{
        Name:   sessionName,
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validAuthHandler(w http.ResponseWriter, r *http.Request) bool {
    var ret = false
    cookie, err := r.Cookie(sessionName)
    if err != nil || cookie.Value != "authenticated" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return false
    }
    ret = true
    return ret
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	username := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
    r.ParseForm()
    u := r.FormValue("username")
    p := r.FormValue("password")

    if u == username && p == password {
        // Cookie にログイン情報を保存
        cookie := http.Cookie{
            Name:  sessionName,
            Value: "authenticated",
            Path:  "/",
        }
        http.SetCookie(w, &cookie)
        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
    }
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}	

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", helloHandler)
    //http.HandleFunc("/api/admin/content_list", handler.AdminListHandler)
    //http.HandleFunc("/api/admin/data_list", handler.AdminDataList)
    http.HandleFunc("/api/data/getone", handler.GetoneTodosHandler)
    http.HandleFunc("/api/data/list", handler.ListTodosHandler)
	http.HandleFunc("/api/data/create", handler.CreateTodoHandler)
	http.HandleFunc("/api/data/update", handler.UpdateTodoHandler)
	http.HandleFunc("/api/data/delete", handler.DeleteTodoHandler)
	http.HandleFunc("/api/login", loginHandler)

	http.HandleFunc("/login", loginPage)
    http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
