package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/NguyenMinhKhanhBK/simple_bank/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	conf := asynq.Config{
		Queues: map[string]int{
			QueueCritical: 6,
			QueueDefault:  3,
		},
	}
	server := asynq.NewServer(redisOpt, conf)
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (rp *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload with error %w, stop retry: %w", err, asynq.SkipRetry)
	}

	user, err := rp.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: send email to user
	logrus.WithField("user", user.Username).WithField("task_type", task.Type()).Info("processed task")

	return nil
}

func (rp *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TASK_SEND_VERIFY_EMAIL, rp.ProcessTaskSendVerifyEmail)
	return rp.server.Start(mux)
}
