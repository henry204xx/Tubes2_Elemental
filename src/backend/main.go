// backend/main.go
package main

import (
	"backend/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/query", api.QueryHandler)
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
