package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"time"

	"github.com/caseymrm/menuet"
	"google.golang.org/api/calendar/v3"
)

var nextTenEvents *calendar.Events

// Example Zoom URL: https://zoom.us/j/930721398
var reZoomURL = regexp.MustCompile(`zoom.us/j/(\d+)$`)

func noMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(items, menuet.MenuItem{
		Text:     "Recent posts",
		FontSize: 12,
	})
	return items
}

func menuItems() []menuet.MenuItem {
	// The below was just an attempt to not make go panic...didn't work
	items := []menuet.MenuItem{}

	for _, event := range nextTenEvents.Items {
		// date := event.Start.DateTime
		// if date == "" {
		// 	date = event.Start.Date
		// }
		ts, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		items = append(items, menuet.MenuItem{
			Text:     fmt.Sprintf("%-15s %s", ts.Format("3:04 PM"), event.Summary),
			Children: zoomer(event.Location),
		})
	}

	return items
}

func zoomer(zoomURL string) func() []menuet.MenuItem {
	zoomMatch := reZoomURL.FindStringSubmatch(zoomURL)
	if zoomURL == "" || len(zoomMatch) == 0 {
		return nil
	}

	zoomScheme := fmt.Sprintf("zoommtg://zoom.us/join?confno=%s", zoomMatch[1])

	return func() []menuet.MenuItem {
		return []menuet.MenuItem{
			{
				Text: zoomURL,
				Clicked: func() {
					exec.Command("open", zoomScheme).Run()
				},
			},
		}
	}

}

func calendarSync(srv *calendar.Service) {
	ticker := time.NewTicker(1 * time.Minute)
	for ; true; <-ticker.C {
		fmt.Println("⌚️ Syncing events")
		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		// fmt.Println(err)
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		nextTenEvents = events
		if len(events.Items) == 0 {
			fmt.Println("No upcoming events found.")
		}

		menuet.App().MenuChanged()
	}

}

func fakeData() []menuet.MenuItem {
	return []menuet.MenuItem{
		{
			Text:     fmt.Sprintf("%-15s%s", "1:00 PM", "Do some stuff"),
			Children: zoomer(""),
		},
		{
			Text:     fmt.Sprintf("%-15s%s", "1:30 PM", "Important meeting"),
			Children: zoomer("https://zoom.us/j/12345678"),
		},
		{
			Text:     fmt.Sprintf("%-15s%s", "3:00 PM", "Keyboard Cat 🐈"),
			Children: zoomer("https://zoom.us/j/12345678"),
		},
	}
}

func main() {
	fmt.Println("👋 Booting up...")

	// Ensure that we've authenticated
	// If you just want to play around, uncomment these next few lines and replace
	// menuItems with fakeData below when setting menuet.App().Children below.
	srv := authorizeCalendar()
	go calendarSync(srv)

	app := menuet.App()
	app.SetMenuState(&menuet.MenuState{
		Title: "🗓",
	})
	app.Name = "ZoomIt!"
	app.Label = "com.github.dacort.zoomit"
	app.Children = menuItems
	app.RunApplication()
}

// Arrange that main.main runs on main thread.
// https://github.com/golang/go/wiki/LockOSThread
func init() {
	runtime.LockOSThread()
}
