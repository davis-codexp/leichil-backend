package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var fcm_client *messaging.Client

type PushNotificationRequest struct {
	DeviceToken string            `json:"device_token"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data,omitempty"`
}

type MultiplePushNotificationRequest struct {
	DeviceTokens []string          `json:"device_tokens"`
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	Data         map[string]string `json:"data,omitempty"`
}

func InitFirebase() {
	var ctx = context.Background()
	var opt option.ClientOption = option.WithCredentialsFile("fierbase.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println("Unable to initialize Firebase App: ", err.Error())
		return
	}
	fcm_client, err = app.Messaging(ctx)
	if err != nil {
		log.Println("Unable to initialize Firebase Client: ", err.Error())
		return
	}
	log.Println("Firebase Initialized!!")
}

func SendPushNotification(device_token, title, body string, data map[string]string) bool {
	var message = &messaging.Message{
		Token: device_token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}
	message_id, err := fcm_client.Send(context.Background(), message)
	if err != nil {
		log.Println("Error in sending push notification")
		return false
	}
	log.Println("Push Notification Message ID: ", message_id)
	return true
}

func SendMultiplePushNotification(device_tokens []string, title, body string, data map[string]string) bool {
	var message = &messaging.MulticastMessage{
		Tokens: device_tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}
	br, err := fcm_client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Println("Error in sending multicast push notifications: ", err.Error())
		return false
	}
	log.Printf("Messages Sent! Success: %d, Failure: %d", br.SuccessCount, br.FailureCount)
	return true
}

func SendMultipleSilentPushNotification(device_tokens []string, data map[string]string) bool {
	var message = &messaging.MulticastMessage{
		Tokens: device_tokens,
		Data:   data,
	}
	br, err := fcm_client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Println("Error in sending multicast push notifications: ", err.Error())
		return false
	}
	log.Printf("Messages Sent! Success: %d, Failure: %d", br.SuccessCount, br.FailureCount)
	return true
}
