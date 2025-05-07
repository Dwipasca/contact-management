package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var menus = []string{
	"Add Contact",
	"Add Multiple Contact",
	"Edit Contact",
	"Delete Contact",
	"Show List Contact",
	"Search Contact",
	"Export Data To JSON",
	"Import Data From JSON",
	"Export Data To CSV",
	"Import Data From CSV",
}

const (
	SearchByID = "id"
	SearchByName = "name"
	SearchByEmail = "email"
)

type Contact struct {
	ID		int
	Name	string
	Email	string
	Phone	string
}

var contacts []Contact
var NextID = 1

func main() {

	// create a new scanner to read user input from keyboard
	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenus()
		fmt.Print("\nChoose Menu: ")
		
		// wait for user input till the key enter is pressed
		// then the input will stored internaly by scanner
		scanner.Scan()
		
		// get input from scanner
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "1":
			// add contact
			setTitle(menus[0]) 	
			addSingleContact(scanner)

		case "2":
			// add multiple contact
			setTitle(menus[1]) 
			addMultipleContacts(scanner)

		case "3":
			// edit contact
			setTitle(menus[2]) 
			editExistingContact(scanner)

		case "4":
			// delete contact
			setTitle(menus[3]) 
			deleteExistingContact(scanner)

		case "5":
			// show list contact
			setTitle(menus[4]) 
			showListContact(contacts...)
		case "6":
			// search contact
			setTitle(menus[5]) 
			searchExistingContact(scanner)
		case "7":
			// export to json
			setTitle(menus[6]) 
			exportToJSON(scanner)
		case "8":
			// import from json
			setTitle(menus[7]) 
			importFromJSON(scanner)
		case "9":
			// export to csv
			setTitle(menus[8]) 
			exportToCSV(scanner)
		case "10":
			// import from csv
			setTitle(menus[9]) 
			importFromCSV(scanner)
		case "0":
			setTitle("Exit from the program")
			return
		default:
			setTitle("error: Invalid input please input a number")
		}
	
	}


}

// -----------------------------------------------

func addSingleContact(scanner *bufio.Scanner) {
	name := promptRequiredInput(scanner, "Name")
	email := promptRequiredInput(scanner, "Email")
	phone := promptRequiredInput(scanner, "Phone")

	addContact(name, email, phone)
	setRespond("add a new contact", "success")
}

func addMultipleContacts(scanner *bufio.Scanner) {
	var count int
	for {
		countStr := promptInput(scanner, "How many contacts to add")
		num, err := strconv.Atoi(countStr)
		if err != nil || num <= 0 {
			setRespond("invalid input, please try it again", "error")
		} else {
			count = num
			break
		}	
	}

	var newContacts []Contact

	for i := 1; i <= count; i++ {
		fmt.Printf("\n---- New Contact %d ----\n", i)
		name := promptRequiredInput(scanner, "Name")
		email := promptRequiredInput(scanner, "Email")
		phone := promptRequiredInput(scanner, "Phone")
		newContacts = append(newContacts, Contact{
			Name: name,
			Email: email,
			Phone: phone,
		})
	}

	addMultipleContact(newContacts...)
	setRespond("add multiple contact", "success")
}

func editExistingContact(scanner *bufio.Scanner) {
	var id int

	for {
		idStr := promptInput(scanner, "Set the ID you want to edit")
		num, err := strconv.Atoi(idStr)
		if err != nil {
			setRespond("invalid ID input", "error")
		} else {
			id = num
			break
		}
	}

	name := promptInput(scanner, "Name")
	email := promptInput(scanner, "Email")
	phone := promptInput(scanner, "Phone")

	var pName, pEmail, pPhone *string
	if name != "" {
		pName = &name
	}
	if email != "" {
		pEmail = &email
	}
	if phone != "" {
		pPhone = &phone
	}

	if editContact(id, pName, pEmail, pPhone) {
		setRespond("edited contact", "success")
	} else {
		setRespond("Contact not found", "result")
	}
}

func deleteExistingContact(scanner *bufio.Scanner) {
	var id int

	for {
		idStr := promptRequiredInput(scanner, "set the ID you want to delete")
		num, err := strconv.Atoi(idStr)
		if err != nil {
			setRespond("invalid ID input", "error")
		} else {
			id = num
			break
		}
	}

	if deleteContact(id) {
		setRespond("deleted contact", "success")
	} else {
		setRespond("Contact not found", "result")
	}
}

func searchExistingContact(scanner *bufio.Scanner) {
	searchBy := promptRequiredInput(scanner, "Search By (id, name or email)")
	query := promptRequiredInput(scanner, "Query")

	result := searchContact(searchBy, query)
	showListContact(result...)
}

// ---------------------------------------------

func addContact(name, email, phone string) Contact {
	var newContact Contact

	newContact.ID = NextID
	newContact.Name = name
	newContact.Email = email
	newContact.Phone = phone

	contacts = append(contacts, newContact)
	NextID++

	return newContact
}

func addMultipleContact(newContacts ...Contact) {
	for _, nc := range newContacts {
		nc.ID = NextID
		contacts = append(contacts, nc)
		NextID++
	}
}

func editContact(id int, name, email, phone *string) bool {
	for idx, ctc := range contacts {
		if ctc.ID == id {
			if name != nil {
				contacts[idx].Name = *name
			}
			if email != nil {
				contacts[idx].Email = *email
			}
			if phone != nil {
				contacts[idx].Phone = *phone
			}
			return true
		}
	}
	return false
}

func deleteContact(id int) bool {
	for idx, ctc := range contacts {
		if ctc.ID == id {
			contacts = append(contacts[:idx], contacts[idx+1:]... )
			return true
		}
	}
	return false
}

