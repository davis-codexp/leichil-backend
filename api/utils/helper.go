package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"leichil-backend/message_broker"
	"math/rand"
	"os"
	"strconv"
	"time"

	"crypto/aes"
	"crypto/cipher"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"

	rnd "crypto/rand"
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

func Encrypt(plain_text string) (string, error) {
	var cipher_text string
	var key, _ = base64.StdEncoding.DecodeString(os.Getenv("AES_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("Error in generating Cipher: ", err.Error())
		return cipher_text, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error("Error in generating GCM: ", err.Error())
		return cipher_text, err
	}
	var nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rnd.Reader, nonce); err != nil {
		return cipher_text, err
	}
	var cipher []byte = gcm.Seal(nonce, nonce, []byte(plain_text), nil)
	cipher_text = base64.URLEncoding.EncodeToString(cipher)
	log.Info("CIPHER TEXT: ", cipher_text)
	return cipher_text, err
}

func Decrypt(cipher_text string) (string, error) {
	var plain_text string
	data, err := base64.URLEncoding.DecodeString(cipher_text)
	if err != nil {
		log.Error("Error in Decoding Base64 string: ", err.Error())
		return plain_text, err
	}
	var key, _ = base64.StdEncoding.DecodeString(os.Getenv("AES_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("Error in generating Cipher: ", err.Error())
		return plain_text, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error("Error in generating GCM: ", err.Error())
		return plain_text, err
	}
	var nonce_size = gcm.NonceSize()
	if len(data) < nonce_size {
		return plain_text, fmt.Errorf("cipher text too short")
	}
	nonce, cipher := data[:nonce_size], data[nonce_size:]
	text_byte, err := gcm.Open(nil, nonce, cipher, nil)
	if err != nil {
		log.Error("Error in gcm Open: ", err.Error())
		return plain_text, err
	}
	plain_text = string(text_byte)
	return plain_text, nil
}
