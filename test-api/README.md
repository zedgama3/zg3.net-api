# zg3.net-api
API endpoints for zg3.net

# Notes:

Enable Go telemetry:
`go run golang.org/x/telemetry/cmd/gotelemetry@latest on`

Cleanup dependancies?
`go mod tidy`

Run Application:
```PowerShell
$env:GIN_MODE = "release" # Omit during development.  Can also be set to: debug, release, test - defaults to debug.
go run ./cmd/server/main.go
```

**Purpose**
Create API endpoints for the zg3.net website.

**Tooling**
- Golang
- GIN
- Insomnia

# Links

[Project Layout Example](https://github.com/golang-standards/project-layout)
[Using GraphQL with Golang](https://www.apollographql.com/blog/using-graphql-with-golang)

Example provided by ChatGPT:
```
/myapi
│
├── cmd
│   └── server
│       └── main.go            # Contains the main function, starts the server
│
├── internal
│   ├── app                    # Application logic
│   │   ├── handler            # Request handlers
│   │   │   ├── user.go        # Handlers for user-related requests
│   │   │   └── product.go     # Handlers for product-related requests
│   │   └── middleware         # Middleware components
│   │       └── auth.go        # Authentication middleware
│   │
│   ├── model                  # Data models
│   │   ├── user.go            # User model
│   │   └── product.go         # Product model
│   │
│   └── service                # Business logic and service layer
│       ├── userService.go     # Business logic for user handling
│       └── productService.go  # Business logic for product handling
│
├── pkg                        # Library code that's ok to use by external applications
│   └── utils                  # Utility functions and common libraries
│       └── logger.go          # Logging utility
│
├── config                     # Configuration files and code
│   └── config.go              # Configuration setup and handling
│
├── go.mod                     # Go module file
└── go.sum                     # Go sum file for module verification
```

[Example of API](https://github.com/benbjohnson/wtf)