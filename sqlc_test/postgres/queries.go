// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package postgres

import (
	"context"
	// "database/sql"
	"time"
)

const newUserTransaction = `-- name: newUserTransaction :exec
    INSERT INTO users("name", "phone_number") VALUES ($1, $2)
`

type NewUserTransactionParams struct {
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phone_number"`
}

func (q *Queries) NewUserTransaction(ctx context.Context, arg NewUserTransactionParams) error {
	_, err := q.db.ExecContext(ctx, newUserTransaction, arg.Name, arg.PhoneNumber)
	return err
}

const otpTransaction = `-- name: otpTransaction :exec
    INSERT INTO users_otp("otp", "otp_expiration_time", "phone_number_otp") VALUES ($1, $2, $3)
`

type OtpTransactionParams struct {
	Otp               string         `json:"otp"`
	OtpExpirationTime time.Time      `json:"otp_expiration_time"`
	PhoneNumberOtp    string         `json:"phone_number_otp"`
}

func (q *Queries) OtpTransaction(ctx context.Context, arg OtpTransactionParams) error {
	_, err := q.db.ExecContext(ctx, otpTransaction, arg.Otp, arg.OtpExpirationTime, arg.PhoneNumberOtp)
	return err
}

const verifyOTP = `-- name: verifyOTP :one
    SELECT otp, otp_expiration_time FROM users_otp WHERE phone_number_otp=$1 AND otp=$2
`

type VerifyOTPParams struct {
	PhoneNumberOtp string         `json:"phone_number_otp"`
	Otp            string         `json:"otp"`
}

type VerifyOTPRow struct {
	Otp               string         `json:"otp"`
	OtpExpirationTime time.Time      `json:"otp_expiration_time"`
}

func (q *Queries) VerifyOTP(ctx context.Context, arg VerifyOTPParams) (VerifyOTPRow, error) {
	row := q.db.QueryRowContext(ctx, verifyOTP, arg.PhoneNumberOtp, arg.Otp)
	var i VerifyOTPRow
	err := row.Scan(&i.Otp, &i.OtpExpirationTime)
	return i, err
}
