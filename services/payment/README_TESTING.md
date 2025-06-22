# Payment Service Testing

This document describes the testing setup and approach for the Payment Service.

## Test Structure

The payment service tests are organized into several files:

- `main_test.go` - Tests for the main HTTP endpoints
- `handlers/payment_handler_test.go` - Unit tests for payment handler functions
- `redis/client_test.go` - Tests for Redis client functionality
- `integration_test.go` - Integration tests for the entire service
- `test_helpers.go` - Common test utilities and mock helpers

## Test Types

### Unit Tests
Unit tests focus on testing individual functions in isolation:

- **Handler Tests**: Test individual handler functions with mocked dependencies
- **Redis Tests**: Test Redis client initialization and configuration
- **HTTP Tests**: Test HTTP request/response handling

### Integration Tests
Integration tests verify the complete workflow:

- **Health Endpoint**: Verify service health check
- **Deposit Endpoint**: Test deposit functionality with proper error handling
- **Error Cases**: Test invalid requests and missing parameters

## Dependencies

The test suite uses the following testing libraries:

```go
github.com/stretchr/testify v1.9.0           // Assertions and test suites
github.com/go-redis/redismock/v8 v8.11.5     // Redis mocking
```

## Running Tests

### Run All Tests
```bash
make test
# or
go test ./...
```

### Run Tests with Verbose Output
```bash
make test-verbose
# or
go test -v ./...
```

### Run Tests with Coverage
```bash
make test-coverage
# or
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Run Race Detection Tests
```bash
make test-race
# or
go test -race ./...
```

## Test Environment Setup

The tests use the following environment setup:

1. **Gin Test Mode**: All tests run with `gin.SetMode(gin.TestMode)`
2. **Mock Environment Variables**: Tests set `ALPACA_API_KEY` and `ALPACA_SECRET_KEY`
3. **Mock Redis**: Uses `redismock` to mock Redis operations
4. **Mock HTTP Server**: Uses `httptest.Server` to mock external API calls

## Test Coverage

Current test coverage includes:

- ✅ Health endpoint functionality
- ✅ Deposit endpoint parameter validation
- ✅ Redis client initialization
- ✅ HTTP request/response handling
- ✅ Error handling for missing headers
- ✅ JSON marshaling/unmarshaling
- ✅ Authentication header validation

## Known Limitations

1. **External API Testing**: Tests currently mock HTTP requests to Alpaca API
2. **Redis Integration**: Tests use mock Redis client instead of real Redis instance
3. **Environment Variables**: Some tests depend on environment variable setup

## Future Improvements

1. **Dependency Injection**: Refactor code to inject HTTP client for better testability
2. **Test Database**: Use test Redis instance for integration tests
3. **E2E Tests**: Add end-to-end tests with real external services
4. **Performance Tests**: Add benchmarks for critical paths
5. **Contract Tests**: Add tests to verify API contract compliance

## Mock Setup Examples

### Mock Redis Cache Hit
```go
helper := NewTestHelper()
defer helper.Cleanup()

achDetails := &handlers.ACHDetails{Id: "test-ach-id"}
helper.SetupMockRedisCache("account-123", true, achDetails)
```

### Mock Alpaca API
```go
helper := NewTestHelper()
defer helper.Cleanup()

responses := map[string]interface{}{
    "get_ach": []handlers.ACHDetails{{Id: "existing-ach"}},
    "transfer": map[string]string{"id": "transfer-123"},
}
helper.CreateMockAlpacaServer("account-123", responses)
```

## Troubleshooting

### Common Issues

1. **Redis Connection Errors**: Ensure Redis mock is properly set up
2. **Environment Variables**: Check that test environment variables are set
3. **HTTP Timeouts**: Mock servers should be properly configured
4. **Test Isolation**: Ensure tests clean up after themselves

### Debug Tips

1. Use `-v` flag for verbose test output
2. Check test logs for HTTP request details
3. Verify mock expectations are met
4. Use debugger or print statements for complex scenarios

## Contributing

When adding new tests:

1. Follow existing test patterns and structure
2. Use descriptive test names that explain what is being tested
3. Include both positive and negative test cases
4. Mock external dependencies appropriately
5. Ensure tests are deterministic and don't depend on external state
6. Add proper cleanup in test teardown 