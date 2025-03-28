#!/bin/sh

# Start the Go API in the background
./api &

# Start the Node.js web app in the foreground
cd web
npm run preview -- --host 0.0.0.0 --port 3000