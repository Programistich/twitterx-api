# docs/AGENTS.md

**After all changes, please update any relevant files with the new information.**

## OpenAPI Specification Maintenance

This directory contains the OpenAPI specification for the Twitter API. The specification is manually maintained and served via Swagger UI.

### Files in this Directory

- **openapi.yaml** - OpenAPI 3.0.3 specification (source of truth for API documentation)
- **../fx/API_FxTwitter.md** - External FxTwitter API documentation reference

## How to Update OpenAPI Specification

When adding new features or modifying existing endpoints, follow these steps:

### 1. Adding a New Endpoint

Edit `openapi.yaml` and add the new path under the `paths` section:

```yaml
paths:
  /your/new/{endpoint}:
    get:
      summary: Brief description
      description: Detailed description of what this endpoint does
      operationId: uniqueOperationId
      tags:
        - Tweets
      parameters:
        - name: endpoint
          in: path
          required: true
          description: Parameter description
          schema:
            type: string
            example: "example_value"
      responses:
        '200':
          description: Success response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/YourResponseType'
        '500':
          description: Internal server error
          content:
            text/plain:
              schema:
                type: string
```

### 2. Adding a New Schema/Model

If your endpoint returns a new data structure, add it to `components/schemas`:

```yaml
components:
  schemas:
    YourNewModel:
      type: object
      required:
        - required_field1
        - required_field2
      properties:
        required_field1:
          type: string
          description: Field description
          example: "example_value"
        optional_field:
          type: string
          nullable: true
          description: Optional field description
```

### 3. Modifying Existing Endpoints

When changing endpoint behavior:
1. Update the endpoint description in `openapi.yaml`
2. Update parameter schemas if parameters changed
3. Update response schemas if response structure changed
4. Update examples to reflect new behavior

### 4. Verification Checklist

After updating `openapi.yaml`:

- [ ] YAML syntax is valid (use a YAML validator or run the app to check)
- [ ] All `$ref` references point to existing schemas
- [ ] Examples are provided for all request/response types
- [ ] Required vs optional fields are correctly marked
- [ ] HTTP status codes are appropriate
- [ ] Descriptions are clear and helpful
- [ ] Run the application and verify Swagger UI at http://localhost:8080/
- [ ] Test the "Try it out" functionality in Swagger UI
- [ ] Verify that `/openapi.yaml` endpoint returns the updated spec

### 5. Update Related Documentation

After updating `openapi.yaml`, also update:
- **README.md** - If adding major new endpoints or features
- **AGENTS.md** (root) - If changing architecture or request flow
- **main.go** - Implement the actual endpoint code

## OpenAPI Structure Overview

The `openapi.yaml` file follows this structure:

```
openapi: 3.0.3
├── info: API metadata (title, description, version)
├── servers: API server URLs
├── paths: All API endpoints
│   ├── /users/{username}/tweets
│   └── /users/{username}/tweets/{id}
├── components:
│   └── schemas: Reusable data models
│       ├── TweetsResponse
│       ├── FxTwitterResponse
│       ├── Tweet
│       ├── Author
│       ├── Media
│       └── ... (other models)
└── tags: Endpoint grouping
```

## Common Patterns

### Path Parameters

```yaml
parameters:
  - name: username
    in: path
    required: true
    description: Twitter username (without @)
    schema:
      type: string
      example: elonmusk
```

### Query Parameters

```yaml
parameters:
  - name: limit
    in: query
    required: false
    description: Maximum number of results
    schema:
      type: integer
      default: 20
      example: 50
```

### Request Body

```yaml
requestBody:
  required: true
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/YourRequestType'
      example:
        field1: "value1"
        field2: "value2"
```

### Nested Objects

Use `$ref` to reference existing schemas for cleaner code:

```yaml
properties:
  author:
    $ref: '#/components/schemas/Author'
  media:
    $ref: '#/components/schemas/Media'
```

### Arrays

```yaml
properties:
  tweet_ids:
    type: array
    items:
      type: string
    example:
      - "1234567890123456789"
      - "1234567890123456790"
```

### Nullable Fields

```yaml
properties:
  optional_field:
    type: string
    nullable: true
    description: This field may be null
```

## Swagger UI Configuration

The Swagger UI is configured in `main.go` with these settings:

```go
httpSwagger.Handler(
    httpSwagger.URL("/openapi.yaml"),        // Points to spec file
    httpSwagger.DeepLinking(true),           // Enable deep linking
    httpSwagger.DocExpansion("list"),        // Auto-expand endpoint list
    httpSwagger.DomID("swagger-ui"),         // DOM element ID
)
```

To modify Swagger UI behavior, update these options in `main.go:95-99`.

## Troubleshooting

### Swagger UI shows "Failed to load API definition"

- Check that `docs/openapi/openapi.yaml` exists and is readable
- Verify YAML syntax is valid
- Check application logs for file read errors
- Ensure Docker container has access to `docs/` directory (see `scripts/Dockerfile`)

### Schema validation errors

- Ensure all `$ref` references are correct
- Check that required fields are present in all schemas
- Verify enum values match exactly (case-sensitive)

### Changes not reflecting in Swagger UI

- Hard refresh your browser (Ctrl+Shift+R or Cmd+Shift+R)
- Restart the application
- Check that you're editing the correct `docs/openapi/openapi.yaml` file
- Verify the file is being copied to Docker container if using Docker

## Resources

- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Schema Validator](https://apitools.dev/swagger-parser/online/)
- [YAML Validator](https://www.yamllint.com/)
