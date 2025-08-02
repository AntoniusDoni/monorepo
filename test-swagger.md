# Testing Swagger UI with Authentication

## Steps to Test

1. **Start the server**:
   ```bash
   ENABLE_MODULES=warehouse go run apps/main.go
   ```

2. **Open Swagger UI**:
   Navigate to: http://localhost:8080/swagger/index.html

3. **Test Authentication Flow**:

   ### Step 1: Register with Office
   - Find the `POST /register-with-office` endpoint
   - Click "Try it out"
   - Use this sample data:
   ```json
   {
     "username": "testuser",
     "password": "password123",
     "email": "test@example.com",
     "office_code": "TEST01",
     "office_name": "Test Office",
     "office_address": "123 Test St",
     "office_city": "Test City"
   }
   ```
   - Execute the request

   ### Step 2: Login
   - Find the `POST /login` endpoint
   - Use these credentials:
   ```json
   {
     "username": "testuser",
     "password": "password123"
   }
   ```
   - Copy the `token` from the response

   ### Step 3: Authorize in Swagger
   - Click the **"Authorize"** button (🔒) at the top
   - In the "BearerAuth" field, enter: `Bearer YOUR_TOKEN_HERE`
   - Click "Authorize"

   ### Step 4: Test Protected Endpoints
   - Try `GET /v1/api/offices` - should work with authentication
   - Try `GET /v1/api/offices/active` - should return your test office
   - Try creating another office with `POST /v1/api/offices`

## Expected Results

✅ **Authentication endpoints** (register, login) should work without authorization
✅ **Protected endpoints** should return 401 without token
✅ **Protected endpoints** should work correctly with valid token
✅ **Swagger UI** should show the lock icon (🔒) for protected endpoints
✅ **Authorization button** should be visible at the top of Swagger UI

## Troubleshooting

- If you get 401 errors, make sure the token is properly formatted: `Bearer <token>`
- If endpoints don't appear, ensure `ENABLE_MODULES=warehouse` is set
- If Swagger UI doesn't load, check that the server is running on port 8080