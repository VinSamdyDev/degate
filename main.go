package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// TestService1 is a test service
var TestService1 = NewDegateService(
	"TestService1",
	func(w http.ResponseWriter, r *http.Request) DegateServiceResponse {

		// await for 3 seconds
		time.Sleep(3 * time.Second)
		return NewDegateServiceResponse(w, r, nil)
	},
)

// TestService2 is a test service
var TestService2 = NewDegateService(
	"TestService2",
	func(w http.ResponseWriter, r *http.Request) DegateServiceResponse {

		// await for 3 seconds
		time.Sleep(3 * time.Second)
		return NewDegateServiceResponse(w, r, nil)
	},
)

// TestService3 is a test service
var TestService3 = NewDegateService(
	"TestService3",
	func(w http.ResponseWriter, r *http.Request) DegateServiceResponse {

		// set the header token to the response
		r.Header.Set("Step", "Service 3")

		// await for 3 seconds
		time.Sleep(3 * time.Second)

		// return the response
		res := map[string]interface{}{
			"message": "Hello World",
		}
		w.Header().Set("Content-Type", "application/json")
		// convert map to json
		jsonStr, err := json.Marshal(res)

		if err != nil {
			return NewDegateServiceResponse(w, r, err)
		}

		// write the response
		w.Write(jsonStr)

		return NewDegateServiceResponse(w, r, nil)
	},
)

// TestProcess is a test process
var TestProcess = DegateProcess{
	1: TestService1,
	2: TestService2,
	3: TestService3,
}

// TestGateway is a test gateway
var TestGateway = NewDegateway(8080, []string{"*"})

// Test sets the processes to the gateway and starts the gateway
func Test() {
	TestGateway.SetProcesses(map[string]*DegateProcess{
		"/api/v1/test": &TestProcess,
	})

	TestGateway.Start()
}

func main() {
	Test()
}
