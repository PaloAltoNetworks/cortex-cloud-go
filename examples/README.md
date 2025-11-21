# cortex-cloud-go Examples

This directory contains practical examples demonstrating how to use the cortex-cloud-go SDK.

## Prerequisites

Before running these examples, ensure you have:

1. **Go 1.24.0+** installed
2. **Valid Cortex Cloud API credentials**
3. **Network access** to Cortex Cloud API endpoints

## Setup

### 1. Set Environment Variables

```bash
export CORTEX_API_URL="https://api-tenant.xdr.us.paloaltonetworks.com"
export CORTEX_API_KEY="your-api-key"
export CORTEX_API_KEY_ID="1"
export CORTEX_API_KEY_TYPE="standard"  # or "advanced"
```

### 2. Install Dependencies

```bash
cd cortex-cloud-go
go mod download
```

## Examples

### Basic Usage

**File:** [`basic_usage/main.go`](basic_usage/main.go)

Demonstrates:
- Creating a client with environment variables
- Creating a client with explicit configuration
- Listing asset groups
- Basic error handling

**Run:**
```bash
cd examples/basic_usage
go run main.go
```

**Expected Output:**
```
Found 5 asset groups
  - group-1: Production Servers
  - group-2: Development Servers
  - group-3: Database Servers
  - group-4: Web Servers
  - group-5: Test Servers
```

---

### Custom Request ID

**File:** [`custom_request_id/main.go`](custom_request_id/main.go)

Demonstrates:
- Setting custom request IDs
- Getting or generating request IDs
- Request ID propagation through functions
- Logging with request IDs

**Run:**
```bash
cd examples/custom_request_id
go run main.go
```

**Expected Output:**
```
Found 5 asset groups (request ID: my-workflow-12345)
Using request ID: req_c2aee4202a7069e04e52c59c072818b3
Found 5 asset groups
[req_abc123] Starting workflow
[req_abc123] Executing step1
[req_abc123] Step1 found 5 groups
[req_abc123] Executing step2
[req_abc123] Step2 completed
[req_abc123] Workflow completed
```

**Use Cases:**
- Multi-step workflows with consistent tracking
- Correlating with external tracing systems
- Debugging specific user sessions

---

### Custom User-Agent

**File:** [`custom_user_agent/main.go`](custom_user_agent/main.go)

Demonstrates:
- Setting custom User-Agent
- Appending to default User-Agent
- Getting SDK version information

**Run:**
```bash
cd examples/custom_user_agent
go run main.go
```

**Expected Output:**
```
Example 1: Found 5 asset groups with custom User-Agent
Example 2: Found 5 asset groups with appended User-Agent
User-Agent: cortex-cloud-go/1.0.0 (platform/1.0.0; go1.25.1; darwin/arm64) my-app/2.0.0

SDK Version Information:
  sdk_version: 1.0.0
  sdk_name: cortex-cloud-go
  git_commit: abc123def456
  build_date: 2025-11-21T14:30:00Z
  go_version: go1.25.1
  os: darwin
  arch: arm64
```

**Use Cases:**
- Identifying your application in server logs
- Tracking SDK version adoption
- Custom analytics and monitoring

---

### Debug Logging

**File:** [`debug_logging/main.go`](debug_logging/main.go)

Demonstrates:
- Enabling debug logging
- Viewing full request/response dumps
- Request ID in debug logs

**Run:**
```bash
cd examples/debug_logging
go run main.go
```

**Expected Output:**
```
---[ REQUEST req_abc123 ]-----------------------------
GET /asset_groups HTTP/1.1
Host: api-tenant.xdr.us.paloaltonetworks.com
User-Agent: cortex-cloud-go/1.0.0 (platform/1.0.0; go1.25.1; darwin/arm64)
X-Request-ID: req_abc123
Content-Type: application/json
x-xdr-auth-id: 1
Authorization: ***
-----------------------------------------------------

---[ RESPONSE req_abc123 ]----------------------------
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 1234
...
-----------------------------------------------------

Found 5 asset groups
```

**Use Cases:**
- Troubleshooting API issues
- Understanding request/response flow
- Debugging authentication problems

---

## Common Patterns

### Error Handling

```go
import "github.com/PaloAltoNetworks/cortex-cloud-go/errors"

groups, err := client.ListAssetGroups(ctx, nil)
if err != nil {
    // Check for API errors
    if apiErr, ok := err.(*errors.APIError); ok {
        log.Printf("API Error: %s (code: %s, status: %d)", 
            apiErr.Message, apiErr.ErrorCode, apiErr.StatusCode)
    }
    
    // Check for SDK errors
    if sdkErr, ok := err.(*errors.SDKError); ok {
        log.Printf("SDK Error: %s (code: %s)", 
            sdkErr.Message, sdkErr.Code)
    }
    
    return err
}
```

### Context with Timeout

```go
import "time"

// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// API call will timeout after 30 seconds
groups, err := client.ListAssetGroups(ctx, nil)
```

### Request ID Propagation

```go
func parentFunction(ctx context.Context, client *platform.Client) error {
    // Request ID automatically propagates to child functions
    return childFunction(ctx, client)
}

func childFunction(ctx context.Context, client *platform.Client) error {
    // Uses same request ID as parent
    requestID := client.GetRequestID(ctx)
    log.Printf("[%s] Processing...", requestID)
    
    groups, err := client.ListAssetGroups(ctx, nil)
    return err
}
```

## Troubleshooting

### Import Errors

If you see import errors when running examples:

```bash
# From the cortex-cloud-go root directory
make work
make build
```

### Authentication Errors

If you get 401 Unauthorized errors:

1. Verify your API credentials are correct
2. Check that `CORTEX_API_KEY_TYPE` is set correctly (`"standard"` or `"advanced"`)
3. Ensure your API key has the required permissions

### Network Errors

If you get connection errors:

1. Verify the API URL is correct
2. Check network connectivity to Cortex Cloud
3. Verify firewall rules allow outbound HTTPS

## Additional Resources

- [REQUEST_TRACKING.md](../REQUEST_TRACKING.md) - Complete request tracking guide
- [DEVELOPER.md](../DEVELOPER.md) - Development guide
- [README.md](../README.md) - SDK overview

## Contributing

To add a new example:

1. Create a new directory under `examples/`
2. Add a `main.go` file with your example
3. Update this README with a description
4. Ensure the example follows Go best practices
5. Test the example with real API credentials

---

**Last Updated:** 2025-11-21  
**SDK Version:** 1.0.0