package db

import (
	"context"
)

type ForgotPasswordTxParams struct {
	UpdateUserVerificationCodeParams
	AfterCreate func(user User) error
}

// CreateTxResult is the result of the transfer transaction
type ForgotPasswordTxResult struct {
	User User
}

func (store *SQLStore) ForgotPasswordTx(ctx context.Context, arg ForgotPasswordTxParams) (ForgotPasswordTxResult, error) {
	var result ForgotPasswordTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.UpdateUserVerificationCode(ctx, arg.UpdateUserVerificationCodeParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})
	return result, err
}