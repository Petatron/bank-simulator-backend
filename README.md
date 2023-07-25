# bank-simulator-backend
[![codecov](https://codecov.io/github/Petatron/bank-simulator-backend/branch/main/graph/badge.svg?token=SiNvQCTzQo)](https://codecov.io/github/Petatron/bank-simulator-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/Petatron/bank-simulator-backend)](https://goreportcard.com/report/github.com/Petatron/bank-simulator-backend)
[![Go Reference](https://pkg.go.dev/badge/github.com/Petatron/bank-simulator-backend.svg)](https://pkg.go.dev/github.com/Petatron/bank-simulator-backend)
[![Build Status](https://dev.azure.com/Petatron/bank_simulator_backend/_apis/build/status%2FPetatron.bank-simulator-backend?branchName=main)](https://dev.azure.com/Petatron/bank_simulator_backend/_build/latest?definitionId=6&branchName=main)

This project aimed to build a bank simulator system. The system features include:

- User account CUDR.
- Retrive the bank account info owned by users.
- Making transation between two accounts.
- Retrive transation records from user.

For backend system, typically using `Golang` as programming language.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)

## Installation

Get latest production code.

```bash
git clone https://github.com/Petatron/bank-simulator-backend.git
```

Install basic required environment.

- [Go installation](https://go.dev/dl/)
- [Docker installation](https://www.docker.com/)
- [PostgresSQL installation](https://www.postgresql.org/download/)

It is recommand to use `Homebrew` to manage and install if you are using Linux or try to use them on terminal. (Please make sure you have installed `Homebrew` before run below commands.)

```bash
brew install go
brew install docker
brew install postgresql
```

## Usage

### MakeFile

The project used MakeFile to set up docker image and migrate database to docker container.

```bash
# Pull docker image and start docker container.
make postgres
# Create database and set up username and password.
make createdb
# Remove database
make dropdb
# Create tables
make migrateup
# Remove database
make migratedown
# Gnerate Go code
make sqlc
```

### SQLc

The project used [SQLc](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html) to generate type-safe database connection Go code from SQL.

How to use SQLc in project?

```bash
make sqlc
```
