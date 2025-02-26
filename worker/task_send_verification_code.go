package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	db "github.com/techschool/simplebank/db/sqlc"
	mymail "github.com/techschool/simplebank/mail"
)

const TaskSendVerificationCodeEmail = "task:send_verification_code_email"

type SendVerificationCodePayload struct {
	Username string `json:"username"`
}

type PayloadSendVerificationCodeEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerificationCodeEmail(
	ctx context.Context,
	payload *PayloadSendVerificationCodeEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerificationCodeEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerificationCodeEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendVerificationCodePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", err)
	}

	user, err := processor.store.GetUserByUsername(ctx, payload.Username)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return fmt.Errorf("failed to get user: %w", err)
		}
	}

	verificationCodeEmail, err := processor.store.UpdateUserVerificationCode(ctx, db.UpdateUserVerificationCodeParams{
		Email: user.Email,
		VerificationCode: pgtype.Text{
			String: user.VerificationCode.String,
			Valid: true,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create verification code email: %w", err)
	}

	msg := mymail.Message{
		Subject: "Forgot Password Code",
		VerifyUrl: "",
		Content: fmt.Sprintf(`Hello %s,<br />
		Please use this code for password reset!<br/ >
		Verification Code: %s
		`, user.FirstName, verificationCodeEmail.VerificationCode.String),
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