package main

import "net/http"

type DegateServiceResponse struct {
	ResWriter http.ResponseWriter
	Req       *http.Request
	Err       error
}

// Error implements error
func (DegateServiceResponse) Error() string {
	panic("unimplemented")
}

type DegateService struct {
	name    string
	handler func(http.ResponseWriter, *http.Request) DegateServiceResponse
}

// NewDegateService creates a new DegateService instance
func NewDegateService(name string, start func(http.ResponseWriter, *http.Request) DegateServiceResponse) *DegateService {
	return &DegateService{
		name:    name,
		handler: start,
	}
}

// NewDegateServiceResponse creates a new DegateServiceResponse instance
func NewDegateServiceResponse(resWriter http.ResponseWriter, req *http.Request, err error) DegateServiceResponse {
	return DegateServiceResponse{
		ResWriter: resWriter,
		Req:       req,
		Err:       err,
	}
}
