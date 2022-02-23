package main

type SomeBullShitServiceContract interface {
}

type SomeBullShitService struct {
	Name string
}

func NewSomeBullShitService(name string) *SomeBullShitService {
	return &SomeBullShitService{Name: name}
}
