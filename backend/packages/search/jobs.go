package search

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeSearchIndex  = "search:index"
	TypeSearchDelete = "search:delete"
)

type IndexPayload struct {
	Index string                 `json:"index"`
	ID    string                 `json:"id"`
	Data  map[string]interface{} `json:"data"`
}

type DeletePayload struct {
	Index string `json:"index"`
	ID    string `json:"id"`
}

/**
 * HandleSearchIndexTask processes the background indexing task.
 */
func HandleSearchIndexTask(ctx context.Context, t *asynq.Task, manager *SearchManager) error {
	var p IndexPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return manager.Index(p.Index, p.ID, p.Data)
}

/**
 * HandleSearchDeleteTask processes the background deletion task.
 */
func HandleSearchDeleteTask(ctx context.Context, t *asynq.Task, manager *SearchManager) error {
	var p DeletePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return manager.Delete(p.Index, p.ID)
}

/**
 * DispatchIndexTask helper to queue an indexing job.
 */
func DispatchIndexTask(client *asynq.Client, index string, id string, data map[string]interface{}) error {
	payload, _ := json.Marshal(IndexPayload{Index: index, ID: id, Data: data})
	task := asynq.NewTask(TypeSearchIndex, payload)
	_, err := client.Enqueue(task)
	return err
}

/**
 * DispatchDeleteTask helper to queue a deletion job.
 */
func DispatchDeleteTask(client *asynq.Client, index string, id string) error {
	payload, _ := json.Marshal(DeletePayload{Index: index, ID: id})
	task := asynq.NewTask(TypeSearchDelete, payload)
	_, err := client.Enqueue(task)
	return err
}
