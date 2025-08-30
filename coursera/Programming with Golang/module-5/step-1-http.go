package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Hello World")
	value := r.PathValue("name")
	fmt.Println(value)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/hello/", helloHandler)
	fmt.Println("Listening on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
