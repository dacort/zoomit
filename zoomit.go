package main

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/caseymrm/menuet"
	"google.golang.org/api/calendar/v3"
)

var nextTenEvents *calendar.Events

// Example Zoom URL: https://zoom.us/j/930721398
var reZoomURL = regexp.MustCompile(`zoom.us/j/(\d+)$`)

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}

	// We may not have synced the calendar items yet
	if nextTenEvents == nil {
		return append(items, menuet.MenuItem{
			Text: "Syncing...",
		})
	}

	var date string
	for _, event := range nextTenEvents.Items {
		ts, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		// Handle all-day events
		if fmt.Sprintf("%s", event.Start.DateTime) == "" {
			ts, _ = time.Parse("2006-01-02", event.Start.Date)
		}
		td := ts.Format("Monday")
		if td != date {
			items = append(items, menuet.MenuItem{
				Text:     strings.ToUpper(td),
				FontSize: 10,
			})
			date = td
		}
		items = append(items, menuet.MenuItem{
			Text:     fmt.Sprintf("  %-15s %s", ts.Format("03:04 PM"), event.Summary),
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
		log.Println("⌚️ Syncing events")
		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		nextTenEvents = events
		if len(events.Items) == 0 {
			log.Println("No upcoming events found.")
		}

		menuet.App().MenuChanged()
	}

}

func main() {
	log.Println("👋 Booting up...")

	// Ensure that we've authenticated
	// If you just want to play around, uncomment these next few lines and replace
	// menuItems with fakemenuItems below when setting menuet.App().Children below.
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

	// While we're here, set up logging!
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "zoomit")
	if e == nil {
		log.SetOutput(io.MultiWriter(logwriter, os.Stdout))
	}
}
