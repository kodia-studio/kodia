# Kodia Architecture: The Service Provider System

Kodia uses a modular "Kernel" architecture inspired by modern framework design. This approach balances a "batteries-included" experience with high extensibility.

## The Kodia Kernel (`App`)

The heart of a Kodia application is the `App` struct (located in `pkg/kodia`). It manages the lifecycle of the application and acts as a central registry for all features.

### Application Lifecycle

1.  **Registering**: Providers are registered. Each provider's `Register()` method is called to initialize services and bind dependencies to the app's internal container.
2.  **Booting**: After all providers are registered, the `Boot()` method of each provider is called. This is where routes and event listeners are registered, ensuring all services they depend on are already available.
3.  **Running**: The application starts its HTTP server and listens for incoming requests.
4.  **Shutdown**: The application handles graceful shutdown with a configurable timeout.

## Service Providers

A Service Provider is a piece of code that "teaches" Kodia a new feature. It implements the `ServiceProvider` interface:

```go
type ServiceProvider interface {
    Name() string
    Register(app *App) error
    Boot(app *App) error
}
```

### Official (Built-in) Providers
Kodia comes pre-packaged with several official providers to ensure you're ready to go from day one:
- **`DatabaseProvider`**: Sets up GORM and connection pooling.
- **`AuthProvider`**: Provides the entire JWT authentication system.
- **`HttpProvider`**: Configures Gin, Security, and CORS.
- **`InfraProvider`**: Handles Mail, Storage, and Background Workers.

## Dependency Injection

Kodia uses a simple, thread-safe internal container to share services between providers. You can store anything in the container during registration and retrieve it during booting.

```go
// Storing a service
app.Set("my_service", services.NewMyService())

// Retrieving a service (with automatic panic if missing)
service := app.MustGet("my_service").(ports.MyService)
```

## Extending with Third-Party Plugins

To create your own plugin:
1.  Implement the `ServiceProvider` interface in your package.
2.  Register your provider in `cmd/server/main.go`.
3.  That's it! Your plugin can now add routes, workers, and events to any Kodia application.
