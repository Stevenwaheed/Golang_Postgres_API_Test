package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"context"
	"postgres"
	_ "github.com/lib/pq"
	"database/sql"
)


func createConnection() *sql.DB{
	
	connStr := "user=postgres password=S1122001:) dbname=postgres sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return conn
}


func newUserTransactionMain(c * gin.Context){

	name, _ := c.GetQuery("name")
	phoneNumber, _ := c.GetQuery("phone_number")
	
	conn := createConnection()
	defer conn.Close()

	tx, _ := conn.Begin()
	db := postgres.New(conn)

	err := db.NewUserTransaction(context.Background(), postgres.NewUserTransactionParams{
		Name:    name ,
		PhoneNumber:    phoneNumber,
	})
	defer tx.Rollback()

	if err != nil{
		panic(err)
	}

	c.IndentedJSON(200, gin.H{"Name":name, "PhoneNumber": phoneNumber, "Result":"Added Successfully"})
	tx.Commit()
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


func otpTransactionMain(c * gin.Context){

	phoneNumber, _ := c.GetQuery("phone_number")
	otp := generateOTP()

	conn := createConnection()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Something Goes Wrong"})
		panic(err)
	}
	   
	db := postgres.New(conn)

	err = db.OtpTransaction(context.Background(), postgres.OtpTransactionParams{
		PhoneNumberOtp:    	phoneNumber,
		Otp: 				otp,
		OtpExpirationTime:  time.Now(),
	})
	defer tx.Rollback()   // Rollback after this function commit the result to make sure that there is no other operations

	tx.Commit()

	if err != nil{
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Connection is busy or closed"})
		panic(err)
	} else{
		c.IndentedJSON(http.StatusOK, gin.H{"Phone Number":phoneNumber, "Result":"Added Successfully, your OTP is "+ otp +", and it will expire after 1 min"})
	}
}


func verifyOTPMain(c * gin.Context){

	phoneNumber, _ := c.GetQuery("phone_number")
	otp, _ := c.GetQuery("otp")

	conn := createConnection()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Something Goes Wrong"})
		panic(err)
	}
	   
	db := postgres.New(conn)

	rows, err := db.VerifyOTP(context.Background(), postgres.VerifyOTPParams{
		PhoneNumberOtp:    	phoneNumber,
		Otp: 				otp,
	})
	defer tx.Rollback()   // Rollback after this function commit the result to make sure that there is no other operations

	
	if time.Now().Minute() - rows.OtpExpirationTime.Minute() > 1{
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Expired"})
	} else {
		if otp == rows.Otp{
			c.IndentedJSON(http.StatusOK, gin.H{"Phone Number":phoneNumber, "Result":"Verifed", "Now":time.Now().Minute() % 2, "stored":rows.OtpExpirationTime.Minute()%2})
		} else if otp != rows.Otp{
			c.IndentedJSON(http.StatusNotFound, gin.H{"Phone Number":phoneNumber, "OTP Result":"Not Found"})
		}
	}

	err = tx.Commit()
	if err != nil{
		c.IndentedJSON(http.StatusForbidden, gin.H{"Phone Number":phoneNumber, "Result":"Connection is busy or closed"})
		panic(err)
	}
}



func main(){
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.POST("/api/users", newUserTransactionMain)
	router.POST("/api/users/generateotp", otpTransactionMain)
	router.POST("/api/users/verifyotp", verifyOTPMain)
	router.Run("localhost:5000")	
}

