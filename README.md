# Catalog Service

A Go-based REST API for managing catalog services, supporting search, create, update, delete, and detail retrieval, with OpenSearch as the backend.

---

## Prerequisites

- **Go 1.24**
- Docker (for running OpenSearch via Docker Compose).

---

## Setup

1. **Clone the repository**
   ```sh
   git clone https://github.com/anidok/catalog-service.git
   cd catalog-service
   ```

2. **Configuration**
   - Edit `application.yaml` for environment-specific settings (OpenSearch host, ports, etc).

3. **Dependencies**
   - Install Go modules:
     ```sh
     make deps
     ```

4. **Generate Mocks (for unit tests)**
   - If you change interfaces or want to regenerate mocks:
     ```sh
     make generate-mocks
     ```

5. **Linting**
   - To check code quality and style:
     ```sh
     make lint
     ```

---

## Running OpenSearch (Docker Compose)

For integration tests and local development, use Docker Compose:

```sh
docker-compose up -d
```

- OpenSearch will be available on port **9200**.
- The OpenSearch Dashboard will be accessible at [http://localhost:5601](http://localhost:5601).

---

## Database Migration

Apply DB/index migrations:

```sh
make migrate
```

---

## Initial Data Ingestion

Load initial test data into OpenSearch:

```sh
make ingest
```

---

## Running Application

  ```sh
  make run-api
  ```

---

## Running Tests

- **Unit Tests**
  ```sh
  make unit-test
  ```

- **Integration Tests** (requires OpenSearch running)
  ```sh
  make integration-test
  ```

---

## API Endpoints & Sample cURL

### Get Service by ID

```sh
curl -X GET "http://localhost:4000/api/services/<id>" \
  -H "X-Correlation-ID: test-corr-id"
```

### Search Services

#### 1. Without any search query and pagination parameters (defaults to page=1, limit=10)
```sh
curl -X GET "http://localhost:4000/api/services" \
  -H "X-Correlation-ID: test-corr-id"
```

#### 2. Without search query but with page and limit
```sh
curl -X GET "http://localhost:4000/api/services?page=2&limit=5" \
  -H "X-Correlation-ID: test-corr-id"
```

#### 3. With search query matching service name
```sh
curl -X GET "http://localhost:4000/api/services?q=forex" \
  -H "X-Correlation-ID: test-corr-id"
```

#### 4. With search query matching some description part
```sh
curl -X GET "http://localhost:4000/api/services?q=special%20rates" \
  -H "X-Correlation-ID: test-corr-id"
```

#### 5. With search query matching exact one document
```sh
curl -X GET "http://localhost:4000/api/services?q=Forex%20Card" \
  -H "X-Correlation-ID: test-corr-id"
```

#### 6. With search query not matching anything
```sh
curl -X GET "http://localhost:4000/api/services?q=nonexistentquery" \
  -H "X-Correlation-ID: test-corr-id"
```

### Create Service

```sh
curl -X POST "http://localhost:4000/api/services" \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: test-corr-id" \
  -d '{
    "name": "Forex Card",
    "description": "Forex card for students with special rates",
    "versions": [
      { "version_number": "1.0", "Details": "Initial" }
    ]
  }'
```

### Update Service (append version, update description)

```sh
curl -X PUT "http://localhost:4000/api/services/<id>" \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: test-corr-id" \
  -d '{
    "description": "Updated description",
    "versions": [
      { "version_number": "2.0", "Details": "Second release" }
    ]
  }'
```

### Delete Service

```sh
curl -X DELETE "http://localhost:4000/api/services/<id>" \
  -H "X-Correlation-ID: test-corr-id"
```

---

## Design Considerations & Trade-offs

- **OpenSearch as Backend:**  
  Chosen for full-text search and scalability. All service data is stored and queried in OpenSearch.

- **API Design:**  
  RESTful endpoints with clear separation for search, detail, create, update, and delete.  
  Only `description` and `versions` can be updated; `name` is immutable after creation.

- **Validation:**  
  Strict validation for required fields, versioning, and update constraints.

- **Structured Logging:**  
  All API and repository operations use structured logging for better traceability and debugging. Correlation IDs are supported for end-to-end request tracing.

- **Configuration:**  
  All environment-specific and sensitive settings (like OpenSearch hosts, timeouts, etc.) are managed via a central `application.yaml` config file, making the service easy to configure for different environments.

- **Testing:**  
  - **Unit tests** mock OpenSearch and focus on logic.
  - **Integration tests** require a running OpenSearch instance (via Docker Compose).
  - Test data is loaded/cleaned up automatically.

- **Error Handling:**  
  - Consistent error response format with codes and causes.  
  - The service includes a panic recovery middleware that catches panics, logs the error, and returns a standard error response instead of crashing the application.  
  - Graceful shutdown is supported to ensure all in-flight requests are completed before the server exits.

- **Assumptions:**  
  - Service `name` is unique and immutable.
  - Version numbers are strings and must be provided for each version.
  - OpenSearch index is managed externally or via migration scripts.

- **Trade-offs:**  
  - No partial updates (PATCH); only full update for allowed fields.
  - No authentication/authorization included.
  - No pagination links except for "next".
  - Logging and config are kept simple and extensible, but do not include advanced features like dynamic reload or log aggregation out of the box.

---

## Notes

- Make sure OpenSearch is running before running integration tests or ingesting data.
- Use the provided Makefile targets for common tasks.
- See `testdata/services.json` for example data used in tests.
- See `data.jsonl` for example data used for initial ingestion.

---
