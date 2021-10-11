package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/devinmarder/go-qa-service/event"
	"github.com/devinmarder/go-qa-service/repository"
)

// repo is the repository interface used by the handelers.
var repo repository.Repository

// Body represents the top-level fields of the request.
type Body struct {
	Payload repository.ServiceCoverage `json:"payload"`
}

// updateHandler is used for updating the repo with the coverage statistics provided in the request.
// The request body must be json formatted with fields compatible with the Body struct.
func updateHandler(w http.ResponseWriter, r *http.Request) {
	var body Body

	// Decodes reqest body.
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Updates repository.
	err := repo.UpdateServiceCoverage(body.Payload)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "service name: %v \ncoverage: %v", body.Payload.ServiceName, body.Payload.Coverage)
}

// webHandler writes an html formatted response of the items contained in the repository.
func webHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Services QA Status</h1> <dl> <hr>")
	// Get list of service coverage opjects.
	serviceCoverage, err := repo.ListServiceCoverage()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Generate HTML response.
	for _, sc := range serviceCoverage {
		fmt.Fprintf(w, "<dt>%v</dt><dd>coverage: %v</dd> <hr>", sc.ServiceName, sc.Coverage)
	}
	fmt.Fprint(w, "<dl>")
}

// apiHandler write a json formatted respose of the items contained in the repository.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	scl, err := repo.ListServiceCoverage()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Generate Json response.
	b, _ := json.Marshal(scl)
	fmt.Fprint(w, string(b))
}

// attachEvent wraps a provided handler function and writes its body to the provided channel.
func attachEvent(fn func(http.ResponseWriter, *http.Request), eventChan chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		// Replace body with tee to write to make stream available for event.
		r.Body = io.NopCloser(tee)
		fn(w, r)
		// Writes request to event to channel
		eventChan <- buf.String()
	}
}

func main() {
	// Configure the repository.
	repository.ConfigureRepository(&repo)

	// Create channel for writing events and start gorutine to produce events written to channel.
	eventChan := make(chan string)
	go event.RunEventProducer(eventChan)

	// Start a test listener to log events.
	go event.RunEventlistener()

	// get port from args.
	port := os.Args[1]

	// Add handlers to the default HTTP ServeMux.
	http.HandleFunc("/", attachEvent(updateHandler, eventChan))
	http.HandleFunc("/stats", webHandler)
	http.HandleFunc("/api/stats", apiHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
