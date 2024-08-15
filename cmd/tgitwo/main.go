package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
)

var port = flag.Int("port", 9200, "specify the port")

func init() {
	flag.Parse()
}

func main() {
	os.Exit(mainerr())
}

func mainerr() int {
	err := initializeServer()
	if err != nil {
		log.Fatal(err.Error())
		return 1
	}

	return 0
}

func initializeServer() error {
	router := chi.NewMux()
	portStr := fmt.Sprintf(":%d", *port)
	srv := http.Server{
		Addr:    portStr,
		Handler: router,
	}

	router.Get("/text/generate/ping", http.HandlerFunc(pingHandler))
	router.Post("/text/generate", http.HandlerFunc(textGenerator))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	wg.Wait()

	return nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Create a map to hold the response data
	response := map[string]string{"response": "pong"}

	// Encode the response map to JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}

func textGenerator(w http.ResponseWriter, r *http.Request) {
	var responseData = struct {
		Result string `json:"result"`
	}{
		Result: "Welcome to the text generator two. \n Text generator two APIs are powerful tools that enhance productivity, creativity, and efficiency across various domains. \n They enable businesses, developers, educators, and creatives to automate tasks, generate personalized content, and provide instant responses, all while saving time and resources. \n As AI continues to advance, the demand for text generator APIs is likely to grow, driving further innovation and integration into everyday applications.",
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
