#!/bin/bash

function execute_pipe() {
  local commands=("$@")

  for command in "${commands[@]}"; do
    echo -e "\033[1;33m[ $command ]\033[0m"
    if ! eval "$command" ; then
      return 1
    fi
  done
}

commands=(
  "git pull origin main"
  "cd ./frontend"
  "npm install"
  "npm audit fix"
  "npm run build"
  "cd ../"
  "go mod download"
  "go clean -cache"
  "go build -o ./migrations ./cmd/migrations/main.go"
  "./migrations apply"
  "go build -o ./main ./cmd/main/main.go"
  "pkill -f ./main"
  "./main"
)

execute_pipe "${commands[@]}"