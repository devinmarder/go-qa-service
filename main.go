package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	//<extract service and coverage stats>
	//<write service and coverage stats to repository>
	fmt.Fprintf(w, "updated service")
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	//<get list of services and their coverages from repository>
	//<generate html formatted response>
	fmt.Fprintf(w, "<h1>list of services and their coverage</h1>")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	//<get list of services and their coverages from repository>
	//<generate json encoded response>
	fmt.Fprint(w, `{"services": [{"service_name": "example_name", "coverage": 23}]`)
}

func main() {
	port := os.Args[1]
	http.HandleFunc("/", updateHandler)
	http.HandleFunc("/stats", webHandler)
	http.HandleFunc("/api/stats", apiHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
