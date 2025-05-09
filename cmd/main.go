package main

import (
	"github.com/Dwipasca/contact-management/internal/handler"
	"github.com/Dwipasca/contact-management/internal/repository"
	"github.com/Dwipasca/contact-management/internal/usecase"
)

func main() {
	repo := repository.NewContactRepository()
	service := usecase.NewContactService(repo)
	handler := handler.NewContactHandler(service)

	handler.ShowMainMenu()
}