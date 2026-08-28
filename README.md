Markdown
# Agnos Backend Assignment - Hospital Middleware API

A robust RESTful API middleware built with Go, designed to serve as a secure gateway between hospital staff and external Hospital Information Systems (HIS). 

This project strictly adheres to **Clean Architecture** principles and includes full infrastructure setup via **Docker Compose**, featuring an **Nginx** reverse proxy and a **PostgreSQL** database.

## 🚀 Tech Stack & Libraries
* **Language:** Go (1.25)
* **Framework:** Gin Web Framework
* **Database:** PostgreSQL (with `pgcrypto` for UUID generation)
* **ORM:** GORM
* **Security:** JWT (JSON Web Tokens) & bcrypt for password hashing
* **Infrastructure:** Docker, Docker Compose, Nginx

## 📁 Project Structure (Clean Architecture)
The source code is organized to decouple the business logic from external frameworks:

```text
.
├── cmd/api/            # Application entry point (main.go)
├── domain/             # Core business entities and interfaces
├── delivery/http/      # HTTP Handlers, Routing, and JWT Middleware
├── usecase/            # Business logic and rules (e.g., fallback to HIS API)
├── repository/         # Data persistence (Postgres) and External HTTP clients
├── nginx/              # Nginx reverse proxy configuration
├── init.sql            # Database schema and mock data seeding
└── docker-compose.yml  # Infrastructure as Code

```

🛠 Prerequisites
Docker and Docker Compose installed on your machine.
No need to install Go or PostgreSQL locally; Docker handles everything.

🏃‍♂️ How to Run
Clone the repository and navigate to the project root.
Start the services using Docker Compose:
Bash
docker-compose up -d
(Note: On the very first run, PostgreSQL will automatically execute init.sql to create tables and seed mock data).
The API is now accessible via Nginx at http://localhost (Port 80).

To stop the services and remove the database volume (reset data):
Bash
docker-compose down -v

📖 API Documentation & Usage
1. Create Staff
Creates a new staff member account.
Endpoint: POST /staff/create
cURL Example:
Bash
curl -X POST http://localhost/staff/create \
-H "Content-Type: application/json" \
-d '{"username": "admin_win", "password": "password123", "hospital": "HOSPITAL_A"}'

2. Login
Authenticates a staff member and returns a JWT token.
Endpoint: POST /staff/login
cURL Example:
Bash
curl -X POST http://localhost/staff/login \
-H "Content-Type: application/json" \
-d '{"username": "admin_win", "password": "password123", "hospital": "HOSPITAL_A"}'
(Copy the token from the response to use in the Search API).

3. Search Patient
Searches for a patient. If the patient is not found in the local database and a national_id or passport_id is provided, the middleware will attempt to fetch the data from the external HIS API and cache it locally.
Endpoint: GET /patient/search
Headers: Authorization: Bearer <your_jwt_token>
cURL Example (Search by Name):
Bash
curl -X GET "http://localhost/patient/search?first_name=สมชาย" \
-H "Authorization: Bearer <your_jwt_token>"
cURL Example (Search by National ID - triggers HIS API fallback):
Bash
curl -X GET "http://localhost/patient/search?national_id=1100112233445" \
-H "Authorization: Bearer <your_jwt_token>"

🧪 Running Unit Tests
Unit tests are implemented for the Usecase layer to verify business logic and middleware behavior using Dependency Injection and Mock Repositories.
To run the tests locally (requires Go installed):
Bash
go test ./usecase -v