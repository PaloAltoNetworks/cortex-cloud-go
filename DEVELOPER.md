# Development

This documentation will guide you through the process of setting up your dev environment for modifying, building and running the Cortex Cloud Go SDK locally 

## Dependencies

#### Required

  - Go 1.24.0+

#### Optional

  - Copywrite v0.22.0+
    - `brew tap hashicorp/tap`
    - `brew install hashicorp/tap/copywrite`

<!-- TODO: list additional optional dependencies (gosec) with installation instructions -->

## Workspace

Run the `work` recipe to initialize/update the Go workspace file:
```
make work
```

## Build

Use the `build` recipe to build all modules:
```
make build
```

## Testing

Use the `test` recipe to execute both the unit and acceptance test suites:
```
make test
```

Note that the acceptance test suite requires a Cortex Cloud API key and key ID with the appropriate permissions. See the [Acceptance Tests](#acceptance-tests) section below for more information.

### Unit Tests

Use the `test-unit` recipe to execute the unit test suite:
```
make test-unit
```
You may also override the `TEST_PACKAGE` variable to execute the unit tests for a specific package:
```
make test-unit TEST_PACKAGE=internal/app
```

### Acceptance Tests

Use the `test-acc` recipe to run the acceptance test suite:
```
make test-acc
```

To run the acceptance test suite, you must have a Cortex Cloud API key and key ID.
<!-- TODO: add link to API key creation doc page -->

If you are running the entire acceptance test suite, your API key must be associated with a user that has Instance Administrator permissions.

You must also provide the API URL for the Cortex Cloud tenant.
<!-- TODO: add guidance for finding API URL -->

Once you have obtained your API key and URL, set the following environment variables with the appropriate values:
  - `CORTEX_API_URL`
  - `CORTEX_API_KEY`
  - `CORTEX_API_KEY_ID`
  - `CORTEX_API_KEY_TYPE` (set to `"standard"` or `"advanced"`)

## Request Tracking

The SDK automatically adds request tracking headers to all API calls for improved observability and debugging.

### X-Request-ID Header

Every request includes a unique `X-Request-ID` header for tracing:

```
X-Request-ID: req_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
```

**Format:** `req_<32-hex-characters>` (128-bit cryptographically secure random)

You can provide your own request ID via context:

```go
import "github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"

ctx := context.Background()
ctx = client.WithRequestID(ctx, "my-custom-request-id")

// Use ctx in API calls
assetGroups, err := platformClient.ListAssetGroups(ctx, req)
```

**Use cases:**
- Correlating with external tracing systems (OpenTelemetry, Jaeger)
- Multi-step workflows with consistent tracking
- Debugging specific user sessions

### User-Agent Header

The SDK sends a versioned User-Agent identifying:
- SDK version
- Module and version
- Go version
- Operating system and architecture

**Format:**
```
cortex-cloud-go/<sdk-version> (<module>/<module-version>; go<go-version>; <os>/<arch>)
```

**Example:**
```
User-Agent: cortex-cloud-go/1.0.0 (platform/1.0.0; go1.25.1; darwin/arm64)
```

You can customize the User-Agent:

```go
import "github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"

client, err := platform.NewClient(
    config.WithAgent("my-app/1.0.0"),
    // ... other options
)
```

Or append to the default User-Agent:

```go
import (
    "github.com/PaloAltoNetworks/cortex-cloud-go/version"
    "github.com/PaloAltoNetworks/cortex-cloud-go/platform"
)

customUA := version.UserAgentWithCustom(
    platform.ModuleName,
    platform.Version,
    "my-app/2.0.0",
)

client, err := platform.NewClient(
    config.WithAgent(customUA),
    // ... other options
)
```

### Logging Request IDs

Request IDs are automatically logged at INFO level:

```
2025-11-21T14:30:00Z INFO API request started [map[endpoint:/asset_groups method:GET request_id:req_abc123]]
2025-11-21T14:30:01Z INFO API request completed [map[request_id:req_abc123 status_code:200]]
```

Enable DEBUG logging to see full request/response dumps with request IDs:

```bash
export CORTEX_LOG_LEVEL=debug
```

**Debug output:**
```
---[ REQUEST req_abc123 ]-----------------------------
GET /asset_groups HTTP/1.1
Host: api-tenant.xdr.us.paloaltonetworks.com
User-Agent: cortex-cloud-go/1.0.0 (platform/1.0.0; go1.25.1; darwin/arm64)
X-Request-ID: req_abc123
...
-----------------------------------------------------

---[ RESPONSE req_abc123 ]----------------------------
HTTP/1.1 200 OK
Content-Type: application/json
...
-----------------------------------------------------
```

### Request ID Propagation

Request IDs propagate through context across function calls:

```go
func processAssetGroups(ctx context.Context, client *platform.Client) error {
    // Request ID from parent context is automatically used
    groups, err := client.ListAssetGroups(ctx, nil)
    if err != nil {
        return err
    }

    for _, group := range groups {
        // Same request ID used for all calls in this context
        if err := processGroup(ctx, client, group); err != nil {
            return err
        }
    }
    return nil
}
```

### Performance Impact

Request tracking has minimal overhead:
- Request ID generation: ~1μs
- User-Agent generation: ~234ns
- Total overhead: < 0.001%

### Additional Documentation

For comprehensive information on request tracking, including examples, best practices, and troubleshooting, see [REQUEST_TRACKING.md](REQUEST_TRACKING.md).

