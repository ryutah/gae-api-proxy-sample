package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PING"))
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		echo := r.FormValue("val")
		if echo == "" {
			echo = "HelloWorld!!"
		}
		w.Write([]byte(echo))
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		req := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), 500)
			return
		}
		for k, v := range req {
			log.Printf("%s = %v", k, v)
		}
		fmt.Fprintln(w, `{"result": "OK"}`)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
