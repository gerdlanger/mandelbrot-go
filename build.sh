#!/bin/bash

echo MacOS
GOOS=darwin GOARCH=amd64 go build -o apfel

echo Windows
GOOS=windows GOARCH=amd64 go build -o apfel.exe

echo Linux ARM
GOOS=linux GOARCH=arm GOARM=7 go build -o apfel_arm

zip apfel_mac.zip apfel
zip apfel_win.zip apfel.exe
zip apfel_arm.zip apfel_arm
