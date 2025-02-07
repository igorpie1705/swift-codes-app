# Swift Codes API

## Introduction
This project is an API for managing SWIFT codes. It supports CRUD operations for SWIFT codes and stores them in a PostgreSQL database.

## Requirements
- **Docker** and **Docker Compose** (for running the application)
- **Postman** (optional, for testing the endpoints)
- **Git** (to clone the repository)
- **cURL** (optional, for testing the API in the terminal)

## Installation and Launch
### 1. Clone the Repository
First, open a terminal (Linux/Mac) or PowerShell (Windows) and type:
```sh
git clone https://github.com/igorpie1705/swift-codes-app.git
cd swift-codes-app
```

### 2. Start the Application in Docker
If you already have Docker installed, run the following command:
```sh
docker-compose up --build
```
This will start:
- PostgreSQL as the database
- The API in a Go container

The application will be accessible at **http://localhost:8080**.

## Project Structure
```
.
├── database/        # Database connection handling
├── handlers/        # API business logic
├── models/          # Model definitions
├── tests/           # Integration and unit tests
├── main.go          # Application entry point
└── docker-compose.yml  # Docker configuration
```

## API Endpoints

### Get a SWIFT Code
**GET /v1/swift-codes/:swift-code**
```sh
curl -X GET http://localhost:8080/v1/swift-codes/ABCDEFGXXX
```

### Get SWIFT Codes for a Country
**GET /v1/swift-codes/country/:countryISO2code**
```sh
curl -X GET http://localhost:8080/v1/swift-codes/country/PL
```

### Add a New SWIFT Code
**POST /v1/swift-codes**
```sh
curl -X POST http://localhost:8080/v1/swift-codes -H "Content-Type: application/json" -d '{
    "swiftCode": "NEWBANKGBXXX",
    "bankName": "New Bank",
    "address": "789 New St",
    "countryISO2": "GB",
    "countryName": "United Kingdom",
    "isHeadquarter": true
}'
```

### Delete a SWIFT Code
**DELETE /v1/swift-codes/:swift-code**
```sh
curl -X DELETE http://localhost:8080/v1/swift-codes/NEWBANKGBXXX
```

## Testing the API in Postman
1. Open **Postman**  
2. Add a new **GET/POST/DELETE** request  
3. Use the URL **http://localhost:8080/v1/swift-codes**  
4. Click "Send" and check the response  

## Running Tests
You can run tests manually in the terminal:
```sh
go test ./tests -v
```

### Running Tests in Docker
To run tests inside the Docker container, type:
```sh
docker-compose run app go test ./tests -v
```
This will verify the correct operation of the API and the database connection.

## Conclusion
If you have any questions or encounter issues, feel free to reach out! :)
