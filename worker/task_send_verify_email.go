package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	mymail "github.com/Shegs-Emma/sisikemi_server/mail"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
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

// func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
// 	var payload SendVerifyEmailPayload
// 	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
// 		return fmt.Errorf("failed to unmarshal task payload: %w", err)
// 	}

// 	user, err := processor.store.GetUserByUsername(ctx, payload.Username)

// 	if err != nil {
// 		if errors.Is(err, db.ErrRecordNotFound) {
// 			return fmt.Errorf("failed to get user: %w", err)
// 		}
// 	}

// 	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
// 		Username: user.Username,
// 		Email: user.Email,
// 		SecretCode: util.RandomString(32),
// 	})

// 	if err != nil {
// 		return fmt.Errorf("failed to create verify email: %w", err)
// 	}

// 	// verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
// 	verifyUrl := fmt.Sprintf("%s/verify_email?email_id=%d&secret_code=%s",processor.config.ClientEndpoint ,verifyEmail.ID, verifyEmail.SecretCode)

// 	msg := mymail.Message{
// 		Subject: "Welcome to Sisikemi Fashion",
// 		VerifyUrl: verifyUrl,
// 		Content: fmt.Sprintf(`Hello %s,<br />
// 		Thank you for registering with us!<br/ >
// 		Please <a href="%s">click here</a> to verify your email address.<br/>
// 		`, user.FirstName, verifyUrl),
// 		To: user.Email,
// 		From: processor.config.FromAddress,
// 		FromName: processor.config.FromName,
// 	}

// 	if err := processor.sendSmtpMailer.SendSMTPEmail(msg); err != nil {
// 		return fmt.Errorf("failed to send email: %w", err)
// 	}

// 	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
// 		Str("email", user.Email).Msg("processed task")
// 	return nil
// }

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
		return fmt.Errorf("failed to query user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	// Safely trim trailing slash from ClientEndpoint
	baseURL := processor.config.ClientEndpoint
	verifyUrl := fmt.Sprintf("%sverify_email?email_id=%d&secret_code=%s",
		baseURL, verifyEmail.ID, verifyEmail.SecretCode,
	)

	// Prepare the email message
	// msg := mymail.Message{
	// 	Subject:   "Welcome to Sisikemi Fashion",
	// 	VerifyUrl: verifyUrl, // plain text fallback
	// 	Content: fmt.Sprintf(`
	// 		<!DOCTYPE html>
	// 		<html>
	// 		<body>
	// 			<p>Hello %s,</p>
	// 			<p>Thank you for registering with us!</p>
	// 			<p>Please <a href="%s">click here</a> to verify your email address.</p>
	// 			<p>If the link doesn't work, copy and paste this URL into your browser:</p>
	// 			<p>%s</p>
	// 		</body>
	// 		</html>`,
	// 		user.FirstName, verifyUrl, verifyUrl,
	// 	),
	// 	To:       user.Email,
	// 	From:     processor.config.FromAddress,
	// 	FromName: processor.config.FromName,
	// }

	msg := mymail.Message{
		Subject:   "Welcome to Sisikemi Fashion",
		VerifyUrl: verifyUrl, // plain text fallback
		Content: fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
				<title>Email Verification</title>
			</head>
			<body>
				<p>Hello %s,</p>
				<p>Thank you for registering with us!</p>
				<p>
				Please <a href="%s">click here</a> to verify your email address.
				</p>
				<p>
				If the link doesn’t work, copy and paste this URL into your browser:
				</p>
				<p><a href="%s">%s</a></p>
			</body>
			</html>
		`, user.FirstName, verifyUrl, verifyUrl, verifyUrl),
		To:       user.Email,
		From:     processor.config.FromAddress,
		FromName: processor.config.FromName,
	}

	// Send the email
	if err := processor.sendSmtpMailer.SendSMTPEmail(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")

	return nil
}
