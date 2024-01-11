package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
)


type User struct{
	Name  				 string	   `json:"name"`
	PhoneNumber  		 string    `json:"phone_number"`
}


type OTP struct{
	OTP   				 string               `json:"otp"`
	OTPExpirationTime    time.Time       `json:"otp_expration_time"`
}


var dbURL = "postgres://postgres:S1122001:)@localhost:5432/postgres"


func openConnection() *pgx.Conn{
	// Open a connection with the database
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		panic(err)
	}
	
	conn, err := pgx.ConnectConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}
	
	return conn
}


func createTable(c * gin.Context){
	var createTables = `
    CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		phone_number TEXT NOT NULL UNIQUE
	)
	
	CREATE TABLE IF NOT EXISTS users_otp (
		phone_number_otp TEXT NOT NULL REFERENCES users(phone_number),
		otp TEXT,
		otp_expiration_time TIMESTAMP
	)`

	conn := openConnection()
	conn.Query(context.Background(), createTables)

	c.IndentedJSON(http.StatusOK, "Tables created successfully.")
}



func setUserInfo(userName, phoneNumber string)User{
	var user User

	user.Name = userName
	user.PhoneNumber = phoneNumber

	return user
}



func setOTPInfo(OTP_var string, OTPExpirationTime time.Time) OTP{
	var otp OTP

	otp.OTP = OTP_var
	otp.OTPExpirationTime = OTPExpirationTime

	return otp
}



func checkPhoneNumberDuplicate(phoneNumber string, phoneNumbers pgx.Rows) bool{
	for _, phoneNum := range phoneNumbers.RawValues(){
		if phoneNumber == string(phoneNum){
			return true
		}
	}
	return false
}


func newUserTransaction(conn *pgx.Conn, userName, phoneNumber string) error{
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())     // Rollback after this function commit the result to make sure that there is no other operations

	var insertCommend = "INSERT INTO users(name, phone_number) VALUES ($1, $2)"
	_, err = tx.Exec(context.Background(), insertCommend, userName, phoneNumber)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil{
		return err
	}

	return nil
}


func createNewUser(c * gin.Context){
	var user User

	name, _ := c.GetQuery("name")
	phoneNumber, _ := c.GetQuery("phone_number")
	user = setUserInfo(name, phoneNumber)

	conn := openConnection()
	err := newUserTransaction(conn, user.Name, user.PhoneNumber)

	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Connection is busy or closed"})
		panic(err)
	} else{
		c.IndentedJSON(http.StatusOK, gin.H{"Name":name, "Phone Number":phoneNumber, "Result":"Added Successfully"})
	}

	defer conn.Close(context.Background())    // Close a connection in the end of the program
}



func generateOTP()string {
	var otp []int
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 4; i++{
		otp = append(otp, rand.Intn(10))
	}

	str := ""
	for i:=0; i<len(otp); i++{
		str += strconv.Itoa(otp[i])
	}
	
	return str
}


func otpTransaction(conn *pgx.Conn, phoneNumber string, otp OTP) error{
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())    // Rollback after this function commit the result to make sure that there is no other operations

	var insertOTPCommend = "INSERT INTO users_otp(otp, otp_expiration_time, phone_number_otp) VALUES ($1, $2, $3)"
	_, err = tx.Exec(context.Background(), insertOTPCommend, otp.OTP, otp.OTPExpirationTime, phoneNumber)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil{
		return err
	}

	return nil
}


func insertOTP(c * gin.Context){
	phoneNumber, _ := c.GetQuery("phone_number")
	otp := generateOTP()
	
	// Open a connection with the database
	conn := openConnection()
	defer conn.Close(context.Background()) // Close a connection in the end of the program
	
	otp_obj := setOTPInfo(otp, time.Now())
	err := otpTransaction(conn, phoneNumber, otp_obj)
	
	if err != nil{
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Connection is busy or closed"})
		panic(err)
	} else{
		c.IndentedJSON(http.StatusOK, gin.H{"Phone Number":phoneNumber, "Result":"Added Successfully, your OTP is "+ otp +", and it will expire after 1 min"})
	}
}


func verifyOTP(c * gin.Context){
	phoneNumber, _ := c.GetQuery("phone_number")
	otp, _ := c.GetQuery("otp")

	// Open a connection with the database
	conn := openConnection()
	defer conn.Close(context.Background())  // Close a connection in the end of the program

	tx, err := conn.Begin(context.Background())
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Something Goes Wrong"})
		panic(err)
	}
	defer tx.Rollback(context.Background())    // Rollback after this function commit the result to make sure that there is no other operations

	// Database select commend
	selectCommand := "SELECT otp, otp_expiration_time FROM users_otp WHERE phone_number_otp=$1 AND otp=$2"
	rows, err := conn.Query(context.Background(), selectCommand, phoneNumber, otp)
	if err != nil{
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Something Goes Wrong"})
		panic(err)
	}
	defer rows.Close()


	for rows.Next(){
		var val_otp string
		var otpExpirationTime time.Time
		
		err1 := rows.Scan(&val_otp, &otpExpirationTime)
		if err1 != nil{
			c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Error While Scanning"})
		}

		if time.Now().Minute() - otpExpirationTime.Minute() > 1{
			c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Expired"})
		} else {
			if otp == val_otp{
				c.IndentedJSON(http.StatusOK, gin.H{"Phone Number":phoneNumber, "Result":"Verifed"})
			} else if otp != val_otp{
				c.IndentedJSON(http.StatusNotFound, gin.H{"Phone Number":phoneNumber, "OTP Result":"Not Found"})
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil{
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Connection is busy or closed"})
		panic(err)
	}
}



func main(){
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.POST("/api/createtable", createTable)
	router.POST("/api/users", createNewUser)
	router.POST("/api/users/generateotp", insertOTP)
	router.POST("/api/users/verifyotp", verifyOTP)
	router.Run("localhost:5000")	
}