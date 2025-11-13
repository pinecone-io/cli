# E2E Integration Tests

Opt-in integration tests that exercise the built `pc` binary against real Pinecone endpoints.

Run locally:

```bash
# build to a temp binary automatically
PC_E2E=1 go test ./test/e2e -tags=e2e -v

# or use a prebuilt binary (e.g., goreleaser artifact)
PC_E2E=1 PC_BIN=./dist/pc go test ./test/e2e -tags=e2e -v

# use pc installed on your local $PATH
PC_E2E=1 PC_E2E_USE_PATH=1 go test ./test/e2e -tags=e2e -v
```

Environment variables:

- PC_E2E=1 to enable tests
- PC_E2E_DEBUG=1 to log commands and outputs
- PINECONE_ENVIRONMENT=production|staging (default: production)
- Service account flow:
  - PINECONE_CLIENT_ID
  - PINECONE_CLIENT_SECRET
- API key flow:
  - PINECONE_API_KEY (project-scoped)
- Target context (optional for target/api-key tests):
  - PC_E2E_ORG_ID
  - PC_E2E_PROJECT_ID
- Serverless params (optional):
  - PC_E2E_CLOUD (default: aws)
  - PC_E2E_REGION (default: us-east-1)
  - PC_E2E_DIMENSION (default: 8)
