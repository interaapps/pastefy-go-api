# Pastefy Go API Client

```go
package main

import pastefy "github.com/interaapps/pastefy-go-api"

func main() {

	client := pastefy.NewClient()

	// Not required for just fetching pastes or folders
	client.SetApiToken("...")

	createPaste := pastefy.Paste{
		Title:     "test.js",
		Content:   `console.log("Hey");`,
	}
	createdPaste, _ := client.CreatePaste(createPaste)
	println(createdPaste.RawUrl)
	
	paste, _ := client.GetPaste("ZA8U8CCQ")
	println(paste.Content)
	paste.Content += `\nconsole.log("There!")`
	client.SavePaste(paste)
	
	// Getting folder
	folder, _ := client.GetFolder("abcdefgh")
	for _, folderPaste := range folder.Pastes {
		println(folderPaste.Title)
	}


	// Getting current logged in user
	user, _ := client.GetUser()
	if user.LoggedIn {
		println("Hello: " + user.Name)
	}
	
	
	// Edit encrypted pastes
	password := "password"
	paste, _ = paste.Decrypt(password)
	paste.Content = "Hey"
	paste, _ = paste.Encrypt(password)
	_, err := client.SavePaste(paste) // (edited paste, error)
	if err == nil {
		println("Successful!")
    }
}
```