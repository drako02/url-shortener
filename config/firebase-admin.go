package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

var FirebaseAuth *auth.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("/home/drako/Downloads/url-shortner-784fc-firebase-adminsdk-fbsvc-c7b0277b29.json")
	var err error
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
}

func InitFirebaseAuth() {
	var err error
	FirebaseAuth, err = FirebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
}

