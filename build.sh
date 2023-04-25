#!/bin/bash

set -euo pipefail

GOOS=windows GOARCH=amd64 go build -o bin/sagemaker-pipelines-dashboard-amd64.exe
GOOS=windows GOARCH=386 go build -o bin/sagemaker-pipelines-dashboard-386.exe
GOOS=darwin GOARCH=amd64 go build -o bin/sagemaker-pipelines-dashboard-amd64-darwin
GOOS=linux GOARCH=amd6 go build -o bin/sagemaker-pipelines-dashboard-amd64-linux
GOOS=linux GOARCH=386 go build -o bin/sagemaker-pipelines-dashboard-386-linux
