package worker

import (
	"context"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/mail"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendVerificationCodeEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store db.Store
	sendSmtpMailer mail.EmailSendSMTPEmailSender
	config util.Config
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, sendSmtpMailer mail.EmailSendSMTPEmailSender, config util.Config) TaskProcessor {
	server := asynq.NewServer(
		redisOpt, 
		asynq.Config{
			Queues: map[string] int{
				QueueCritical: 10,
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store: store,
		sendSmtpMailer: sendSmtpMailer,
		config: config,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(TaskSendVerificationCodeEmail, processor.ProcessTaskSendVerificationCodeEmail)

	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}