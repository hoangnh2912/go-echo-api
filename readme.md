# ![RealWorld Example App](logo.png)

> ### Golang/Echo codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged fullstack application built with **Golang/Echo** including CRUD operations, authentication, routing, pagination, and more.

## Getting started

### Install Golang (go1.11+)

Please check the official golang installation guide before you start. [Official Documentation](https://golang.org/doc/install)
Also make sure you have installed a go1.11+ version.

For more info and detailed instructions please check this guide: [Setting GOPATH](https://github.com/golang/go/wiki/SettingGOPATH)

### Clone the repository

Clone this repository:

```bash
➜ git clone https://github.com/hoangnh2912/go-echo-api.git
```

Switch to the repo folder

```bash
➜ cd go-echo-api
```

### Copy .env.example to .env

```bash
➜ cp .env.example .env
```

### Install dependencies

```bash
➜ go mod download
```

### Run development

```bash
➜ make dev
```

### Build

```bash
➜ go build
```

### Tests

```bash
➜ go test ./...
```

## Folder structure

```
📦go-echo-api
 ┣ 📂controllers
 ┃ ┣ 📜article.go
 ┃ ┗ 📜user.go
 ┣ 📂db
 ┃ ┗ 📜db.go
 ┣ 📂docs
 ┃ ┣ 📜docs.go
 ┃ ┣ 📜swagger.json
 ┃ ┗ 📜swagger.yaml
 ┣ 📂middleware
 ┃ ┗ 📜jwt.go
 ┣ 📂model
 ┃ ┣ 📜article.go
 ┃ ┗ 📜user.go
 ┣ 📂requests
 ┃ ┣ 📜article.go
 ┃ ┗ 📜user.go
 ┣ 📂responses
 ┃ ┣ 📜article.go
 ┃ ┗ 📜user.go
 ┣ 📂services
 ┃ ┣ 📜article.go
 ┃ ┗ 📜user.go
 ┣ 📂utils
 ┃ ┣ 📜errors.go
 ┃ ┣ 📜jwt.go
 ┃ ┗ 📜utils.go
 ┣ 📜.dockerignore
 ┣ 📜.example.env
 ┣ 📜.gitignore
 ┣ 📜Dockerfile
 ┣ 📜LICENSE
 ┣ 📜Makefile
 ┣ 📜docker-compose.dev.yml
 ┣ 📜docker-compose.prod.yml
 ┣ 📜go.mod
 ┣ 📜go.sum
 ┣ 📜logo.png
 ┣ 📜main.go
 ┗ 📜readme.md
```
