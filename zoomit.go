package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/caseymrm/menuet"
)

func openZoom() {
	fmt.Println("Here we gooooooo....")
	exec.Command("open", "zoommtg://zoom.us/join?confno=123456789").Run()
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{
		{
			Text:    fmt.Sprintf("ZOOM IT! %d events", numEvents),
			Clicked: openZoom,
		},
	}
	return items
}

var numEvents = 0

func main() {
	fmt.Println("👋 Booting up...")
	// First authorize the user's gcal
	srv := authorizeCalendar()
	t := time.Now().Format(time.RFC3339)
	fmt.Println(srv, t)
	// _, err := srv.Events.List("primary").ShowDeleted(false).
	// 	SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	// fmt.Println(err)
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	// }
	// numEvents = len(events.Items)

	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "🗓",
	})
	menuet.App().Children = menuItems
	menuet.App().RunApplication()
}
