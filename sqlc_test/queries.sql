-- name: newUserTransaction :exec
    INSERT INTO users("name", "phone_number") VALUES ($1, $2);

-- name: otpTransaction :exec
    INSERT INTO users_otp("otp", "otp_expiration_time", "phone_number_otp") VALUES ($1, $2, $3);

-- name: verifyOTP :one
    SELECT otp, otp_expiration_time FROM users_otp WHERE phone_number_otp=$1 AND otp=$2;
