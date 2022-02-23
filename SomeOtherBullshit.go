package main

type SomeOtherBullshitServiceContract interface {
}

type SomeOtherBullshitService struct {
	Name string
}

func NewSomeOtherBullshitService(name string) *SomeOtherBullshitService {
	return &SomeOtherBullshitService{Name: name}
}
