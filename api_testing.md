# API Testing Guide: Auth Module

This guide describes how to test the Login and Logout endpoints for the Task List & Timesheet Backend.

## Prerequisites

1.  **Database**: Ensure the PostgreSQL database is running on port 5433 (as configured in `docker-compose.yml`).
2.  **Migration & Seed**: The database must have the initial schema and at least one test user.
    ```bash
    # Apply initial migration
    cat migrations/0001_initial.up.sql | docker exec -i task-list-timesheet-be-db-1 psql -U postgres -d taskdb
    
    # Apply name column migration
    cat migrations/0002_add_name_to_users.up.sql | docker exec -i task-list-timesheet-be-db-1 psql -U postgres -d taskdb
    
    # Seed test user (password: password123)
    docker exec -i task-list-timesheet-be-db-1 psql -U postgres -d taskdb -c "INSERT INTO users (name, email, password) VALUES ('Test User', 'test@example.com', '\$2a\$10\$sCdsZptkklMc9kHRgQhNAuEsI2KOL8Nxdj7VYuk.BsA/oFgRXvjWO');"
    ```

## 1. Start the Application

```bash
make run
```
The server will start on `http://localhost:8080`.

## 2. Test Register

```bash
curl -i -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ades", "email":"ades@example.com", "password":"password123"}'
```

**Expected Response (201 Created):**
```json
{
  "data": "Ok"
}
```

## 3. Test Login

Send a `POST` request to `/api/users/login` with valid credentials.

```bash
curl -i -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com", "password":"password123"}'
```

**Expected Response (200 OK):**
```json
{
  "data": {
    "token": "your-received-token-here"
  }
}
```

## 3. Test Logout

Send a `POST` request to `/api/users/logout` with the token received from login in the `Authorization` header.

```bash
curl -i -X POST http://localhost:8080/api/users/logout \
  -H "Authorization: Bearer <your-received-token-here>"
```

**Expected Response (200 OK):**
```json
{
  "data": "success"
}
```

## 5. Test Get User Info

```bash
curl -i -X GET http://localhost:8080/api/users/info \
  -H "Authorization: Bearer <your-received-token-here>"
```

**Expected Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "Ades",
    "email": "ades@example.com",
    "expired_token": "2026-04-01T11:47:27Z"
  }
}
```

## 6. Verification

After logout, verify that the token has been removed from the database:
```bash
docker exec -it task-list-timesheet-be-db-1 psql -U postgres -d taskdb -c "SELECT * FROM user_tokens;"
```
The result should show `(0 rows)`.

## 5. Test Invalid Credentials

```bash
curl -i -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com", "password":"wrongpassword"}'
```

**Expected Response (401 Unauthorized):**
```json
{
  "error": "Unauthorized"
}
```
