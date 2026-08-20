#!/bin/zsh

file=".env"
if [ -e "$file" ]; then
    export $(cat .env)
else
    file=".env.example"
    if [ -e "$file"]; then
        export $(cat .env.example)
    fi
fi

go build -o ./app main.go

./app
