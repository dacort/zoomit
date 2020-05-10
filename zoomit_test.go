package main

import (
	"testing"

	"google.golang.org/api/calendar/v3"
)

func TestBasicURL(t *testing.T) {
	expectedURI := "zoommtg://zoom.us/join?confno=123584"
	rawLink := "https://zoom.us/j/123584"
	zoom := extractZoomURL(rawLink)

	if zoom.clientURL != expectedURI {
		t.Errorf("Did not extract basic URL: got `%s`, want `%s`", zoom.clientURL, expectedURI)
	}
}

func TestPasswordURL(t *testing.T) {
	expectedURI := "zoommtg://zoom.us/join?confno=8857663&pwd=Thisw0uldbeapassword"
	rawLink := "https://bigco.zoom.us/j/8857663?pwd=Thisw0uldbeapassword"
	zoom := extractZoomURL(rawLink)

	if zoom.clientURL != expectedURI {
		t.Errorf("Did not extract password URL: got `%s`, want `%s`", zoom.clientURL, expectedURI)
	}
}

func TestEventLocationExtract(t *testing.T) {
	rawLink := "https://zoom.us/j/123584"
	expectedURI := "zoommtg://zoom.us/join?confno=123584"
	event := &calendar.Event{
		Location: rawLink,
	}

	zm := findZoomURLInEvent(event)
	if (zm.originalURL != rawLink) {
		t.Errorf("Did not extract URL from event: got `%s`, want `%s`", zm.originalURL, rawLink)
	}
	if (zm.clientURL != expectedURI) {
		t.Errorf("Did not extract URL from event: got `%s`, want `%s`", zm.clientURL, expectedURI)
	}
}

func TestEventConfenceExtract(t *testing.T) {
	rawLink := "https://zoom.us/j/123584"
	expectedURI := "zoommtg://zoom.us/join?confno=123584"
	event := &calendar.Event{
		Location: "",
		ConferenceData: &calendar.ConferenceData{
			EntryPoints: []*calendar.EntryPoint{
				{
					EntryPointType: "video",
					Uri: rawLink,
				},
			},
		},
	}

	zm := findZoomURLInEvent(event)
	if (zm.originalURL != rawLink) {
		t.Errorf("Did not extract URL from event: got `%s`, want `%s`", zm.originalURL, rawLink)
	}
	if (zm.clientURL != expectedURI) {
		t.Errorf("Did not extract URL from event: got `%s`, want `%s`", zm.clientURL, expectedURI)
	}
}