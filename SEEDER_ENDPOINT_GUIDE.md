# Database Seeder Endpoint Guide

## Overview

The database seeder has been moved from automatic execution on startup to an on-demand endpoint. This provides better control over when initial data is populated and is more suitable for production environments.

## Endpoints

### 1. Run Database Seeder
- **Method**: `POST`
- **Endpoint**: `/admin/seed`
- **Authentication**: Required (Bearer token)
- **Description**: Executes the database seeder to populate initial data

#### Request Example:
```bash
curl -X POST http://localhost:8080/admin/seed \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

#### Response Example:
```json
{
  "success": true,
  "data": {
    "message": "Database seeder executed successfully",
    "status": "completed"
  }
}
```

### 2. Health Check
- **Method**: `GET`
- **Endpoint**: `/admin/health`
- **Authentication**: Not required
- **Description**: Returns the health status of the application

#### Request Example:
```bash
curl -X GET http://localhost:8080/admin/health
```

#### Response Example:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "message": "Application is running properly"
  }
}
```

## Usage Instructions

### 1. **First Time Setup**
After deploying the application, run the seeder to populate initial data:

```bash
# First, login to get a JWT token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-password"
  }'

# Use the token to run the seeder
curl -X POST http://localhost:8080/admin/seed \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

### 2. **Development Environment**
For development, you can run the seeder whenever you need to reset or populate test data:

```bash
# Set environment variables
export ENABLE_MODULES=warehouse

# Start the application
go run apps/main.go

# In another terminal, run the seeder
curl -X POST http://localhost:8080/admin/seed \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. **Production Environment**
In production, run the seeder only once after initial deployment:

```bash
# Run seeder after deployment
curl -X POST https://your-domain.com/admin/seed \
  -H "Authorization: Bearer YOUR_PRODUCTION_TOKEN"
```

## Security Considerations

### **Authentication Required**
- The seeder endpoint requires authentication to prevent unauthorized access
- Only authenticated users can execute the seeder
- Use strong JWT tokens in production

### **Rate Limiting** (Recommended)
Consider adding rate limiting to the admin endpoints:

```go
// Example rate limiting middleware (not implemented yet)
adminGroup.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5))) // 5 requests per minute
```

### **IP Whitelisting** (Optional)
For production environments, consider restricting admin endpoints to specific IP addresses.

## Error Handling

### **Common Error Responses**

#### 401 Unauthorized
```json
{
  "success": false,
  "error": "unauthorized access"
}
```

#### 500 Internal Server Error
```json
{
  "success": false,
  "error": "database connection failed"
}
```

## Implementation Details

### **File Structure**
```
shared/
└── handler/
    └── admin_handler.go    # Admin endpoints handler
apps/
└── main.go                # Updated with admin routes registration
```

### **Code Changes Made**

1. **Removed Automatic Seeder Execution**
   ```go
   // Before: Automatic execution on startup
   seeder.Seed(dbInstance)
   
   // After: Available via endpoint
   // Database seeder is now available via endpoint: POST /admin/seed
   ```

2. **Added Admin Handler**
   ```go
   type AdminHandler struct {
       db *gorm.DB
   }
   
   func (h *AdminHandler) RunSeeder(c echo.Context) error {
       seeder.Seed(h.db)
       return c.JSON(http.StatusOK, response)
   }
   ```

3. **Registered Admin Routes**
   ```go
   func registerAdminRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, authService *service.AuthService) {
       adminHandler := shandler.NewAdminHandler(db)
       adminGroup := e.Group("/admin")
       adminGroup.Use(authMiddleware.Middleware)
       adminGroup.POST("/seed", adminHandler.RunSeeder)
       e.GET("/admin/health", adminHandler.HealthCheck)
   }
   ```

## Benefits of This Approach

### **1. Better Control**
- Seeder runs only when explicitly requested
- No automatic data population on every startup
- Suitable for production environments

### **2. Security**
- Authentication required for seeder execution
- Prevents unauthorized data manipulation
- Audit trail through request logs

### **3. Flexibility**
- Can run seeder multiple times if needed
- Easy to integrate with deployment scripts
- Health check endpoint for monitoring

### **4. Production Ready**
- No risk of accidental data overwriting on restart
- Controlled data initialization process
- Better separation of concerns

## Swagger Documentation

The admin endpoints are documented in Swagger and available at:
- **URL**: `http://localhost:8080/swagger/`
- **Tags**: `admin`

## Monitoring and Logging

The seeder execution is logged with the following messages:
```
Running database seeder...
Database seeding completed
Admin routes registered successfully
```

Monitor these logs to track seeder execution and ensure proper functionality.

## Future Enhancements

Consider adding these features in future versions:
1. **Seeder Status Endpoint** - Check if seeder has been run
2. **Selective Seeding** - Run specific parts of the seeder
3. **Rollback Functionality** - Undo seeder changes
4. **Seeder History** - Track when and by whom seeder was executed
5. **Backup Before Seeding** - Automatic backup before running seeder