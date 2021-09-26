package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/devinmarder/go-qa-service/repository"
)

var repo repository.LocalRepository

type Body struct {
	Payload repository.ServiceCoverage `json:"payload"`
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	repo.UpdateServiceCoverage(body.Payload.ServiceName, body.Payload.Coverage)
	fmt.Fprintf(w, "service name: %v \ncoverage: %v", body.Payload.ServiceName, body.Payload.Coverage)
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	//<generate html formatted response>
	fmt.Fprintf(w, "<h1>list of services and their coverage</h1>")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(repo.ListServiceCoverage())
	fmt.Fprint(w, string(b))
}

func main() {
	port := os.Args[1]
	http.HandleFunc("/", updateHandler)
	http.HandleFunc("/stats", webHandler)
	http.HandleFunc("/api/stats", apiHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
