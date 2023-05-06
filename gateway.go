package main

import (
	"fmt"
	"log"
	"net/http"
)

type DegateWay struct {
	port      int
	cors      []string
	processes map[string]*DegateProcess
}

// NewDegateway creates a new Degateway instance
func NewDegateway(port int, cors []string) *DegateWay {
	// check if cors is empty
	if len(cors) == 0 {
		cors = []string{"*"}
	}

	return &DegateWay{
		port:      port,
		cors:      cors,
		processes: make(map[string]*DegateProcess),
	}
}

// Set the DegateProcess instances to the Degateway instance
func (dgw *DegateWay) SetProcesses(processes map[string]*DegateProcess) {
	dgw.processes = processes
}

// Define a middleware function to wrap the HTTP handler with CORS headers
func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Access-Control-Allow-Origin header to allow any domain to access the resource
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Check if the HTTP method is allowed
		if r.Method == "OPTIONS" {
			// Set the Access-Control-Allow-Methods header to allow any method
			w.Header().Set("Access-Control-Allow-Methods", "*")

			// Set the Access-Control-Allow-Headers header to specify which headers are allowed
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			// Set the status code to 204 No Content and return
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Start the Degateway instance
func (dgw *DegateWay) Start() {
	// Create a router to handle the endpoints
	router := http.NewServeMux()

	// loop through the processes
	for path, process := range dgw.processes {
		router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			// loop through the DegateService instances order by the key
			var (
				serviceRes = NewDegateServiceResponse(w, r, nil)
			)
			for _, service := range *process {
				// wait for the service to complete and return the response to the client
				fmt.Printf("Processing service %s...\n", service.name)
				serviceRes = service.handler(serviceRes.ResWriter, serviceRes.Req)
				if serviceRes.Err != nil {
					// if error occurs, break the loop
					fmt.Printf("Service %s error: %s\n", service.name, serviceRes.Err.Error())
					break
				}

				fmt.Printf("Service %s completed\n", service.name)
			}
		})
	}

	// Wrap the handler function with the CORS middleware
	corsHandler := allowCORS(router)

	// start the server
	servePort := fmt.Sprintf(":%d", dgw.port)
	log.Fatal(http.ListenAndServe(servePort, corsHandler))
}
