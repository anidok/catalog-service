# Catalog Service

A Go-based REST API for managing catalog services, supporting search, create, update, delete, and detail retrieval, with OpenSearch as the backend.

---

## Prerequisites

- **Go 1.24** or higher
- **Docker** (for running OpenSearch via Docker Compose).

---

## Setup
Follow this section to run app on your local.
1. **Run OpenSearch**
   - Run docker compose:
     ```sh
     make compose-up
     ```
     Wait for a few seconds for opensearch to setup. Check if opensearch dashboard is accessible at  http://localhost:5601

2. **Run Application**
     ```sh
     make run
     ```
     The API server will be accessible at port 4000.
---

## API Endpoints & Sample cURL

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
### Get Service by ID

```sh
curl -X GET "http://localhost:4000/api/services/<id>" \
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

## Authentication/Authorization Using Kong API Gateway
Kong is used for authentication and authorization (JWT + ACL).  
Kong runs on port **8000** (proxy) and **8001** (admin).  

### Start all services (including Kong):

1. **Scale down previous docker compose**
   ```sh
   make compose-down
   ```

2. **Scale up new compose file containing images for application and kong**
   ```sh
   make compose-up-kong
   ```
   
3. **Generate a JWT token:**
   ```sh
   make jwt-generate
   ```
   
2. **Verify a token (Optional):**
   ```sh
   make jwt-verify token=<your-token>
   ```
   
#### Example cURL with JWT

```sh
curl -X GET "http://localhost:8000/api/services" \
  -H "X-Correlation-ID: test-corr-id" \
  -H "Authorization: Bearer <your-jwt-token>"
```

#### Kong JWT Auth

- All `/api/*` endpoints are protected by JWT and ACL.
- Use the following JWT secret for testing:
  - **key:** `kong-jwt-auth`
  - **secret:** `some-key`
  - **consumer:** `jwt-user`
  - **group:** `catalog-group`

You can generate a JWT token for testing using the secret above.

#### Scale down
```sh
make compose-down-kong
```
---


## Running OpenSearch (Docker Compose)

For integration tests and local development, use Docker Compose:

```sh
docker-compose up -d
```

- OpenSearch will be available on port **9200**.
- The OpenSearch Dashboard will be accessible at [http://localhost:5601](http://localhost:5601).Wait for a few seconds for the dashboard to load.

---

## Running Tests

- **Setup**
    - Install mockery (Optional, unless you're changing interfaced and want to regenerate mocks)
        ```sh
        make install-mockery
        ```

- **Unit Tests**
    - Generate Mocks
        ```sh
        make generate-mocks
        ```
    - Run unit tests
        ```sh
        make unit-test
        ```

- **Integration Tests** (requires OpenSearch running)
  ```sh
  make integration-test
  ```

---

## Running Linter

- **Install golangci-lint**
  ```sh
  make install-golangci-lint
  ```
  
- **To check code quality and style**
    ```sh
    make lint
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
- Kong will block requests without a valid JWT and proper ACL group.
- All other API usage remains the same, just add the JWT header.

---
