# Antrea Renovate Demo - Secure API Server

This repository demonstrates the use of Renovate for vulnerability-focused dependency management in Go projects. It implements a simple RESTful API server with an intentionally vulnerable dependency to showcase how Renovate can be configured to update only dependencies with security vulnerabilities.

## What it does

This application implements a RESTful API server with:

- JWT authentication for securing endpoints
- In-memory storage for item management
- CRUD operations for items via REST endpoints
- Health check endpoint

The server is built using:
- `github.com/gin-gonic/gin` (v1.6.3) - A web framework for building APIs, **with a known security vulnerability**
- `github.com/golang-jwt/jwt/v4` - For JWT authentication

## Project Structure

```
antrea-renovate-demo/
├── cmd/
│   └── apiserver/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── handlers.go
│   ├── auth/
│   │   └── jwt.go
│   └── storage/
│       └── storage.go
├── go.mod
├── go.sum
├── README.md
└── renovate.json
```

## How to run

1. Clone the repository:
```bash
git clone https://github.com/akshayw1/antrea-renovate-demo.git
cd antrea-renovate-demo
```

2. Run the server:
```bash
go run cmd/apiserver/main.go
```

The server runs on port 8080 by default. You can customize this by setting the `PORT` environment variable.

## API Endpoints

### Authentication
- `POST /login` - Get a JWT token
  ```json
  {
    "username": "your_username",
    "password": "your_password"
  }
  ```

### Protected Endpoints (require Authorization header with Bearer token)
- `GET /api/items` - List all items
- `GET /api/items/:id` - Get a specific item 
- `POST /api/items` - Create a new item
  ```json
  {
    "id": "item1",
    "name": "Test Item",
    "value": "Test Value"
  }
  ```
- `DELETE /api/items/:id` - Delete an item

### Health Check
- `GET /health` - Simple health check endpoint

## Vulnerable Dependency

This project intentionally uses a vulnerable version of `github.com/gin-gonic/gin` (v1.6.3) which is affected by:

- **CVE-2020-28483**: A denial-of-service vulnerability in the Gin framework
  - Link: https://nvd.nist.gov/vuln/detail/CVE-2020-28483
  - Description: In the Gin framework before v1.6.3, an attacker can cause a denial of service by sending specially crafted URL paths (with lots of slashes) to cause the router to consume excessive CPU and memory for path matching.

## Renovate Configuration

The key to this demo is the Renovate configuration in `renovate.json`:

```json
{
  "extends": ["config:base"],
  "packageRules": [
    {
      "matchManagers": ["gomod"],
      "vulnerabilityAlerts": true,
      "updateTypes": ["patch", "minor", "major"],
      "groupName": "Go Vulnerability Updates"
    }
  ]
}
```

This configuration:

- Only updates Go dependencies when security vulnerabilities are detected (`vulnerabilityAlerts: true`)
- Groups all vulnerability-related updates into a single PR (`groupName: "Go Vulnerability Updates"`)
- Works with all update types: patch, minor, and major
- Uses Renovate's base configuration for sensible defaults

This approach demonstrates how to maintain secure dependencies in active release branches without being overwhelmed by non-security-related updates.

## Solution for Antrea's Requirements

This configuration addresses Antrea's requirements from issue #7155:

1.  **Selective updates**: Only dependencies with security vulnerabilities get updated
2. **Branch-specific rules**: Can be extended with `matchBaseBranches` to apply different rules to release branches
3.  **Reduced PR noise**: Groups related security updates into single PRs
4.  **Easy configuration**: Simple, declarative config in standard JSON format

When Renovate runs with this configuration, it automatically creates a PR that updates the vulnerable Gin version to a secure one, without updating any other dependencies that don't have security issues.

## Results

After configuring Renovate and adding a vulnerable dependency:

1. Renovate detected the vulnerable version of Gin (v1.6.3)
2. It created a single PR named "Update Go Vulnerability Updates" https://github.com/akshayw1/antrea-renovate-demo/pull/3
3. The PR updated Gin to a secure version (v1.10.1)
4. The PR also updated related transitive dependencies to maintain compatibility

This demonstrates Renovate's ability to focus only on security-critical updates, which is perfect for active release branches where you want to minimize changes while ensuring security.