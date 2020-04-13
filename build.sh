#!/bin/bash

echo MacOS
GOOS=darwin GOARCH=amd64 go build -o apfel

echo Windows
GOOS=windows GOARCH=amd64 go build -o apfel.exe

echo Linux ARM
GOOS=linux GOARCH=arm GOARM=7 go build -o apfel_arm

