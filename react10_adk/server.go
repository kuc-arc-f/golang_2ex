package main

import (
    //"bytes"
    //"encoding/json"
    "fmt"
    "html/template"
    //"io"
    "log"
    "net/http"
    "time"
    "path"

    "react10/handler"
)
var username = "admin"
var password = "password"
var sessionName = "gosession"

type HomeData struct {
    Title  string
}
// LoggerMiddleware はリクエストのログを記録するミドルウェアです。
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // 次のハンドラまたはミドルウェアを呼び出す
		log.Printf("[%s] %s %s %s", r.Method, r.RequestURI, time.Since(start), r.RemoteAddr)
	})
}

// AuthMiddleware は簡易的な認証を行うミドルウェアです。
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "password" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
// ログインページ
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

// ログイン処理
func loginHandler(w http.ResponseWriter, r *http.Request) {
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
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    } else {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
    }
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


// 認証されたユーザーのみ表示
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    if validAuthHandler(w, r) == true {
        fmt.Fprint(w, "<h1>Welcome to your dashboard!</h1><a href='/logout'>Logout</a>")
    }
}

// ログアウト処理
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

// MyHandler は実際のリクエストを処理するハンドラです。
func MyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Go Middleware!")
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	// "/public/" プレフィックスを外して、中身を返すようにする
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/api/adk_init", handler.AdkInitialHandler)
	http.HandleFunc("/api/adk_run", handler.AdkRunHandler)

    http.HandleFunc("/foo", helloHandler)
    http.HandleFunc("/dashboard", dashboardHandler)
    http.HandleFunc("/goodbye", goodbyeHandler)
    http.HandleFunc("/login", loginPage)
    http.HandleFunc("/logout", logoutHandler)

    fmt.Println("Server running on :8080")
    log.Println("Listening on :8080...")

    http.ListenAndServe(":8080", nil)
}
