# Kodia Plugin Development Guide 🐨🔌

Welcome to the Kodia Plugin Ecosystem. This guide will teach you how to build and distribute high-performance extensions for the Kodia Fullstack Framework.

## What is a Plugin?

In Kodia, a Plugin is a package that implements the `Plugin` interface. It can:
- Register new services to the central container.
- Add HTTP routes to the application.
- Listen to and dispatch application-wide events (Hooks).
- Register background workers and cleanup tasks.

## The Plugin Interface

Every plugin must implement the following methods:

```go
type Plugin interface {
    Name() string
    Metadata() PluginMetadata
    Register(app *kodia.App) error
    Boot(app *kodia.App) error
}
```

### Metadata

Metadata provides information about your plugin:
```go
func (p *MyPlugin) Metadata() kodia.PluginMetadata {
    return kodia.PluginMetadata{
        ID:      "com.example.audit-logger",
        Name:    "Audit Logger",
        Version: "1.0.0",
        Author:  "John Doe",
        Description: "An institutional-grade audit logger for Kodia.",
    }
}
```

## Hook System (Event Driven)

Hooks allow your plugin to communicate with the kernel and other plugins.

### Listening to Hooks
Plugins should register listeners in their `Boot()` method:

```go
func (p *MyPlugin) Boot(app *kodia.App) error {
    app.Hooks.Listen("user.created", func(data any) {
        user := data.(*models.User)
        app.Log.Info("New user detected by plugin", zap.String("email", user.Email))
    })
    return nil
}
```

### Dispatching Hooks
You can also trigger your own events:

```go
app.Hooks.Dispatch("audit.log.created", logEntry)
```

## Registering Routes

If your plugin provides an API, implement the `RouterProvider` interface:

```go
func (p *MyPlugin) RegisterRoutes(router *gin.Engine, app *kodia.App) error {
    group := router.Group("/api/my-plugin")
    group.GET("/status", p.handleStatus)
    return nil
}
```

## Best Practices

1.  **Prefixing**: Always prefix your container keys and hook names (e.g., `my_plugin.service`).
2.  **Statelessness**: Try to keep your plugin logic stateless, using the `App` container for persistence and shared state.
3.  **Graceful Cleanup**: If your plugin opens connections or starts background tasks, use `app.RegisterCleanupTask()` to ensure they are stopped correctly.

© 2026 Kodia Studio. "Build like a user, code like a pro."
