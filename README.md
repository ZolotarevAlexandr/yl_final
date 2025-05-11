# Distributed Calculator

This project implements a distributed system for evaluating arithmetic expressions asynchronously.
It splits an input expression into independent tasks that can be processed concurrently.
The system is divided into two main components:

- Orchestrator: Accepts expressions via HTTP, parses them, decomposes them into tasks, and manages the execution order.
- Agent: Polls the orchestrator for tasks, performs the computations (with simulated delays), and reports the results back.


## How to run

Prerequisites:
- Go 1.24 or later
- Docker (optional)

Set environmental variables (or leave default ones defined in settings.go files):
- `ORCHESTRATOR_PORT` - Port of the orchestrator (default: `"8080"`)
- `ORCHESTRATOR_GRPC_PORT` - Port of the orchestrator gRPC (default: `50051`)
- `TIME_ADDITION_MS` – Delay (in milliseconds) for addition (default: `1000`)
- `TIME_SUBTRACTION_MS` – Delay for subtraction (default: `1000`)
- `TIME_MULTIPLICATIONS_MS` – Delay for multiplication (default: `1000`)
- `TIME_DIVISIONS_MS` – Delay for division (default: `1000`)
- `COMPUTING_POWER` – Number of concurrent agent goroutines to run (default: `2`)

### Run as separate modules:
- Run orchestrator:
    - ```bash
      cd .\orchestrator\cmd
      ```
    - ```bash
      go run main.go
      ```
- Run agent:
    - ```bash
      cd .\agent\cmd
      ```
    - ```bash
      go run main.go
      ```

### Run both modules simultaneously:
**Run locally with:**
```bash
go run main.go
```

> [!TIP]
> You might need to create following go.work in root to run it locally
> ```
> go 1.24
> toolchain go1.24.0
> use (
>     .
>     agent
>     calculator
>     db
>     orchestrator
>     grpc
> )
>```

**Or run in Docker container:**
```bash
docker build --tag "calculator" .
```
```bash
docker run -p 8080:8080 "calculator"
```

## API endpoints:

1. #### POST /api/v1/register
   Description:
   Registers a new user.
   
   **Successful Request (200 OK):**
   - Request:
     ```bash
     curl -X POST http://localhost:8080/api/v1/register \
          -H "Content-Type: application/json" \
          -d '{"login": "user", "password": "password"}'
     ```
   **Conflict (409 Conflict):**
   - When occurs:  
     Login already exists

2. #### /api/v1/login
   Description:
   Authenticates a user and returns a JWT token.
   
   **Successful Request (200 OK):**
   - Request:
     ```bash
     curl -X POST http://localhost:8080/api/v1/login
          -H "Content-Type: application/json"
          -d '{"login": "user", "password": "password"}'
     ```
   - Response:
     ```json
     {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDcwNjY0NDV9.WZa6TJksGVREZcJhSyUIBDMEMSkxloV7bij3ACEPxsc"
     }
     ```
   **Failed login (401 Unauthorized):**
   - When occurs:
     Credentials are not valid

> [!NOTE]
> All following endpoints now require authorization header  
> `Authorization: Bearer YOUR_TOKEN_HERE`

3. #### POST /api/v1/calculate
    Description:
    Submits an arithmetic expression for evaluation. The expression is parsed, converted into tasks, and stored for asynchronous processing.

    **Successful Request (201 Created):**
    - Request:
      ```bash
      curl -X POST http://localhost:8080/api/v1/calculate \
           -H "Content-Type: application/json" \
           -H "Authorization: Bearer YOUR_TOKEN_HERE" \
           -d '{"expression": "2+2*2"}'
      ```
    - Response:
      ```json
      {
          "id": "unique-expression-id"
      }       
      ```
    - When occurs:  
      The expression is valid and tasks are generated successfully.

    **Invalid Data (422 Unprocessable Entity):**
    - Request:
      ```bash
      curl -X POST http://localhost:8080/api/v1/calculate \
           -H "Content-Type: application/json" \
           -H "Authorization: Bearer YOUR_TOKEN_HERE" \
           -d '{"expression": "2+"}'
      ```
    - Response:  
      Status 422 (Unprocessable Entity) with message "error processing expression"
    - When occurs:  
      When given expression is not valid (e.g. missing operand)

    **Internal Error (500 Internal Server Error):**
    - When occurs:  
      An unexpected error occurs during tokenization, parsing, or task generation

4. #### GET /api/v1/expressions
   Description:  
   Retrieves a list of all expressions with their statuses and results.

   **Successful Request (200 OK):**
    - Request:
      ```bash
      curl http://localhost:8080/api/v1/expressions \
           -H "Authorization: Bearer YOUR_TOKEN_HERE"
      ```
    - Response:
      ```json
      {
        "expressions": [
          {
            "id": "a1d39298-d20d-4fa6-8d73-fe3cde5738e7",
            "expression": "2+2*2",
            "status": "done",
            "result": 6
          },
          {
            "id": "a1d39298-d20d-4fa6-8d73-fe3cde5738e7",
            "expression": "2+2*2",
            "status": "pending"
          }
        ]
      }
      ```
    This endpoint returns 200 unless a severe internal error occurs (500).

5. #### GET /api/v1/expressions/:id
   Description:  
   Retrieves a specific expression by its ID.

   **Successful Request (200 OK):**
    - Request:
      ```bash
      curl http://localhost:8080/api/v1/expressions/uuid \
           -H "Authorization: Bearer YOUR_TOKEN_HERE"
      ```
    - Response:
      ```json
      {
        "expression":
          {
            "id":"eec27e72-3654-4205-bc56-4bb7f97ac7cc",
            "expression":"(2+2)*(3+3)",
            "status":"done",
            "result":24
          }
      }
      ```
    - When occurs:  
      Existing id is given

   **Expression Not Found (404 Not Found):**
    - Request:
      ```bash
      curl http://localhost:8080/api/v1/expressions/non-existing-uuid
      ```
    - Response:
      Code 404 with message "not found"
    - When occurs:  
      Non existing id is given

Internal API for communicating between orchestrator and agents is made
using gRPC
