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
	DeviceID *string            `json:"device_id"`
	Title    *string            `json:"title"`
	Body     *string            `json:"body"`
	Data     *map[string]string `json:"data,omitempty"`
}

type MultiplePushNotificationRequest struct {
	DeviceIDs *[]string          `json:"device_ids"`
	Title     *string            `json:"title"`
	Body      *string            `json:"body"`
	Data      *map[string]string `json:"data,omitempty"`
}

func InitFirebase() {
	var ctx = context.Background()
	var opt option.ClientOption = option.WithCredentialsFile("firebase.json")
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

func SendPushNotification(device_id, title, body string, data map[string]string) bool {
	var message = &messaging.Message{
		Token: device_id,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}
	message_id, err := fcm_client.Send(context.Background(), message)
	if err != nil {
		log.Println("Error in sending push notification", err.Error())
		return false
	}
	log.Println("Push Notification Message ID: ", message_id)
	return true
}

func SendMultiplePushNotification(device_ids []string, title, body string, data map[string]string) bool {
	var message = &messaging.MulticastMessage{
		Tokens: device_ids,
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

func SendMultipleSilentPushNotification(device_ids []string, data map[string]string) bool {
	var message = &messaging.MulticastMessage{
		Tokens: device_ids,
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
