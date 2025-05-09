package ui

import (
	"fmt"
	"strings"

	"github.com/Dwipasca/contact-management/internal/domain"
)

var Menus = []string{
	"Add Contact",
	"Add Multiple Contact",
	"Edit Contact",
	"Delete Contact",
	"Show List Contact",
	"Search Contact",
	"Export Contacts",
	"Import Contacts",
}

func PrintMenu() {
	SetTitle("contact service")
	for idx, mn := range Menus {
		SetMenu(idx+1, mn)
	}
	SetMenu(0, "Exit")
}

func PrintContacts(contacts ...domain.Contact) {
	fmt.Println("\n-- Contact List --")
	for _, ctc := range contacts {
		fmt.Println("ID: ", ctc.ID)
		fmt.Println("Name: ", ctc.Name)
		fmt.Println("Email: ", ctc.Email)
		fmt.Println("Phone: ", ctc.Phone)
	}
}

func SetTitle(text string) {
	fmt.Println()
	fmt.Println("=================")
	fmt.Println(strings.ToUpper(text))
	fmt.Println("=================")
}

func SetRespond(text, format string) {
	fmt.Println("----------------------")
	switch format {
	case "error":	
		fmt.Println("ERROR: "+text)
	case "success":
		fmt.Println("SUCCESS: "+text)
	case "result":
		fmt.Println("RESULT: "+text)
	default:
		fmt.Println("Invalid format, please type success, error or result")
	}
	fmt.Println("----------------------")
}

func SetMenu(num int, menu string) {
	fmt.Printf("%d. %s \n", num, menu)
}