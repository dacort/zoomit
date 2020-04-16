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

- [ ] Crashes the first time you authenticate
- [ ] Crashes if you click on the icon before the calendar syncs
- [ ] Formatting needs some prettification on the menu
- [ ] Needs a proper downloadable app and Google app

## Features to add

- [ ] Show dates/days in addition to times
- [ ] Option to just open Zoom for you...why wait! (ala [@ConnorPM](https://twitter.com/ConnorPM/status/1250473781707132928?s=20))