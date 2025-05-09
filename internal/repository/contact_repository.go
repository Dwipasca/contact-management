package repository

import "github.com/Dwipasca/contact-management/internal/domain"

type ContactRepository  interface {
	GetAll() ([]domain.Contact, error)
	GetByID(id int) (domain.Contact, error)
	GetByName(name string) ([]domain.Contact, error)
	GetByEmail(email string) (domain.Contact, error)

	Save(contact domain.Contact) error
	SaveAll(contacts []domain.Contact) error
	Update(contact domain.Contact) error
	Delete(id int) error

	ExportToJSON(filename string) error
	ExportToCSV(filename string) error
	ImportFromJSON(filename string) ([]domain.Contact, error)
	ImportFromCSV(filename string) ([]domain.Contact, error)
}