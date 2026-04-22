# Background Workers

Some tasks in your application may take too long to process during a standard HTTP request (e.g., sending a mass email, processing a large image). Kodia handles these using background workers powered by [Asynq](https://github.com/hibiken/asynq) and backed by **Redis**.

## Core Concepts

-   **Client**: Submits tasks to the queue.
-   **Server (Worker)**: Pulls tasks from the queue and processes them.
-   **Task**: A named unit of work with a JSON payload.

## Defining a Task

Create a task with a unique name and any data it needs to execute. We recommend defining these constants in your service layer.

```go
const TaskWelcomeEmail = "email:welcome"

type WelcomeEmailPayload struct {
    UserID string `json:"user_id"`
}
```

## Enqueueing a Task

To run a task in the background, use the `QueueTask` method in your service. Kodia provides the `asynq.Client` via its port.

```go
payload, _ := json.Marshal(WelcomeEmailPayload{UserID: user.ID})
task := asynq.NewTask(TaskWelcomeEmail, payload)

info, err := s.queue.Enqueue(task)
```

## Handling a Task

You must register a handler for your task. This code will be executed by the background worker process.

```go
func (h *EmailHandler) HandleWelcomeEmail(ctx context.Context, t *asynq.Task) error {
    var p WelcomeEmailPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return err // Asynq will retry if an error is returned
    }
    
    // Logic to send email...
    return nil
}
```

## Registering Handlers

In your `cmd/worker/main.go`, register all your task handlers using a ServeMux:

```go
mux := asynq.NewServeMux()
mux.HandleFunc(TaskWelcomeEmail, emailHandler.HandleWelcomeEmail)

if err := worker.Run(mux); err != nil {
    log.Fatal(err)
}
```

## Advanced Task Management

Kodia provides "Elite" task orchestration features out of the box.

### 1. Scheduled Jobs (Fluent API)
Define cron jobs using a natural language syntax in your console kernel.

```go
// Example: Run a cleanup task every day at 10:00 AM
scheduler.Every(1).Day().At("10:00").Do(TaskCleanup, nil)

// Or use a raw Cron expression
scheduler.Cron("0 0 * * *").Do(TaskBackup, payload)
```

### 2. Job Chaining
Execute a sequence of tasks one after another. If a task fails, the chain stops.

```go
err := s.queue.EnqueueChain(ctx, task1, task2, task3)
```

### 3. Job Batching
Submit a group of tasks to be processed efficiently.

```go
err := s.queue.EnqueueBatch(ctx, []ports.Task{job1, job2, job3})
```

---

## 📊 Monitoring & Failed Jobs

Kodia includes a real-time monitor for your queues.

### Dashboard Integration
Access the administrative dashboard at `/api/admin/queues` to:
- **Monitor**: Real-time throughput and queue depth.
- **Failures**: Inspect failed jobs with full stack traces.
- **Retry**: Re-enqueue failed jobs with a single click.

### Audit Logging (Persistence)
Kodia automatically logs failed jobs to the `failed_jobs` table in PostgreSQL. This ensures you have a permanent record of what went wrong, even if Redis data is expired.

> [!TIP]
> You can create a custom UI to query the `failed_jobs` table for long-term reliability auditing.
