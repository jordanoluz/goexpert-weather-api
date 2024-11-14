# Go Expert Weather API

This is a simple weather API project developed in Go, which provides temperature information based on the postal code (zipcode).

## API URL

The API is published and available at the following URL:

**Base URL**: `https://goexpert-weather-api-1036920645078.southamerica-east1.run.app/weather`

You can use this base URL to test the API endpoints without needing to run the project locally.

## How to Run the Project Locally

### 1. Clone the repository

Clone this repository to your local machine:

```git clone https://github.com/jordanoluz/goexpert-weather-api.git```

### 2. Navigate to the project directory

```cd goexpert-weather-api```

### 3. Run Docker Compose

To run the project, use Docker Compose to build and start the containers:

```docker-compose up -d --build```

This will build the Docker image and start the container with the application.

### 4. Automated Tests

The Docker Compose setup will automatically run tests using the go test command before starting the application.

### 5. Test the API

Once the container is running, you can test the API using the following curl commands:

- To check the temperature for a **valid zipcode**:

```curl "http://localhost:8080/weather?zipcode=93010001"```

Expected response: Temperature data for the given zipcode.

- To test an **invalid zipcode**:

```curl "http://localhost:8080/weather?zipcode=1234"```

Expected response: `invalid zipcode`. 

- To test a zipcode **zipcode not found**:

```curl "http://localhost:8080/weather?zipcode=78654321"```

Expected response: `can not find zipcode`.