#!/bin/bash

echo MacOS
GOOS=darwin GOARCH=amd64 go build -o apfel

echo Windows
GOOS=windows GOARCH=amd64 go build -o apfel.exe

echo Linux x64
GOOS=linux GOARCH=amd64 go build -o apfel_linux

echo Linux ARM
GOOS=linux GOARCH=arm GOARM=7 go build -o apfel_arm

zip apfel_mac.zip apfel
zip apfel_win.zip apfel.exe
zip apfel_linux.zip apfel_linux
zip apfel_arm.zip apfel_arm
