# Distributed Task Board System

A distributed task management system built with Go (Microservices) and React.

## Services

- **Auth Service**: Handles user registration, login, and role management.
- **Task Service**: Manages task creation, assignment, and status updates.
- **User Service**: Manages user profiles and syncs role updates.
- **Notification Service**: Real-time notifications via WebSockets.
- **Frontend**: React-based dashboard for Admins and Users.

## Setup Instructions

1.  **Prerequisites**:
    - Go 1.21+
    - Node.js & npm
    - MongoDB
    - RabbitMQ

2.  **Running the Backend**:
    ```bash
    # Terminal 1
    cd auth-service && go run main.go

    # Terminal 2
    cd user-service && go run main.go

    # Terminal 3
    cd Task-service && go run main.go

    # Terminal 4
    cd notification-service && go run main.go
    ```

3.  **Running the Frontend**:
    ```bash
    cd taskboard-frontend
    npm start
    ```

## Test Credentials (for Proof/Demo)

Role synchronization is implemented between Auth and User services.

### Admin User
*   **Email**: `Admin@test.com`
*   **Password**: `@dminT3st`

*(Note: If this user does not exist, please register a new user with these credentials and promote them to admin using the database or API)*

### Regular User
*   **Email**: `ayat@gmail.com`
*   **Password**: `Abcd123#`
