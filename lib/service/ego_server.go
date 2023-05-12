package service

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func RunEgoServer() {
	fmt.Print("ego-server started on port 8080 \n")

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/run", RunHandler)
	// http.Handle("/", http.FileServer(http.Dir("static")))

	http.ListenAndServe(":8080", nil)
}
