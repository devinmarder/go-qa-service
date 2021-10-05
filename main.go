package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/devinmarder/go-qa-service/repository"
)

var repo repository.Repository

type Body struct {
	Payload repository.ServiceCoverage `json:"payload"`
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := repo.UpdateServiceCoverage(body.Payload)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "service name: %v \ncoverage: %v", body.Payload.ServiceName, body.Payload.Coverage)
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Services QA Status</h1> <dl> <hr>")
	serviceCoverage, err := repo.ListServiceCoverage()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, sc := range serviceCoverage {
		fmt.Fprintf(w, "<dt>%v</dt><dd>coverage: %v</dd> <hr>", sc.ServiceName, sc.Coverage)
	}
	fmt.Fprint(w, "<dl>")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	scl, err := repo.ListServiceCoverage()
	if err != nil {
		log.Print()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(scl)
	fmt.Fprint(w, string(b))
}

func main() {
	repository.ConfigureRepository(&repo)
	port := os.Args[1]
	http.HandleFunc("/", updateHandler)
	http.HandleFunc("/stats", webHandler)
	http.HandleFunc("/api/stats", apiHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
