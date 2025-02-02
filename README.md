Task Management API

RESTful APIs built with Golang and the Gin framework for managing tasks. 
This API supports public and protected endpoints, JWT-based authentication, database integration with MySQL, and robust error handling.

Features
- Task CRUD Operations: Create, read, update, and delete tasks.
- Public & Protected Routes: Public endpoints for task retrieval and protected routes secured with JWT.
- JWT Authentication: Token-based security for sensitive operations.
- Input Validation: Ensures data integrity using Gin's built-in validation.
- Integration & Unit Tests: Tests for both API endpoints and core logic.
- Docker Support: Easily deployable with Docker and Docker Compose.

Prerequisites:
- Go 1.23
- Docker

Installation:
- git clone https://github.com/rajnishkmishra/task_service.git
- cd task_service
- go mod tidy

Running the Application:
- For running the application run below command from your terminal
<pre> docker-compose up --build </pre>
- This will spin up the API along with the MySQL database.
- And the output in terminal will look something like this in below screeshot.

<img width="700" alt="Screenshot 2025-02-03 at 12 48 16 AM" src="https://github.com/user-attachments/assets/f5256c0f-0144-4e24-851a-c6ee55611acc" />

API Endpoints:
- Public Endpoints:
    - GET /v1/login - Login
    - GET /v1/public/tasks - List all tasks
- Protected Endpoints: (Requires JWT)
    - GET /v1/tasks/{id} - Retrieve a task by ID
    - POST /v1/tasks - Create a new task
    - PUT /v1/tasks/{id} - Update an existing task
    - DELETE /v1/tasks/{id} - Delete a task

Authentication
- Send the JWT token in the token header
  <pre> token:Bearer {access_token} </pre>

Import the provided Postman collection: TaskManagement.postman_collection.json

Steps to hit APIs:
- First hit /v1/login APIs, This API required phone_number in the request.
- /v1/login API will give you an access_token in response.

<img width="700" alt="Screenshot 2025-02-03 at 12 06 48 AM" src="https://github.com/user-attachments/assets/efcb3f5b-13d1-4019-961c-2af7fc473833" />

- Copy this acces_token and use it to hit protected endpoints. In header you need to pass key as "token" and value should be "Bearer {access_token}"

<img width="700" alt="Screenshot 2025-02-03 at 12 09 21 AM" src="https://github.com/user-attachments/assets/2b54a036-586a-420f-b12d-f8f4888bb43e" />
        
-  Public endpoints will work without any authentication.

Run the tests:
    - To run unit test navigate to services/task_service/test/task_service_test and run below command
    <pre>go test</pre>
    - This will run all the test present in this go file.
    - Similarly to run integration test, navigate to test package and run below command.
    <pre>go test</pre>
