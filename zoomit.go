package main

import (
	"fmt"
	"os/exec"

	"github.com/caseymrm/menuet"
)

func openZoom() {
	fmt.Println("Here we gooooooo....")
	exec.Command("open", "zoommtg://zoom.us/join?confno=123456789").Run()
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{
		{
			Text:    "ZOOM IT!",
			Clicked: openZoom,
		},
	}
	return items
}

func main() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "🗓",
	})
	menuet.App().Children = menuItems
	menuet.App().RunApplication()
}
