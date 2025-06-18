# Authentication Setup Guide

The application now uses HTTP-only cookies for authentication with real JWT tokens.

## Backend Changes

### 1. User Management Service
- **HTTP-only cookies**: Login now sets an `auth_token` cookie instead of returning the token in the response
- **Logout endpoint**: Added `/logout` endpoint that clears the authentication cookie
- **User info endpoint**: Added `/me` endpoint to get current user information
- **JWT middleware**: Added middleware to protect the `/me` endpoint

### 2. API Gateway
- **Cookie support**: Updated middleware to read JWT tokens from cookies
- **CORS configuration**: Updated to support credentials (cookies) with specific origins
- **Route protection**: Added `/user/*` routes to protected routes group

## Frontend Changes

### 1. Authentication Context
- **AuthContext**: Created a centralized authentication state management
- **HTTP-only cookies**: All API calls use `credentials: 'include'` to send cookies
- **Auto-authentication**: Checks authentication status on app load
- **Real API integration**: Uses actual API endpoints through the gateway

### 2. Components Updated
- **LoginForm**: Real authentication with error handling
- **RegisterForm**: Real registration with auto-login after success
- **Sidebar**: Shows user initials from actual user data
- **UserMenu**: Displays real user information
- **AuthScreen**: Simplified to remove unnecessary props

## Running the Application

### Prerequisites
1. Ensure PostgreSQL is running and configured
2. Ensure all backend services are running (user-management, api-gateway, etc.)

### Environment Variables
For the frontend, create a `.env` file with:
```
REACT_APP_API_URL=http://localhost:3000
```

### Development Servers
1. **API Gateway**: `http://localhost:3000`
2. **Frontend**: `http://localhost:3001` (React dev server)
3. **User Management**: `http://localhost:8080`

### Testing Authentication
1. Navigate to `http://localhost:3001`
2. Register a new user with valid information
3. Login with the registered credentials
4. Check that user initials appear in the sidebar
5. Verify logout functionality works properly

## Security Features
- **HTTP-only cookies**: Tokens are not accessible via JavaScript, preventing XSS attacks
- **CORS with credentials**: Properly configured for secure cross-origin requests
- **JWT validation**: All protected routes validate JWT tokens
- **Automatic logout**: Invalid or expired tokens automatically log users out

## API Endpoints
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration  
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/user/me` - Get current user info (protected) 