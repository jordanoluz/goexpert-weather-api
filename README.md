# Go Expert Weather API

This is a simple weather API project developed in Go that provides temperature information based on the postal code (zipcode), using Open Telemetry and Zipkin for tracing and spans.

## How to Run the Project Locally

### 1. Clone the repository

Clone this repository to your local machine:

```
git clone https://github.com/jordanoluz/goexpert-weather-api.git
```

### 2. Navigate to the project directory

Change into the project directory:

```
cd goexpert-weather-api
```

### 3. Run Docker Compose

To run the project using Docker, use Docker Compose to build and start the containers:

```
docker compose up -d --build
```

This will build the Docker images and start the application containers.

### 4. Automated Tests

The Docker Compose setup will automatically run tests using the go test command before starting the application.

### 5. Test the API

Once the container is running, you can test the API using the following `curl` commands:

- To check the POST request for a **valid zipcode**:

```
curl -X POST http://localhost:8181/weather \
     -H "Content-Type: application/json" \
     -d '{
         "cep": "29902555"
     }'

```

Expected response: City name and temperature data for the given zipcode.

### 6. Check the Tracing

After making the requests mentioned in step 5, you can access the following URL to view the tracing and span data:

```
http://localhost:9411
```

This will display the tracing results, allowing you to monitor the flow of requests through OpenTelemetry and Zipkin.
