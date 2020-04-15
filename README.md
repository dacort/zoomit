# Zoom It!

A simple menubar app that lists upcoming meetings and allows you to open Zoom with ease. 🤝 📹

## Overview

There's a few things necessary to get this working:

- Authentication with calendar service
- Menubar app
- Zoom Opener

## OK...

Let's try and do #2 and #3, that should be easy enough.

## Sweet!

That works, now for google calendar?

Follow along on https://developers.google.com/calendar/quickstart/go

Copy `credentials.json` to local dir and run the test code.

Almost works, but we get a core dump when we try to fetch calendar events. :(