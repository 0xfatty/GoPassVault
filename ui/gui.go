package ui

import (
	"fmt"
	"log"

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

	// Check if master password is set
	if !v.IsMasterPasswordSet() {
		dialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK_CANCEL, "Set your master password")
		dialog.SetTitle("First-time Setup")

		entry, err := gtk.EntryNew()
		if err != nil {
			log.Fatal("Unable to create the entry:", err)
		}
		entry.SetVisibility(false)

		contentArea, err := dialog.GetContentArea()
		if err != nil {
			log.Fatal("Unable to get content area:", err)
		}

		contentArea.Add(entry)
		dialog.ShowAll()

		response := dialog.Run()

		if response == gtk.RESPONSE_OK {
			masterPassword, err := entry.GetText()
			if err != nil {
				log.Fatal("Unable to get text:", err)
			}
			v.SetMasterPassword(masterPassword)
		}
		dialog.Destroy()

		if response != gtk.RESPONSE_OK {
			return
		}
	}

	masterPasswordEntry, _ := gtk.EntryNew()
	masterPasswordEntry.SetPlaceholderText("Master Password")
	grid.Attach(masterPasswordEntry, 0, 0, 2, 1)

	serviceEntry, _ := gtk.EntryNew()
	serviceEntry.SetPlaceholderText("Service")
	grid.Attach(serviceEntry, 0, 1, 1, 1)

	passwordEntry, _ := gtk.EntryNew()
	passwordEntry.SetPlaceholderText("Password")
	grid.Attach(passwordEntry, 1, 1, 1, 1)

	addButton, _ := gtk.ButtonNewWithLabel("Add")
	grid.Attach(addButton, 0, 2, 1, 1)

	getButton, _ := gtk.ButtonNewWithLabel("Get")
	grid.Attach(getButton, 1, 2, 1, 1)

	addButton.Connect("clicked", func() {
		masterPassword, _ := masterPasswordEntry.GetText()
		if !v.CheckMasterPassword(masterPassword) {
			fmt.Println("Master password incorrect")
			return
		}

		service, _ := serviceEntry.GetText()
		password, _ := passwordEntry.GetText()

		v.AddEntry(service, password)
		v.Save()
	})

	getButton.Connect("clicked", func() {
		masterPassword, _ := masterPasswordEntry.GetText()
		if !v.CheckMasterPassword(masterPassword) {
			fmt.Println("Master password incorrect")
			return
		}

		service, _ := serviceEntry.GetText()
		password, err := v.GetEntry(service)
		if err != nil {
			fmt.Println("Service not found")
			return
		}

		passwordEntry.SetText(password)
	})

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.ShowAll()
	gtk.Main()
}
