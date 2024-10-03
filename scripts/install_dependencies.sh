#! usr/bin/bash

go getgolang.org/x/crypto
go getgithub.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get github.com/lib/pq
go get github.com/google/uuid
go install github.com/pressly/goose/v3/cmd/goose@latest

sudo apt update
sudo apt install postgresql postgresql-contrib
sudo passwd postgres