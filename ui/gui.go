package ui

import (
	"fmt"

	"github.com/0xfatty/GoPassVault/pkg/vault"
	"github.com/gotk3/gotk3/gtk"
)

// ShowUI displays the graphical user interface
func ShowUI(v *vault.Vault) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		panic(err)
	}
	win.SetTitle("Go Password Vault")
	win.SetDefaultSize(400, 200)

	grid, err := gtk.GridNew()
	if err != nil {
		panic(err)
	}
	win.Add(grid)

	masterPasswordEntry, _ := gtk.EntryNew()
	masterPasswordEntry.SetPlaceholderText("Master Password")
	grid.Add(masterPasswordEntry)

	serviceEntry, _ := gtk.EntryNew()
	serviceEntry.SetPlaceholderText("Service")
	grid.Add(serviceEntry)

	passwordEntry, _ := gtk.EntryNew()
	passwordEntry.SetPlaceholderText("Password")
	grid.Add(passwordEntry)

	textView, _ := gtk.TextViewNew()
	grid.Attach(textView, 0, 3, 2, 1)

	saveButton, _ := gtk.ButtonNewWithLabel("Save")
	saveButton.Connect("clicked", func() {
		masterPassword, _ := masterPasswordEntry.GetText()
		service, _ := serviceEntry.GetText()
		password, _ := passwordEntry.GetText()

		// Implement validation and saving logic
		// Update textView with result or error message
		text, _ := textView.GetBuffer()
		if v.CheckMasterPassword(masterPassword) {
			v.AddEntry(service, password)
			v.Save()
			text.SetText("Password saved successfully.")
		} else {
			text.SetText("Invalid master password.")
		}
	})
	grid.Attach(saveButton, 0, 4, 1, 1)

	retrieveButton, _ := gtk.ButtonNewWithLabel("Retrieve")
	retrieveButton.Connect("clicked", func() {
		masterPassword, _ := masterPasswordEntry.GetText()
		service, _ := serviceEntry.GetText()

		// Implement validation and retrieval logic
		// Update textView with result or error message
		text, _ := textView.GetBuffer()
		if v.CheckMasterPassword(masterPassword) {
			password, err := v.GetEntry(service)
			if err != nil {
				text.SetText("Error retrieving password.")
			} else {
				text.SetText(fmt.Sprintf("Password for %s is %s.", service, password))
			}
		} else {
			text.SetText("Invalid master password.")
		}
	})
	grid.Attach(retrieveButton, 1, 4, 1, 1)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.ShowAll()
	gtk.Main()
}
