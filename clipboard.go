package main

import (
	"log"
	"os/exec"
	"time"
)

func monitor() {
	lastClientURL := ""
	ticker := time.NewTicker(1 * time.Second)
	for ; true; <-ticker.C {
		log.Println("Debug: Getting clipboard contents")
		clipboard, err := getClipboard()
		if err != nil {
			log.Printf("Error getting clipboard contents: %v", err)
			continue
		}
		meeting := extractZoomURL(clipboard)
		if meeting != nil && meeting.clientURL != lastClientURL {
			log.Println("Detected Zoom meeting, launching!")
			exec.Command("open", meeting.clientURL).Run()
			lastClientURL = meeting.clientURL
		}
	}

}

func getClipboard() (string, error) {
	pasteCmd := exec.Command("pbpaste")
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
