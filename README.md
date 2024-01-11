<p align="left">
  
  #### 1- Use `pgx` as the PostgreSQL driver:
  In this project, I used pgx to connect gin API with PostgreSQL to perform some transactions like INSERT, UPDATE, DELETE elsements from tables.
  
  #### 2- Use transactions for queries in Golang:
  Transactions are units of work that consist of multiple database operations that are executed as a single, atomic, and consistent unit. Transactions ensure data      integrity and the ACID properties: 
  Atomicity, Consistency, Isolation, and Durability.
      BEGIN TRANSACTION or START TRANSACTION: Starts a new transaction.
      COMMIT: Commits the transaction, making all its changes permanent.
      ROLLBACK: Rolls back (undoes) the transaction, canceling all its changes.

  #### 3- Utilize `sqlc` to generate type safe Go code from SQL queries:
  In this project, I used "SQLC" tool, and It's designed to generate type-safe Go code based on SQL queries and schema definitions.
  
  #### 4- Implement functions to:
        Create a new user with proper data validation.   ```POST /api/users```
        Generate a new OTP for a user.                   ```POST /api/users/generateotp```
        Verify OTP for a user.                           ```POST /api/users/verifyotp```
</p>



`POST /api/users`: Create a new user.
  1- Accepts JSON payload with `name` and `phone_number`.
  2- Ensure that `phone_number` is unique; if not, return a 400 error.
  3- Store the user in the database.

![Screenshot 2024-01-11 062912](https://github.com/Stevenwaheed/Golang_Postgres_API_Test/assets/83607748/e0186d9b-13dc-4f6f-9a9c-12e3e80b0407)

<p align="center">
  <img src="https://github.com/Stevenwaheed/Golang_Postgres_API_Test/assets/83607748/1e1b7ff5-1522-45e1-a165-7d7f8f6fa566"/>
</p>

`POST /api/users/generateotp`: Generate OTP for a user.
  1- Accepts JSON payload with `phone_number`.
  2- If the `phone_number` does not exist, return a 404 error.
  3- Generate a random 4 digit OTP and set its expiration time to 1 minute from the current time.

![Screenshot 2024-01-11 063110](https://github.com/Stevenwaheed/Golang_Postgres_API_Test/assets/83607748/b9e06c1e-af57-473b-ba67-a5332ae23bd4)

<p align="center">
  <img src="https://github.com/Stevenwaheed/Golang_Postgres_API_Test/assets/83607748/10205f45-9332-4684-ad4e-8361937f103d"/>
</p>


`POST /api/users/verifyotp`: Verify OTP for a user.
  1- Accepts JSON payload with `phone_number` and `otp`.
  2- Check if the OTP is correct and not expired (compare with `otp_expiration_time`).
  3- If the OTP is correct and not expired, return a success message.
  4- If the OTP is incorrect, return an error message.
  5- If the OTP is expired, return an error message indicating that the OTP has expired.


```
  (Verified)
```
![Screenshot 2024-01-11 063203](https://github.com/Stevenwaheed/Golang_Postgres_API_Test/assets/83607748/64462ad5-89b6-4ea6-b51e-e505f81eee14)


```
  (Expired)
```
![Screenshot 2024-01-11 063225](https://github.com/Stevenwaheed/Golang_Postgres_API_Test/assets/83607748/f6f45693-09a5-4040-9106-b6161572cc1d)


