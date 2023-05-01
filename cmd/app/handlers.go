package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cuongld2.com/api/pkg/models"
	"github.com/gorilla/mux"
	"solace.dev/go/messaging"
	"solace.dev/go/messaging/pkg/solace/config"
	"solace.dev/go/messaging/pkg/solace/resource"
)

const TopicPrefix = "services/meetings"

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Get all meetings stored
	meetings, err := app.meetings.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Convert meeting list into json encoding
	b, err := json.Marshal(meetings)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("meetings have been listed")

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	// Find meeting by id
	m, err := app.meetings.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("meeting not found")
			return
		}
		// Any other error will send an internal server error
		app.serverError(w, err)
	}

	// Convert meeting to json encoding
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Have been found a meeting")

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Define meeting model
	var m models.Meeting
	// Get request information
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	// Insert new meeting
	insertResult, err := app.meetings.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	// Send message to Solace broker

	// Configuration parameters
	brokerConfig := config.ServicePropertyMap{
		config.TransportLayerPropertyHost:                os.Getenv("TransportLayerPropertyHost"),
		config.ServicePropertyVPNName:                    os.Getenv("ServicePropertyVPNName"),
		config.AuthenticationPropertySchemeBasicUserName: os.Getenv("AuthenticationPropertySchemeBasicUserName"),
		config.AuthenticationPropertySchemeBasicPassword: os.Getenv("AuthenticationPropertySchemeBasicPassword"),
	}
	messagingService, err := messaging.NewMessagingServiceBuilder().FromConfigurationProvider(brokerConfig).WithTransportSecurityStrategy(config.NewTransportSecurityStrategy().WithoutCertificateValidation()).
		Build()

	if err != nil {
		panic(err)
	}

	// Connect to the messaging serice
	if err := messagingService.Connect(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the broker? ", messagingService.IsConnected())

	//  Build a Direct Message Publisher
	directPublisher, builderErr := messagingService.CreateDirectMessagePublisherBuilder().Build()
	if builderErr != nil {
		panic(builderErr)
	}

	startErr := directPublisher.Start()
	if startErr != nil {
		panic(startErr)
	}

	fmt.Println("Direct Publisher running? ", directPublisher.IsRunning())

	//  Prepare outbound message payload and body
	messageBody := "New meeting has been created by user is: " + m.SenderUserID + " with category is: " + m.MeetingCategory + " and title is: " + m.MeetingTitle
	messageBuilder := messagingService.MessageBuilder().
		WithProperty("application", "meetings").
		WithProperty("language", "go").
		WithProperty("category", m.MeetingCategory)

	println("Subscribe to topic ", TopicPrefix+"/>")

	if directPublisher.IsReady() {
		message, err := messageBuilder.BuildWithStringPayload(messageBody)
		if err != nil {
			panic(err)
		}
		publishErr := directPublisher.Publish(message, resource.TopicOf(TopicPrefix+"/senderUserId/"+m.SenderUserID+"/meetingCategory/"+m.MeetingCategory+"/meetingTitle/"+m.MeetingTitle+"/"))
		if publishErr != nil {
			panic(publishErr)
		}
	}

	app.infoLog.Printf("New meeting have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete meeting by id
	deleteResult, err := app.meetings.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Have been eliminated %d meeting(s)", deleteResult.DeletedCount)
}
