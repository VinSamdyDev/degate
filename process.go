package main

type DegateProcess map[int]*DegateService

// NewDegateProcess creates a new DegateProcess instance
func NewDegateProcess() *DegateProcess {
	return &DegateProcess{}
}