func showListContact(listContact ...Contact) {

	if len(listContact) <= 0 {
		setRespond("data is not found", "result")
	}

	for _, ctc := range listContact {
		fmt.Println("-----------------")
		fmt.Println("ID: ", ctc.ID)
		fmt.Println("Name: ", ctc.Name)
		fmt.Println("Email: ", ctc.Email)
		fmt.Println("Phone: ", ctc.Phone)
	}
}

func searchContact(searchBy, query string) []Contact {	
	var result [] Contact

	searchBy = strings.ToLower(searchBy)
	query = strings.ToLower(query)

	for _, ctc := range contacts {
		switch searchBy {
		case SearchByID:
			if strings.Contains(strconv.Itoa(ctc.ID), query) {
				result = append(result, ctc)
			}
		case SearchByName:
			if strings.Contains(ctc.Name, query) {
				result = append(result, ctc)
			}
		case SearchByEmail:
			if strings.Contains(ctc.Email, query) {
				result = append(result, ctc)
			}
		}
	}
	return result
}

func exportToJSON(scanner *bufio.Scanner) {

	fileName := promptRequiredInput(scanner, "set filename to export to json")

	// create the "data" folder if it does not exist
	// os.ModePerm = 0777 (read/write/execute permissions for all users)
	os.MkdirAll("data", os.ModePerm)

	// join folder and filename into a full file path
	// ex: data/contacts.json
	filePath := filepath.Join("data", fileName)

	// Convert contacts slice into JSON format
	// "" means no prefix, "  " means 2-space indentation
	data, err := json.MarshalIndent(contacts,"", "  ")
	if err != nil {
		setRespond("failed to convert into json","error")
	}

	// Write the JSON data to the specified file with 0644 permissions
	// 0644 = owner can read/write, others can only read
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		setRespond("failed to write file", "error")
	}

	setRespond("contacts successfully saved to JSON", "success")
}

func importFromJSON(scanner *bufio.Scanner) {
	fileName := promptRequiredInput(scanner, "input filename to import from JSON")

	filePath := filepath.Join("data", fileName)

	// read data from the json file
	data, err := os.ReadFile(filePath)
	if err != nil {
		setRespond(err.Error(),"error")
		return
	}

	// create temporary slice to hold data parsed from json
	var dataFromJson []Contact

	// decode json data into temporary slice
	if err := json.Unmarshal(data, &dataFromJson); err != nil {
		setRespond(err.Error(),"error")
		return
	}

	// assign new uniques id to each data from json
	// to avoid duplicate id
	for i := range dataFromJson {
		dataFromJson[i].ID = NextID
		NextID++
	}

	// add new data into contacts slice
	contacts = append(contacts, dataFromJson...)

	setRespond("contacts successfully import from JSON", "success")
}

func exportToCSV(scanner *bufio.Scanner) {
	fileName := promptRequiredInput(scanner, "set file name to export to CSV")
	os.MkdirAll("data",os.ModePerm)
	filePath := filepath.Join("data", fileName)

	// create the file to write csv data
	file, err := os.Create(filePath)
	if err != nil {
		setRespond(err.Error(), "error")
	}
	//ensure file is closed after write the file
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)
	// ensure buffered data is writing to csv file
	defer writer.Flush()

	// write header for each columns
	writer.Write([]string{"ID", "Name", "Email", "Phone"})

	// write datas in csv based on slice contacts
	for _, ctc := range contacts {
		record := []string{
			strconv.Itoa(ctc.ID),
			ctc.Name,
			ctc.Email,
			ctc.Phone,
		}
		writer.Write(record)
	}

	setRespond("contacts successfully saved to CSV", "success")
}

func importFromCSV(scanner *bufio.Scanner) {
	fileName := promptRequiredInput(scanner, "input filename to import from CSV")

	filePath := filepath.Join("data", fileName)

	// open the csv file
	file, err := os.Open(filePath)
	if err != nil {
		setRespond(err.Error(),"error")
		return
	}
	// make sure to file is closed after the function is finished
	defer file.Close() 

	// create new csv reader
	reader := csv.NewReader(file)

	// read all the rows in csv file
	data, err := reader.ReadAll()
	if err != nil {
		setRespond(err.Error(),"error")
		return
	}

	// temporary slices to hold data from csv file
	var dataFromCsv []Contact

	for idx, ctc := range data {

		// because it's a header then we can skip it
		if idx == 0 {
			continue
		}
		
		// insert data from csv to the temporary slice
		// and make sure the id will be unique
		dataFromCsv = append(dataFromCsv, Contact{
			ID: NextID,
			Name: ctc[1],
			Email: ctc[2],
			Phone: ctc[3],
		})

		NextID++
	}

	contacts = append(contacts, dataFromCsv...)
	setRespond("contacts successfully import from CSV", "success")
}

// ----------------------------------------------

func promptInput(scanner *bufio.Scanner, label string) string {
	fmt.Print(label + ": ")
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func promptRequiredInput(scanner *bufio.Scanner, label string) string {
	for {
		input := promptInput(scanner, label)
		if input != "" {
			return input
		}
		setRespond("field "+ label + " must be filled, please try again", "error")
	}
}

func showMenus() {
	setTitle("contact service")
	for idx, mn := range menus {
		setMenu(idx+1, mn)
	}
	setMenu(0, "Exit")
}

func setTitle(text string) {
	fmt.Println()
	fmt.Println("=================")
	fmt.Println(strings.ToUpper(text))
	fmt.Println("=================")
}

func setRespond(text, format string) {
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

func setMenu(num int, menu string) {
	fmt.Printf("%d. %s \n", num, menu)
}