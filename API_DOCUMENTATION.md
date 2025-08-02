# API Documentation

## Overview
This monorepo provides a comprehensive API for managing offices, users, and warehouses. The API includes authentication endpoints and office management capabilities.

## Base URL
```
http://localhost:8080
```

## Swagger Documentation
Access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Using Authentication in Swagger UI
1. First, use the `/login` endpoint to get a JWT token
2. Click the **"Authorize"** button (🔒) at the top of the Swagger UI
3. In the "BearerAuth" field, enter: `Bearer <your-jwt-token>`
   - Example: `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
4. Click **"Authorize"** to apply the token
5. Now you can test all protected endpoints directly in the UI

**Note:** The token will be automatically included in the Authorization header for all subsequent API calls.

## Authentication Endpoints

### 1. Register User (Existing Office)
**POST** `/register`

Register a new user with an existing office ID. The user is automatically assigned the "admin" role.

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "securepassword123",
  "email": "john@example.com",
  "office_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**
```json
{
  "message": "User registered successfully"
}
```

### 2. Register User with Office Creation
**POST** `/register-with-office`

Create a new office and register the first user for that office in a single operation. The user is automatically assigned the "admin" role.

**Request Body:**
```json
{
  "username": "admin_user",
  "password": "securepassword123",
  "email": "admin@newcompany.com",
  "office_code": "NYC01",
  "office_name": "New York Office",
  "office_address": "123 Main St",
  "office_city": "New York",
  "office_phone": "+1-555-0123"
}
```

**Response:**
```json
{
  "message": "Office and user registered successfully",
  "office_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 1
}
```

### 3. User Login
**POST** `/login`

Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "user_identifier": 1,
  "role": "admin",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Office Management Endpoints

All office endpoints require authentication and are prefixed with `/v1/api`.

### 1. Get All Offices
**GET** `/v1/api/offices`

Retrieve paginated list of offices with optional search.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `pageSize` (optional): Items per page (default: 10)
- `searchTerm` (optional): Search term for filtering

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "code": "NYC01",
        "name": "New York Office",
        "address": "123 Main St",
        "city": "New York",
        "phone": "+1-555-0123",
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total_count": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 2. Get Active Offices
**GET** `/v1/api/offices/active`

Retrieve all offices with active status.

### 3. Get Office by ID
**GET** `/v1/api/offices/{id}`

Retrieve a specific office by its UUID.

### 4. Create Office
**POST** `/v1/api/offices`

Create a new office.

**Request Body:**
```json
{
  "code": "LA01",
  "name": "Los Angeles Office",
  "address": "456 Sunset Blvd",
  "city": "Los Angeles",
  "phone": "+1-555-0456",
  "status": "active"
}
```

### 5. Update Office
**PUT** `/v1/api/offices/{id}`

Update an existing office.

### 6. Delete Office
**DELETE** `/v1/api/offices/{id}`

Delete an office by its ID.

## Warehouse Endpoints

### Get All Warehouses
**GET** `/v1/api/warehouses`

Retrieve paginated list of warehouses with optional search.

### Create Warehouse
**POST** `/v1/api/warehouses`

Create a new warehouse.

### Get Warehouse by ID
**GET** `/v1/api/warehouses/{id}`

Retrieve a specific warehouse by its ID.

### Update Warehouse
**PUT** `/v1/api/warehouses/{id}`

Update an existing warehouse.

### Delete Warehouse
**DELETE** `/v1/api/warehouses/{id}`

Delete a warehouse by its ID.

## Error Responses

All endpoints return consistent error responses:

```json
{
  "success": false,
  "error": "Error message description"
}
```

## Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `204 No Content`: Resource deleted successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid token
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Getting Started with Authentication
1. **Register a new office and user**: Use `POST /register-with-office` for initial setup
2. **Login**: Use `POST /login` with your credentials to get a JWT token
3. **Use the token**: Include it in the Authorization header for all subsequent requests

## User Roles

### Admin Role
- Automatically assigned to all newly registered users
- Has access to all API endpoints
- Can manage offices, warehouses, and other resources

**Note:** Currently, all users are assigned the "admin" role upon registration. Future versions may include role-based access control with different permission levels.

## Validation Rules

### User Registration
- Username: 3-50 characters
- Password: Minimum 6 characters
- Email: Valid email format
- Office ID: Valid UUID format

### Office Creation
- Code: 2-10 characters, unique
- Name: 3-100 characters
- Address, City, Phone: Optional fields

## Environment Setup

To enable the warehouse module (including offices), set:
```
ENABLE_MODULES=warehouse
```

## Development

To regenerate Swagger documentation:
```bash
swag init -g apps/main.go -o docs --parseDependency --parseInternal
```

## Testing the API

### Quick Start Example
1. **Create office and user**:
   ```bash
   curl -X POST http://localhost:8080/register-with-office \
     -H "Content-Type: application/json" \
     -d '{
       "username": "admin",
       "password": "password123",
       "email": "admin@company.com",
       "office_code": "HQ01",
       "office_name": "Headquarters",
       "office_address": "123 Business St",
       "office_city": "New York"
     }'
   ```

2. **Login to get token**:
   ```bash
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{
       "username": "admin",
       "password": "password123"
     }'
   ```

3. **Use token to access protected endpoints**:
   ```bash
   curl -X GET http://localhost:8080/v1/api/offices \
     -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
   ```

## Troubleshooting

### Common Issues

1. **"user has no roles assigned" error during login**:
   - This should no longer occur as users are automatically assigned the admin role
   - If it persists, check that the seeder has run and the admin role exists

2. **401 Unauthorized errors**:
   - Ensure the JWT token is properly formatted: `Bearer <token>`
   - Check that the token hasn't expired
   - Verify the user has the required role

3. **Office creation fails**:
   - Check that the office code is unique
   - Ensure all required fields are provided