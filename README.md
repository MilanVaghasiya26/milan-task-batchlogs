# Go Language Base Project

### Things Covered

- Microservice based architecture
- Used zero log(https://github.com/rs/zerolog) for logging an info
- Used .env for reading config info

### Migration Readme

- Run database migration from the root of this project
- Created `migration.sh` file to execute migration cases defined in `migration.go` file
- See the documentation of [Migration](data_model/README.md)

### Seeder Readme

- Run seeder migration from the root of this project
- See the documentation of [Seeder](data_model/README.md)

### RUNNING GO LIKE IN NODEMON

1. fire "go install github.com/mitranim/gow@latest" in your project directory i.e. golang-starter-kit/
2. then go to the individual project using the terminal
   ```bash
     cd project
   ```
3. then fire "gow run . "
4. now change the file it will automatically run

#### Install using source code

#### Prerequisites

- OS: Linux or macOS or Windows
- Go: (Golang)(https://golang.org/dl/) >= v1.21

#### Project Setup
