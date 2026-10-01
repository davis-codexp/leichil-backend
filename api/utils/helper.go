package utils

import (
	"leichil-backend/message_broker"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const DDMMYYYY = "02-01-2006"
const DDMMYYYYhhmmss = "02-01-2006 15:04:05"
const YYYYMMDD = "2006-01-02"
const YYYYMMDDhhmmss = "2006-01-02 15:04:05"

type OtpMessageBody struct {
	MessageID string `json:"message_id"`
	Mobile    string `json:"mobile"`
	Otp       string `json:"otp"`
}

func GeneratePassword() (string, []byte, error) {
	rand.New(rand.NewSource(time.Now().UnixMilli()))
	var temp = rand.Int()
	var password = strconv.Itoa(temp)[0:4]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return password, hash, err
}

func GenerateOTP() string {
	rand.New(rand.NewSource(time.Now().UnixMilli()))
	var temp = rand.Int()
	var password = strconv.Itoa(temp)[0:4]
	return password
}

func GenerateToken(plain_text []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(plain_text, 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func SendOTP(mobile, otp string) {
	var message_id = os.Getenv("OTP_MESSAGE_ID")
	message_broker.SendSms(message_id, mobile, otp)
}
