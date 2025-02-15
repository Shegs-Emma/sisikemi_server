package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/techschool/simplebank/db/sqlc"
	mymail "github.com/techschool/simplebank/mail"
	"github.com/techschool/simplebank/util"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type SendVerifyEmailPayload struct {
	Username string `json:"username"`
}

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendVerifyEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", err)
	}

	user, err := processor.store.GetUserByUsername(ctx, payload.Username)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return fmt.Errorf("failed to get user: %w", err)
		}
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username: user.Username,
		Email: user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	// verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	verifyUrl := fmt.Sprintf("%s/verify_email?email_id=%d&secret_code=%s",processor.config.ClientEndpoint ,verifyEmail.ID, verifyEmail.SecretCode)

	msg := mymail.Message{
		Subject: "Welcome to Sisikemi Fashion",
		VerifyUrl: verifyUrl,
		Content: fmt.Sprintf(`Hello %s,<br />
		Thank you for registering with us!<br/ >
		Please <a href="%s">click here</a> to verify your email address.<br/>
		`, user.FirstName, verifyUrl),
		To: user.Email,
		From: "mightymilan04@gmail.com",
		FromName: "Sisikemi Fashion",
	}

	if err := processor.sendSmtpMailer.SendSMTPEmail(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}