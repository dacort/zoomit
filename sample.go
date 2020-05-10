package main

import (
	"fmt"

	"github.com/caseymrm/menuet"
)

func noMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(items, menuet.MenuItem{
		Text:     "Recent posts",
		FontSize: 12,
	})
	return items
}

func fakeMenuItems() []menuet.MenuItem {
	fakeZoom := &ZoomMeeting{"https://zoom.us/j/12345678", "zoommtg://zoom.us/join?confno=12345678"}
	return []menuet.MenuItem{
		{
			Text:     fmt.Sprintf("%-15s%s", "1:00 PM", "Do some stuff"),
			Children: zoomDetails(nil),
		},
		{
			Text:     fmt.Sprintf("%-15s%s", "1:30 PM", "Important meeting"),
			Children: zoomDetails(fakeZoom),
		},
		{
			Text:     fmt.Sprintf("%-15s%s", "3:00 PM", "Keyboard Cat 🐈"),
			Children: zoomDetails(fakeZoom),
		},
	}
}
