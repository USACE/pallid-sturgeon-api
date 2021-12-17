An Application Programming Interface (API) to manage Pallid Sturgeon data, built with Golang and Deployed on AWS Lambda.

# How to Develop

## Running a Database for Local Development

1. Install Docker and Docker Compose
2. Install postgres database with usace schema by using docker

## Running the GO API for Local Development

Either of these options starts the API at `localhost:8800`. 

**With Visual Studio Code Debugger**

You can use the launch.json file in this repository in lieu of `go run root/main.go` to run the API in the VSCode debugger.  This takes care of the required environment variables to connect to the database.

**Without Visual Studio Code Debugger**

Set the following environment variables and type `go run root/main.go` from the top level of this repository.

    "DB_USER": "",
    "DB_PASS": "",
    "DB_NAME": "",
    "LIB_DIR": "",
    "DB_HOST": "",
    "DB_PORT": "1521",
    "IPPK": "",


Note: When running the API locally, make sure environment variable `LAMBDA` is either **not set** or is set to `LAMBDA=FALSE`.
IPPK is NOT in use currently in AWS Dev.
Use the SIG key when setting IPPK value.