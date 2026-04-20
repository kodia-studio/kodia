# 📋 OpenAPI/Swagger Documentation Guide

Complete guide to auto-generating REST API documentation with Swagger/OpenAPI in Kodia Framework.

**Table of Contents:**
- [What is OpenAPI/Swagger](#what-is-openapiseragger)
- [Setup](#setup)
- [Documenting Endpoints](#documenting-endpoints)
- [Documenting Data Models](#documenting-data-models)
- [Swagger UI](#swagger-ui)
- [Best Practices](#best-practices)
- [Client SDK Generation](#client-sdk-generation)
- [Examples](#examples)

---

## What is OpenAPI/Swagger?

**OpenAPI (Swagger)** is a standard format for describing REST APIs. It provides:

- ✅ **Machine-readable API specification** - Computers can parse and understand your API
- ✅ **Interactive API documentation** - Developers can test endpoints in browser
- ✅ **Client SDK generation** - Auto-generate TypeScript, Python, Java clients
- ✅ **API contract** - Server and client stay in sync
- ✅ **Discoverability** - APIs show up in API directories

---

## Setup

### Installation

```bash
# Install swag command-line tool
go install github.com/swaggo/swag/cmd/swag@latest

# Verify installation
swag version
```

### Dependencies

Add to `backend/go.mod`:

```bash
cd backend
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

---

## Documenting Endpoints

### Basic Endpoint Documentation

Every handler function should have Swagger comments above it:

```go
// @Summary Brief description of what endpoint does
// @Description Longer explanation of the endpoint
// @Tags CategoryName
// @Accept json
// @Produce json
// @Param paramName query string false "Description of parameter"
// @Success 200 {object} ResponseModel "Description of success response"
// @Failure 400 {object} ErrorResponse "Description of error"
// @Router /path/to/endpoint [http-method]
// @Security Bearer
func (h *Handler) MyEndpoint(c *gin.Context) {
    // Implementation
}
```

### Complete Example: Create User

```go
// @Summary Create a new user
// @Description Register a new user account with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /users [post]
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Example (request) {"name": "John Doe", "email": "john@example.com"}
// @Example (response) {"id": "123", "name": "John Doe", "email": "john@example.com"}
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: "Invalid request"})
        return
    }
    
    user, err := h.service.CreateUser(c.Request.Context(), &req)
    if err != nil {
        c.JSON(500, ErrorResponse{Error: "Server error"})
        return
    }
    
    c.JSON(201, user)
}
```

### Annotation Tags

| Tag | Purpose | Example |
|-----|---------|---------|
| `@Summary` | Short description (1 line) | "Create a new user" |
| `@Description` | Long description | "Register a new user account..." |
| `@Tags` | Endpoint category | "Users" |
| `@Accept` | Content types accepted | "json, xml" |
| `@Produce` | Content types returned | "json" |
| `@Param` | Request parameters | "name query string true "User name"" |
| `@Success` | Success response | "200 {object} UserResponse" |
| `@Failure` | Error response | "400 {object} ErrorResponse" |
| `@Router` | Endpoint path & method | "/users [post]" |
| `@Security` | Required security scheme | "Bearer" |

---

## Documenting Data Models

### Using Struct Tags

```go
// User represents a user account
type User struct {
    // User unique identifier (auto-generated)
    // @required
    // @example 550e8400-e29b-41d4-a716-446655440000
    ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

    // User email address
    // @required
    // @example user@example.com
    // @minLength 5
    // @maxLength 255
    Email string `json:"email" validate:"required,email" example:"user@example.com"`

    // User full name
    // @required
    // @example John Doe
    Name string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`

    // User role (user, admin, moderator)
    // @required
    // @enum user,admin,moderator
    // @example user
    Role string `json:"role" enum:"user,admin,moderator" example:"user"`

    // Account creation timestamp (RFC3339)
    // @required
    // @example 2024-04-19T10:00:00Z
    CreatedAt time.Time `json:"created_at" example:"2024-04-19T10:00:00Z"`
}
```

### Documentation Comments for Structs

```go
// CreateUserRequest contains user registration data
// @swagger:model
type CreateUserRequest struct {
    // User email address
    // @required
    // @example user@example.com
    Email string `json:"email" validate:"required,email" example:"user@example.com"`

    // User password
    // @required
    // @minLength 8
    // @maxLength 256
    // @example SecurePassword123!
    Password string `json:"password" validate:"required,min=8,max=256" example:"SecurePassword123!"`

    // User full name
    // @required
    // @minLength 2
    // @maxLength 100
    // @example John Doe
    Name string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
}
```

---

## Swagger UI

### Accessing Swagger UI

After generating docs, access at:

```
http://localhost:8080/api/docs
```

### Features

- 📖 **Interactive Documentation** - See all endpoints
- 🧪 **Try It Out** - Test endpoints directly from browser
- 📥 **Request Preview** - See example requests
- 📤 **Response Preview** - See example responses
- 🔐 **Authentication** - Test with Bearer tokens
- 📊 **Schema Visualization** - See data models

### Download OpenAPI Spec

```
http://localhost:8080/api/docs/swagger.json
```

Use spec file with:
- API client generators
- API monitoring tools
- Documentation generators
- API gateway configurations

---

## Best Practices

### ✅ DO:

```go
// 1. Document all public endpoints
// @Summary Get user by ID
// @Description Retrieve a specific user's information
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
    // ...
}

// 2. Document request parameters
// @Param limit query int false "Number of results (default: 10)" default(10)
// @Param offset query int false "Offset for pagination (default: 0)" default(0)

// 3. Document response models
// @Success 200 {object} UserResponse "User data"
// @Failure 400 {object} ErrorResponse "Validation error"

// 4. Provide examples
// @example (request) {"name": "John", "email": "john@example.com"}

// 5. Document security
// @Security Bearer

// 6. Use clear descriptions
// @Description Retrieve a list of all users with pagination
```

### ❌ DON'T:

```go
// ❌ Minimal documentation
// @Summary User endpoint

// ❌ Missing error responses
// @Success 200 {object} Response

// ❌ No type information
// @Param id query string

// ❌ Unclear descriptions
// @Description Do user stuff

// ❌ Missing examples
// @Success 200 {object} Response
```

---

## Generating Documentation

### Command Line

```bash
# Generate swagger docs
cd backend
swag init -g cmd/server/main.go

# Generate with custom output
swag init -g cmd/server/main.go -o docs/swagger

# View generated files
ls docs/
# Output: docs.go swagger.json swagger.yaml
```

### Makefile Command

Add to Makefile:

```bash
docs:
	@echo "Generating OpenAPI documentation..."
	@cd backend && swag init -g cmd/server/main.go -o docs

docs-view:
	@echo "Opening Swagger UI..."
	@open http://localhost:8080/api/docs
```

### Run:

```bash
make docs
make docs-view
```

---

## Examples

### GET Endpoint (List)

```go
// @Summary List all posts
// @Description Get paginated list of blog posts
// @Tags Posts
// @Produce json
// @Param limit query int false "Number of results" default(10)
// @Param offset query int false "Offset" default(0)
// @Param published query boolean false "Filter by published status"
// @Success 200 {array} PostResponse "List of posts"
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
func (h *PostHandler) List(c *gin.Context) {
    limit := 10
    offset := 0
    
    if l := c.Query("limit"); l != "" {
        limit, _ = strconv.Atoi(l)
    }
    
    posts, err := h.service.GetPosts(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(500, ErrorResponse{Error: "Server error"})
        return
    }
    
    c.JSON(200, posts)
}
```

### POST Endpoint (Create)

```go
// @Summary Create a new post
// @Description Create a blog post with title and content
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "Post data"
// @Success 201 {object} PostResponse "Post created"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /posts [post]
// @Security Bearer
func (h *PostHandler) Create(c *gin.Context) {
    var req CreatePostRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: "Invalid request"})
        return
    }
    
    post, err := h.service.CreatePost(c.Request.Context(), &req)
    if err != nil {
        c.JSON(500, ErrorResponse{Error: "Server error"})
        return
    }
    
    c.JSON(201, post)
}
```

### DELETE Endpoint

```go
// @Summary Delete a post
// @Description Permanently delete a blog post
// @Tags Posts
// @Param id path string true "Post ID"
// @Success 204 "Post deleted"
// @Failure 404 {object} ErrorResponse "Post not found"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /posts/{id} [delete]
// @Security Bearer
func (h *PostHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    
    if err := h.service.DeletePost(c.Request.Context(), id); err != nil {
        c.JSON(500, ErrorResponse{Error: "Server error"})
        return
    }
    
    c.JSON(204, nil)
}
```

### File Upload Endpoint

```go
// @Summary Upload file
// @Description Upload a file (images, documents, etc)
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param description formData string false "File description"
// @Success 200 {object} FileResponse "Upload successful"
// @Failure 400 {object} ErrorResponse "Invalid file"
// @Failure 413 {object} ErrorResponse "File too large"
// @Router /upload [post]
// @Security Bearer
func (h *FileHandler) Upload(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, ErrorResponse{Error: "No file provided"})
        return
    }
    
    // Handle upload...
}
```

---

## Client SDK Generation

### Generate TypeScript Client

```bash
# Install OpenAPI generator
npm install @openapitools/openapi-generator-cli -D

# Generate TypeScript client
npx openapi-generator-cli generate \
  -i http://localhost:8080/api/docs/swagger.json \
  -g typescript-axios \
  -o ./generated/ts-client
```

### Usage

```typescript
// Generated client is type-safe
import { UsersApi } from './generated/ts-client'

const api = new UsersApi()

// TypeScript knows all parameters and return types
const user = await api.getUser({ id: '123' })
console.log(user.email)  // ✅ Type-safe
```

---

## Integration with Documentation

### Link in README

```markdown
## API Documentation

Interactive API documentation available at:
- [Swagger UI](http://localhost:8080/api/docs)
- [OpenAPI Spec](http://localhost:8080/api/docs/swagger.json)
```

### Link in Code Comments

```go
// Handler for /api/users endpoint
// See documentation: GET /api/docs
```

---

## Troubleshooting

### Swagger UI not showing endpoints

1. Check comments are above function (not after)
2. Run `swag init` to regenerate
3. Restart server
4. Check browser cache

### Incorrect response types

1. Verify struct field types match
2. Use struct pointer for objects: `{object} UserResponse`
3. Use array syntax for arrays: `{array} PostResponse`

### Examples not showing

1. Add `// @example` comments
2. Add `example` tags to struct fields
3. Regenerate documentation

---

## Swagger Annotations Cheat Sheet

```go
// Endpoint
// @Summary Summary of endpoint
// @Description Detailed description
// @Tags Category
// @Router /path/{id} [method]

// Input/Output
// @Accept json
// @Produce json
// @Param name in type required "Description"
// @Success 200 {object} Type "Description"
// @Failure 400 {object} ErrorType "Description"

// Security
// @Security Bearer
// @Security BasicAuth
// @Security ApiKey

// Examples
// @Example request {"key": "value"}
// @Example response {"status": "success"}

// Model documentation
// @swagger:model
// type User struct {
//   ID string `json:"id" example:"123"`
// }
```

---

## Resources

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3)
- [Gin Swagger](https://github.com/swaggo/gin-swagger)
- [API Documentation Best Practices](https://swagger.io/blog/api-documentation/)

---

**Next: Generate your API docs with `make docs`! 🚀**
