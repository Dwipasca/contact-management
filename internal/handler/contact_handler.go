package handler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Dwipasca/contact-management/internal/domain"
	"github.com/Dwipasca/contact-management/internal/usecase"
	"github.com/Dwipasca/contact-management/ui"
)


type ContactHandler struct {
	scanner *bufio.Scanner
	service *usecase.ContactService
}

func NewContactHandler(service *usecase.ContactService) *ContactHandler{
	return &ContactHandler{
		scanner: bufio.NewScanner(os.Stdin),
		service: service,
	}
}

func (ch *ContactHandler) ShowMainMenu() {
	for {
		ui.PrintMenu()
		fmt.Print("\nSelect option: ")
		ch.scanner.Scan()
		choice := strings.TrimSpace(ch.scanner.Text())
		ui.ClearScreen()

		switch choice {
		case "1":
			ch.handleAddContact()
		case "2":
			ch.handleAddMultipleContact()
		case "3":
			ch.handleEditContact()
		case "4":
			ch.handleDeleteContact()
		case "5":
			ch.handleListContacts()
		case "6":
			ch.handleSearchContact()
		case "7":
			ch.handleExportContacts()
		case "8":
			ch.handleImportContacts()
		case "0":
			fmt.Println("Exiting application...")
			return
		default:
			ui.SetRespond("Invalid input, please enter a number between 0-8", "error")
		}
	}
}

func (ch *ContactHandler) handleAddContact() {
	ui.SetTitle(ui.Menus[0])

	name := ui.PromptRequiredInput(ch.scanner, "Name")
	email := ui.PromptRequiredInput(ch.scanner, "Email")
	phone := ui.PromptInput(ch.scanner, "Phone")

	err := ch.service.AddContact(name, email, phone)
	if err != nil {
		
		switch {
		case errors.Is(err, usecase.ErrNameRequired),
			errors.Is(err, usecase.ErrEmailRequired),
			errors.Is(err, usecase.ErrInvalidEmail),
			errors.Is(err, usecase.ErrEmailAlreadyExist):
			ui.SetRespond(err.Error(), "error")
		default:
			ui.SetRespond("Something went wrong: "+err.Error(), "error")
		}
		return
	} 

	ui.SetRespond("Successfully added new contact", "success")
}

func (ch *ContactHandler) handleAddMultipleContact() {
	ui.SetTitle(ui.Menus[1])

	var count int
	for {
		countStr := ui.PromptInput(ch.scanner, "How many contacts to add")
		num, err := strconv.Atoi(countStr)
		if err != nil || num <= 0 {
			ui.SetRespond("invalid input, please try it again", "error")
		} else {
			count = num
			break
		}	
	}

	var newContacts []domain.Contact

	for i := 1; i <= count; i++ {
		fmt.Printf("\n---- New Contact %d ----\n", i)
		name := ui.PromptRequiredInput(ch.scanner, "Name")
		email := ui.PromptRequiredInput(ch.scanner, "Email")
		phone := ui.PromptRequiredInput(ch.scanner, "Phone")

		newContacts = append(newContacts, domain.Contact{
			Name: name,
			Email: email,
			Phone: phone,
		})
	}

	err := ch.service.AddMultipleContact(newContacts)
	if err != nil {
		ui.SetRespond(err.Error(), "error")
		return
	}

	ui.SetRespond("Successfully added all contacts", "success")
}

func (ch *ContactHandler) handleEditContact() {
	ui.SetTitle(ui.Menus[2])
	
	// Get contact ID to edit
	idStr := ui.PromptRequiredInput(ch.scanner, "Enter contact ID to edit: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ui.SetRespond("Invalid ID format. Please enter a number", "error")
		return
	}
	
	// Get contact by ID
	contact, err := ch.service.SearchByID(id)
	if err != nil {
		ui.SetRespond("Contact not found", "error")
		return
	}
	
	fmt.Println("\n-- Editing Contact --")
	fmt.Println("(leave empty to keep current)")
	
	fmt.Println("\nCurrent Email:", contact.Name)
	name := ui.PromptInput(ch.scanner, "New Name")
	if name == "" {
		name = contact.Name
	}
	
	fmt.Println("Current Email:", contact.Email)
	email := ui.PromptInput(ch.scanner, "New Email")
	if email == "" {
		email = contact.Email
	}
	
	fmt.Println("Current Phone:", contact.Phone)
	phone := ui.PromptInput(ch.scanner, "New Phone")
	if phone == "" {
		phone = contact.Phone
	}
	
	// Update contact
	err = ch.service.EditContact(id, name, email, phone)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNameRequired),
			errors.Is(err, usecase.ErrEmailAlreadyExist),
			errors.Is(err, usecase.ErrInvalidEmail):
			ui.SetRespond(err.Error(), "error")
		default:
			ui.SetRespond("Failed to update contact: "+err.Error(), "error")
		}
		return
	} 
	ui.SetRespond("Contact updated successfully", "success")
}

