package repository

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Dwipasca/contact-management/internal/domain"
)

type ContactRepositoryImpl struct {
	contacts	[]domain.Contact
	nextID		int
}

func NewContactRepository() *ContactRepositoryImpl {
	return &ContactRepositoryImpl{
		contacts: []domain.Contact{},
		nextID: 1,
	}
}

func (cr *ContactRepositoryImpl) GetAll() ([]domain.Contact, error) {
	return cr.contacts, nil
}

func (cr *ContactRepositoryImpl) GetByID(id int) (domain.Contact, error) {
	for _, ctc := range cr.contacts {
		if ctc.ID == id{
			return ctc, nil
		}
	}

	return domain.Contact{}, nil
}

func (cr *ContactRepositoryImpl) GetByName(name string) ([]domain.Contact, error) {
	var result []domain.Contact
	for _, ctc := range cr.contacts{
		if ctc.Name == name {
			result = append(result, ctc)
		}
	}

	return result, nil
}

func (cr *ContactRepositoryImpl) GetByEmail(email string) (domain.Contact, error) {
	for _, ctc := range cr.contacts {
		if ctc.Email == email{
			return ctc, nil
		}
	}

	return domain.Contact{}, nil
}

func (cr *ContactRepositoryImpl) Save(contact domain.Contact) error {
	contact.ID = cr.nextID
	cr.contacts = append(cr.contacts, contact)
	cr.nextID++
	return nil
}

func (cr *ContactRepositoryImpl) SaveAll(contacts []domain.Contact) error {
	for _, ctc := range contacts {
		if err := cr.Save(ctc); err != nil {
			return fmt.Errorf("failed to save contact %s: %w", ctc.Name, err)
		}
	}
	return nil
}

func (cr *ContactRepositoryImpl) findIndexByID(id int) int {
	for idx, ctc := range cr.contacts {
		if ctc.ID == id {
			return idx
		}
	}
	return -1 // not found
}

func (cr *ContactRepositoryImpl) Update(updated domain.Contact) error {
	idx := cr.findIndexByID(updated.ID)
	if idx == -1 {
		return errors.New("contact is not found")
	}
	cr.contacts[idx] = updated
	return nil
}

func (cr *ContactRepositoryImpl) Delete(id int) error {
	idx := cr.findIndexByID(id)
	if idx == -1 {
		return errors.New("contact is not found")
	}
	cr.contacts = append(cr.contacts[:idx], cr.contacts[idx+1:]... )
	return nil
}

func (cr *ContactRepositoryImpl) ExportToJSON(filename string) error {

	// Convert contacts slice into JSON format
	// "" means no prefix, "  " means 2-space indentation
	data, err := json.MarshalIndent(cr.contacts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal contacts to JSON: %w", err)
	}

	// Write the JSON data to the specified file with 0644 permissions
	// 0644 = owner can read/write, others can only read
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filename, err)
	}

	return nil
}

func (cr *ContactRepositoryImpl) ExportToCSV(filename string) error {
	// create the file to write csv data
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file contacts CSV: %w", err)
	}
	//ensure file is closed after write the file
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	if err := writer.Write([]string{"ID", "Name", "Email", "Phone"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// write data rows
	for _, ctc := range cr.contacts {
		record := []string{
			strconv.Itoa(ctc.ID),
			ctc.Name,
			ctc.Email,
			ctc.Phone,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record for ID %d: %w", ctc.ID, err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return nil
}

func (cr *ContactRepositoryImpl) ImportFromJSON(filename string) ([]domain.Contact, error) {
	
	// read data from json file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("file %s is empty", filename)
	}

	// temporary slice
	var dataFromJSON []domain.Contact
	// decode json data into temporary slice
	if err := json.Unmarshal(data, &dataFromJSON); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	cr.SaveAll(dataFromJSON)

	return dataFromJSON, nil
}

func (cr *ContactRepositoryImpl) ImportFromCSV(filename string) ([]domain.Contact, error) {
	// open the csv file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// make sure to file is closed after the function is finished
	defer file.Close()

	// create new csv reader
	reader := csv.NewReader(file)
	// read all the rows in csv file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var dataFromCSV []domain.Contact
	for idx, dt := range records {
		
		// because it's a header then we can skip it
		if idx == 0 {
			continue
		}

		// insert data from csv to the temporary slice
		dataFromCSV = append(dataFromCSV, domain.Contact{
			// and make sure the id will be unique
			ID: cr.nextID,
			Name: dt[1],
			Email: dt[2],
			Phone: dt[3],
		})
	}

	cr.SaveAll(dataFromCSV)
	return dataFromCSV, nil
}