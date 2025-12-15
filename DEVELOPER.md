# Development

This documentation will guide you through the process of setting up your dev environment for modifying, building and running the Cortex Cloud Go SDK locally 

## Dependencies

#### Required

  - Go 1.24.0+

#### Optional

  - Copywrite v0.22.0+
    - `brew tap hashicorp/tap`
    - `brew install hashicorp/tap/copywrite`

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

If you are running the entire acceptance test suite, your API key must be associated with a user that has Instance Administrator permissions.

You must also provide the API URL for the Cortex Cloud tenant.

Once you have obtained your API key and URL, set the following environment variables with the appropriate values:
  - `CORTEXCLOUD_API_URL`
  - `CORTEXCLOUD_API_KEY`
  - `CORTEXCLOUD_API_KEY_ID`

