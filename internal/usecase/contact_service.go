package usecase

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Dwipasca/contact-management/internal/domain"
	"github.com/Dwipasca/contact-management/internal/repository"
)

type ContactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) *ContactService {
	return &ContactService{
		repo: repo,
	}
}

var (
	ErrNoContacts		 = errors.New("no contacts found")
	ErrNameRequired      = errors.New("name is required")
	ErrEmailRequired     = errors.New("email is required")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrInvalidExportFilename = errors.New("invalid export filename")
	ErrInvalidImportFilename = errors.New("invalid import filename")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (cs *ContactService) GetAllContacts() ([]domain.Contact, error) {
	contacts, err := cs.repo.GetAll()
	if err != nil {
		return  nil, fmt.Errorf("failed to retrieve contacts: %w", err)
	}	

	if len(contacts) == 0 {
		return nil, ErrNoContacts
	}

	return contacts, nil
}

func (cs *ContactService) SearchByID(id int) (domain.Contact, error){
	contact, err := cs.repo.GetByID(id)
	if err != nil {
		return domain.Contact{}, err // file corrupt or something
	}

	if contact.ID == 0 {
		return domain.Contact{}, ErrNoContacts
	}

	return contact, nil
}

func (cs *ContactService) SearchByName(name string) ([]domain.Contact, error) {
	contacts, err := cs.repo.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contacts: %w", err)
	}

	if len(contacts) == 0 {
		return nil, ErrNoContacts
	}

	return contacts, nil
}

func (cs *ContactService) SearchByEmail(email string) (domain.Contact, error){
	contact, err := cs.repo.GetByEmail(email)
	if err != nil {
		return domain.Contact{}, err // file corrupt or something
	}

	if contact.ID == 0 {
		return domain.Contact{}, ErrNoContacts
	}

	return contact, nil

}

func (cs *ContactService) AddContact(name, email, phone string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return ErrNameRequired
	}

	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	if email == "" {
		return ErrEmailRequired
	}

	// check if email is already exists or not
	existing, err := cs.repo.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to check existing email: %w", err)
	}

	if existing.ID != 0 {
		return ErrEmailAlreadyExist
	}

	newContact := domain.Contact{
		Name: name,
		Email: email,
		Phone: phone,
	}


	if err := cs.repo.Save(newContact); err != nil {
		return fmt.Errorf("failed to save contact: %w", err)
	}

	return nil
}

func (cs *ContactService) AddMultipleContact(newContacts []domain.Contact) error {
	var failed []string

	for _, ctc := range newContacts{
		if err := cs.AddContact(ctc.Name, ctc.Email, ctc.Phone); err != nil {
			failed = append(failed, fmt.Sprintf("%s (%s): %v", ctc.Name, ctc.Email, err))
		}
	}

	if len(failed) > 0 {
		return fmt.Errorf("some contacts failed to add:\n%s", strings.Join(failed, "\n"))
	}

	return nil
}

func (cs *ContactService) EditContact(id int, name, email, phone string) error {
	prevContact, err := cs.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("fetching contact failed: %w", err)
	}
	
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return ErrNameRequired
	}

	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	if email != prevContact.Email {
		existing, err := cs.repo.GetByEmail(email)
		if err == nil && existing.ID != id {
			return ErrEmailAlreadyExist
		}
	}

	updated := domain.Contact{
		ID: id,
		Name: name,
		Email: email,
		Phone: phone,
	}

	if err := cs.repo.Update(updated); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

func (cs *ContactService) DeleteContact(id int) error {

	if err := cs.repo.Delete(id); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	return nil
}

func (cs *ContactService) ExportToJSON(filename string) error {
	
	if strings.TrimSpace(filename) == "" || strings.Contains(filename, "..") {
		return ErrInvalidExportFilename
	}

	// create the "data" folder if it does not exist
	// os.ModePerm = 0777 (read/write/execute permissions for all users)
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data folder: %w", err)
	}

	// join folder and filename into a full file path
	// ex: data/contacts.json
	filePath := filepath.Join("data", filename)

	if err := cs.repo.ExportToJSON(filePath); err != nil {
		return fmt.Errorf("failed to export contacts to JSON: %w", err)
	}

	return nil
}

func (cs *ContactService) ExportToCSV(filename string) error {

	if strings.TrimSpace(filename) == "" || strings.Contains(filename, "..") {
		return ErrInvalidExportFilename
	}
	
	// create the "data" folder if it does not exist
	// os.ModePerm = 0777 (read/write/execute permissions for all users)
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data folder: %w", err)
	}

	// join folder and filename into a full file path
	// ex: data/contacts.json
	filePath := filepath.Join("data", filename)

	if err := cs.repo.ExportToCSV(filePath); err != nil {
		return fmt.Errorf("failed to export contacts to CSV: %w", err)
	}

	return nil
}

func (cs *ContactService) ImportFromJSON(filename string) ([]domain.Contact,error) {
	filename = strings.TrimSpace(filename)

	if filename == "" || !strings.HasSuffix(filename, ".json") || strings.Contains(filename, "..") {
		return nil, ErrInvalidImportFilename
	}

	filePath := filepath.Join("data", filename)

	contacts, err := cs.repo.ImportFromJSON(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to import JSON contacts: %w", err)
	}

	return contacts, nil
}

func (cs *ContactService) ImportFromCSV(filename string) ([]domain.Contact,error) {
	return cs.repo.ImportFromCSV(filename)
}

