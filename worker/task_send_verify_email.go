package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

const (
	TASK_SEND_VERIFY_EMAIL = "task:send_verify_email"
)

func (rt *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TASK_SEND_VERIFY_EMAIL, jsPayload, opts...)
	info, err := rt.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"task_id":              info.ID,
		"task_type":            info.Type,
		"task_queue":           info.Queue,
		"task_state":           info.State,
		"task_payload":         string(info.Payload),
		"task_max_retry":       info.MaxRetry,
		"task_next_process_at": info.NextProcessAt,
	}).Info("enqueued task")

	return nil
}
