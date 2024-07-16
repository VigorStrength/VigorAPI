# Vigor Backend

## Introduction

This backend, built with Go, Gin, and MongoDB, serves as the backbone for both the users' mobile app and the SaaS CRM administrators' platform. The project leverages technology to modernize fitness, allowing users to improve their fitness by monitoring workouts, nutrition, and activities. Administrators can seamlessly generate workout and nutritional plans, monitor sales, and manage other tasks. The project is developed with three environments: development, production, and staging.

## Technologies Used

- [Golang](https://golang.org/)
- [Gin Framework](https://gin-gonic.com/)
- [Air Live Reload](https://github.com/cosmtrek/air)
- [MongoDB](https://www.mongodb.com/)
- [Firebase](https://firebase.google.com/)
- [Docker](https://www.docker.com/)

## Installations

### Clone

Clone this project to your local machine:

```sh
    git clone https://github.com/yourusername/vigor-backend.git
```

### Setup

#### Without Docker

1. Setup your local environment variables:

   - `VIGOR_DB_URI`
   - `VIGOR_DB_NAME`
   - `JWT_SECRET_KEY`

   You can use the files **.env.development** and **.env.staging** as example on how to locally setup environment variables globally on your computer.

2. Configure the `.air.toml` file:
   Ensure that you correctly set up the environment under the `env` key:
   - For development: `VIGOR_ENV=dev`
   - For staging/testing: `VIGOR_ENV=test`

For production, the environment setup is intended to be handled within the CI/CD pipeline as contributors merge the code to the main branch.

#### Using Docker

To be determined, not yet entirely set up.

## How to Run

1. Install dependencies:

   ```sh
   go mod tidy
   ```

2. Start the server with live reloading

   ```sh
       air
   ```

3. Start the server without live reloading:

   ```sh
       go run main.go
   ```

## Tests

### On Terminal

Run unit tests:

```sh
    go test -v ./tests/unit_tests/"desired subdirectory"
```

### On Postman

To generate **postamn_collection.json** file, run:

```sh
    go run scripts/generate_postman_collection.go
```

## Not Completed

- Docker Configuration
- CI/CD Configuration
- Some tests need to be rewritten to fit the updated codebase
