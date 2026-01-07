# Distributed Task Management System

This is a microservices-based application designed for scalable, real-time collaboration. It uses Go (Golang) for the backend services, React.js for the frontend, RabbitMQ for asynchronous messaging, and MongoDB for data persistence.

## 1. Prerequisites
Ensure you have the following installed:
-   **Go (Golang)**: v1.19 or higher.
-   **Node.js**: v16 or higher (for React).
-   **MongoDB**: Running locally/Atlas.
-   **RabbitMQ**: Running locally on port `5672`.

## 2. Project Structure
The system consists of 5 main components:
-   `auth-service`: Authentication & JWT.
-   `Task-service`: Task management (Kanban).
-   `user-service`: User profiles.
-   `notification-service`: WebSocket alerts.
-   `taskboard-frontend`: React UI.

## 3. Setup & Run
Open **5 separate terminals** and run the following commands in the project root:

**Terminal 1 (Auth)**:
```bash
cd auth-service && go run main.go
```

**Terminal 2 (Task)**:
```bash
cd Task-service && go run main.go
```

**Terminal 3 (Notification)**:
```bash
cd notification-service && go run main.go
```

**Terminal 4 (User)**:
```bash
cd user-service && go run main.go
```

**Terminal 5 (Frontend)**:
```bash
cd taskboard-frontend && npm start
```

## 4. Usage
1.  Open `http://localhost:3000`.
2.  Register a new account.
3.  Login to see the Kanban board.
4.  Tasks created by Admin users will appear here. Updates will trigger real-time toasts.

## 5. Environment Variables
Ensure `.env` files exist in each service directory with valid `MONGO_URI` and `RABBITMQ_URL`.
## 6. Test Credentials (for Proof/Demo)
Role synchronization is implemented between Auth and User services.

Admin User
Email: Admin@test.com .
Password: @dminT3st .

Regular User
Email: ayat@gmail.com .
Password: Abcd123# .
