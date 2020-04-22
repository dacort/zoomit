# Zoom It!

A simple menubar app that lists upcoming meetings and allows you to open Zoom with ease. 🤝 📹

![Screenshot](images/screenshot.png)

## Overview

There's a few things necessary to get this working:

- Authentication with calendar service
- Menubar app
- Zoom Opener

## Status

🔥 Under active development 🔥

This works! But you need to create your own app as documented [here](https://developers.google.com/calendar/quickstart/go).

Once that's done, build and run the app: `go build && ./zoomit`

Open the URL printed to the console and follow the instructions. 

The app will unfortunately crash after you paste the token back in, but just run `./zoomit` again.

It will sync with Google Calendar every hour and if any of your meeting locations match "zoom.us/j/<ZOOM_ID>",
you'll be able to expand that entry and click on the Zoom link!

Only the next 10 meetings are shown.

## Known Issues

- [x] Crashes the first time you authenticate
- [x] Crashes if you click on the icon before the calendar syncs
- [x] Formatting needs some prettification on the menu
- [ ] Needs a proper downloadable app and Google app
- [ ] App must be validated: https://support.google.com/cloud/answer/7454865

## Features to add

- [ ] Show dates/days in addition to times
- [ ] Option to just open Zoom for you...why wait! (ala [@ConnorPM](https://twitter.com/ConnorPM/status/1250473781707132928?s=20))

## References

https://github.com/golang/oauth2/blob/master/google/example_test.go
https://github.com/golang/go/blob/1abf3aa55bb8b346bb1575ac8db5022f215df65a/src/net/http/server.go#L2783
https://developers.google.com/calendar/quickstart/go
https://martinfowler.com/articles/command-line-google.html
https://developers.google.com/identity/protocols/oauth2/native-app#custom-uri-scheme