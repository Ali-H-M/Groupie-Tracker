package main

import (
	"fmt"
	"groupie-tracker/funcs"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", funcs.Handler)

	fmt.Println("Server running at http://localhost:8080")
	fmt.Println("Press (crtl + c) to stop the program")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
