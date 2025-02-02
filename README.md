Task Management API

RESTful APIs built with Golang and the Gin framework for managing tasks. This API supports public and protected endpoints, JWT-based authentication, database integration with MySQL, and robust error handling.

Features
Task CRUD Operations: Create, read, update, and delete tasks.
Public & Protected Routes: Public endpoints for task retrieval and protected routes secured with JWT.
JWT Authentication: Token-based security for sensitive operations.
Input Validation: Ensures data integrity using Gin's built-in validation.
Integration & Unit Tests: Tests for both API endpoints and core logic.
Docker Support: Easily deployable with Docker and Docker Compose.

Prerequisites:
Go 1.23
Docker

Installation:
git clone https://github.com/rajnishkmishra/task_service.git
cd task_service
go mod tidy

Running the Application:
For running the application run below command from your terminal
docker-compose up --build
This will spin up the API along with the MySQL database.
And the output in terminal will look something like in below screenshot.



API Endpoints:
    Public Endpoints:
        GET /v1/login - Login
        GET /v1/public/tasks - List all tasks
    Protected Endpoints: (Requires JWT)
        GET /v1/tasks/{id} - Retrieve a task by ID
        POST /v1/tasks - Create a new task
        PUT /v1/tasks/{id} - Update an existing task
        DELETE /v1/tasks/{id} - Delete a task

Authentication:
    Send the JWT token in the token header:
        token:Bearer <your_token>

Import the provided Postman collection: TaskManagement.postman_collection.json

Steps to hit APIs:
    1. First hit /v1/login APIs, This API required phone_number in request.
    2. /v1/login API will give you an access_token in response.
    3. Copy this acces_token and use it to hit protected endpoints listed above.
    4. Public endpoints work without any authentication.
