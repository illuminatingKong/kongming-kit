package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message  string              `json:"message"`
	DataCode int64               `json:"data-code"`
	Data     []map[string]string `json:"data"`
	Page     int                 `json:"page"`
	Limt     int                 `json:"limit"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "user data", DataCode: 200, Data: []map[string]string{{"name": "John", "age": "30"}}, Page: 1, Limt: 10}

	// 将结构体转换为JSON格式
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/api/users", handleRequest)

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
