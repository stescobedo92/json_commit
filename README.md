# json_commit

This project is an implementation of an in-memory, concurrent commit log service with a RESTful HTTP API. It's written in Go and follows modern software design practices like dependency injection and structured error handling.

## ✨ Features
- In-Memory Log: Fast and efficient record storage.

- Concurrency Safe: Uses sync.RWMutex to allow multiple simultaneous reads, optimizing performance.

- RESTful API: Clear and predictable endpoints for producing (writing) and consuming (reading) records.

- Structured JSON Responses: All responses, both success and error, follow a consistent JSON format.

- Modular Design: The log's logic is decoupled from the HTTP server, making unit testing and extensibility easier.

🚀 Getting Started
To run the service, simply clone the repository and execute the main entry point.

## 📖 API Endpoints
The API provides two main endpoints to interact with the log.

1. Produce a Record
   Adds a new record to the end of the log.

Endpoint: POST /records

Description: Creates a new record. The data value must be a Base64 encoded string, as the underlying type is []byte.

Request Body:

```json
{
    "record": {
        "value": "SGVsbG8gV29ybGQh"
    }
}
```

Success Response (201 Created):

```json
{
    "offset": 0
}
```

Error Response (400 Bad Request):

```json
{
    "error": "Invalid request payload"
}
```