func (ch *ContactHandler) handleDeleteContact() {
	ui.SetTitle(ui.Menus[3])
	
	// Get all contacts first to check if there are any
	contacts, err := ch.service.GetAllContacts()
	if err != nil {
		ui.SetRespond("Contact not found", "error")
		return
	}
	
	if len(contacts) == 0 {
		ui.SetRespond("No contacts available to delete", "error")
		return
	}
	
	// Display available contacts
	ui.PrintContacts(contacts...)
	
	// Get contact ID to delete
	idStr := ui.PromptRequiredInput(ch.scanner, "\nEnter contact ID to delete")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ui.SetRespond("Invalid ID format, please enter a number", "error")
		return
	}
	
	// Confirm deletion
	confirm := ui.PromptRequiredInput(ch.scanner, "Are you sure you want to delete this contact? (y/n)")
	if strings.ToLower(confirm) != "y" {
		ui.SetRespond("Deletion cancelled", "result")
		return
	}
	
	// Delete contact
	err = ch.service.DeleteContact(id)
	if err != nil {
		ui.SetRespond(err.Error(), "error")
		return
	}

	ui.SetRespond("Contact deleted successfully", "success")
}

func (ch *ContactHandler) handleListContacts() {
	ui.SetTitle(ui.Menus[4])
	
	contacts, err := ch.service.GetAllContacts()
	if err != nil {
		if errors.Is(err, usecase.ErrNoContacts) {
			ui.SetRespond("No contacts available","result")
		} else {
			ui.SetRespond("something went wrong: "+ err.Error(),"error")
		}
		return
	}
	
	ui.PrintContacts(contacts...)
}

func (ch *ContactHandler) handleSearchContact() {
	ui.SetTitle(ui.Menus[5])
	
	fmt.Println("Search by:")
	fmt.Println("1. ID")
	fmt.Println("2. Name")
	fmt.Println("3. Email")
	
	choice := ui.PromptRequiredInput(ch.scanner, "Select option: ")
	
	switch choice {
	case "1":
		idStr := ui.PromptRequiredInput(ch.scanner, "Enter ID: ")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ui.SetRespond("Invalid ID format. Please enter a number", "error")
			return
		}
		
		contact, err := ch.service.SearchByID(id)
		if err != nil {
			if errors.Is(err, usecase.ErrNoContacts) {
				ui.SetRespond("contact with id "+ idStr +" is not found","result")
			}else {
				ui.SetRespond("something went wrong: "+ err.Error(), "error")
			}
			return
		}
		
		ui.PrintContacts(contact)
		
	case "2":
		name := ui.PromptRequiredInput(ch.scanner, "Enter Name: ")
		contacts, err := ch.service.SearchByName(name)
		if err != nil {
			if errors.Is(err, usecase.ErrNoContacts) {
				ui.SetRespond("Contacts with name "+name+" is not found", "result")
			}else {
				ui.SetRespond("something went wrong: "+ err.Error(), "error")
			}
			return
		}
		
		ui.PrintContacts(contacts...)
		
	case "3":
		email := ui.PromptRequiredInput(ch.scanner, "Enter Email: ")
		contact, err := ch.service.SearchByEmail(email)
		if err != nil {
			if errors.Is(err, usecase.ErrNoContacts) {
				ui.SetRespond("contact with email "+ email +" is not found","result")
			}else {
				ui.SetRespond("something went wrong: "+ err.Error(), "error")
			}
			return
		}
		
		ui.PrintContacts(contact)
		
	default:
		ui.SetRespond("Invalid option, please enter a number between 1-3 ", "error")
	}
}

func (ch *ContactHandler) handleExportContacts() {
	ui.SetTitle(ui.Menus[6])
	
	fmt.Println("Export format:")
	fmt.Println("1. JSON")
	fmt.Println("2. CSV")
	
	choice := ui.PromptRequiredInput(ch.scanner, "\nSelect option")
	filename := ui.PromptRequiredInput(ch.scanner, "Enter filename (without extension)")
	
	var err error
	switch choice {
	case "1":
		err = ch.service.ExportToJSON(filename + ".json")
	case "2":
		err = ch.service.ExportToCSV(filename + ".csv")
	default:
		ui.SetRespond("Invalid option", "error")
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidExportFilename):
			ui.SetRespond("Invalid filename, please avoid special characters.", "error")
		default:
			ui.SetRespond("Export failed: "+err.Error(), "error")
		}
		return
	}
	
	ui.SetRespond("Contacts exported successfully to "+filename, "success")
}

func (ch *ContactHandler) handleImportContacts() {
	ui.SetTitle(ui.Menus[7])
	
	fmt.Println("Import format:")
	fmt.Println("1. JSON")
	fmt.Println("2. CSV")
	
	choice := ui.PromptRequiredInput(ch.scanner, "\nSelect option")
	filename := ui.PromptRequiredInput(ch.scanner, "Enter filename (with extension)")
	
	var contacts []domain.Contact
	var err error
	
	switch choice {
	case "1":
		contacts, err = ch.service.ImportFromJSON(filename)
	case "2":
		contacts, err = ch.service.ImportFromCSV(filename)
	default:
		ui.SetRespond("Invalid option", "error")
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidImportFilename):
			ui.SetRespond("Invalid filename, must end with .json or .csv", "error")
		default:
			ui.SetRespond("Import failed: "+err.Error(), "error")
		}
		return
	}
	
	ui.SetRespond(fmt.Sprintf("Successfully imported %d contacts", len(contacts)), "success")
}
