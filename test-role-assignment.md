# Testing Role Assignment

## Overview
This test verifies that users are properly assigned the "admin" role during registration.

## Test Steps

### 1. Start the Server
```bash
ENABLE_MODULES=warehouse go run apps/main.go
```

### 2. Test RegisterWithOffice
```bash
curl -X POST http://localhost:8080/register-with-office \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com",
    "office_code": "TEST01",
    "office_name": "Test Office",
    "office_address": "123 Test St",
    "office_city": "Test City"
  }'
```

**Expected Response:**
```json
{
  "message": "Office and user registered successfully",
  "office_id": "some-uuid",
  "user_id": 1
}
```

### 3. Test Login (Verify Role Assignment)
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "user_identifier": 1,
  "role": "admin",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. Test Regular Register
```bash
# First, get an office ID from the previous test or create one
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "regularuser",
    "password": "password123",
    "email": "regular@example.com",
    "office_id": "OFFICE_ID_FROM_STEP_2"
  }'
```

### 5. Test Regular User Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "regularuser",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "user_identifier": 2,
  "role": "admin",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Success Criteria

✅ **RegisterWithOffice**: Creates user and office successfully
✅ **Role Assignment**: User can login and receives "admin" role
✅ **Regular Register**: Also assigns admin role to new users
✅ **No Role Errors**: No "user has no roles assigned" errors during login

## Troubleshooting

### If you get "user has no roles assigned" error:
1. Check that the seeder ran successfully (look for "Seeding completed." in logs)
2. Verify the admin role exists in the database
3. Check that the `assignAdminRole` function is being called

### If role assignment fails:
1. Check database connectivity
2. Verify the admin role was seeded properly
3. Look for error messages in the server logs

### Database Verification (Optional)
If you have database access, you can verify the role assignment:

```sql
-- Check if admin role exists
SELECT * FROM roles WHERE name = 'admin';

-- Check user roles
SELECT u.username, r.name as role_name 
FROM users u 
JOIN user_roles ur ON u.id = ur.user_id 
JOIN roles r ON ur.role_id = r.id;
```