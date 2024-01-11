CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		phone_number TEXT NOT NULL UNIQUE
	);
	
CREATE TABLE users_otp (
    phone_number_otp TEXT NOT NULL REFERENCES users(phone_number),
    otp TEXT,
    otp_expiration_time TIMESTAMP
);