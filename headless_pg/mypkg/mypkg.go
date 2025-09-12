package mypkg

import(
	"fmt"
	"net/http"
)

func Test() {
 fmt.Println("mypkg.TestHandler")
}


func Test2(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Test2.Hello!")
}
