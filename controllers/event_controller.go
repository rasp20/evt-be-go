package controllers

import (
	"evt-be-go/configs"
	"evt-be-go/models"
	"evt-be-go/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var eventCollection *mongo.Collection = configs.GetCollection(configs.DB, "event-data")
var validate = validator.New()

func CreateEvent(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var event models.Event
	defer cancel()

	//validate the request body
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&event); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newEvent := models.Event{
		Id:          primitive.NewObjectID(),
		Title:       event.Title,
		Start_Date:  event.Start_Date,
		End_Date:    event.End_Date,
		Place:       event.Place,
		City:        event.City,
		Province:    event.Province,
		Country:     event.Country,
		Image_Url:   event.Image_Url,
		Description: event.Description,
		Url_Page:    event.Url_Page,
		Is_Free:     event.Is_Free,
		Promo_Code:  event.Promo_Code,
		Organizer:   event.Organizer,
		Is_Featured: event.Is_Featured,
	}

	result, err := eventCollection.InsertOne(ctx, newEvent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.EventResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

func GetAnEvent(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	eventId := c.Param("eventId")
	var event models.Event
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	err := eventCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&event)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": event}})
}

func EditAnEvent(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	eventId := c.Param("eventId")
	var event models.Event
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	//validate the request body
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&event); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"Title": event.Title, "Start_Date": event.Start_Date, "End_Date": event.End_Date, "Place": event.Place, "City": event.City, "Province": event.Province, "Country": event.Country, "Image_Url": event.Image_Url, "Description": event.Description, "Url_Page": event.Url_Page, "Is_Free": event.Is_Free, "Promo_Code": event.Promo_Code, "Organizer": event.Organizer, "Is_Featured": event.Is_Featured}

	result, err := eventCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//get updated event details
	var updatedEvent models.Event
	if result.MatchedCount == 1 {
		err := eventCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedEvent)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
	}

	return c.JSON(http.StatusOK, responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedEvent}})
}

func DeleteAnEvent(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	eventId := c.Param("eventId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	result, err := eventCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.EventResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "Event with specified ID not found!"}})
	}

	return c.JSON(http.StatusOK, responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "Event successfully deleted!"}})
}

func GetAllEvent(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var events []models.Event
	defer cancel()

	results, err := eventCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleEvent models.Event
		if err = results.Decode(&singleEvent); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		events = append(events, singleEvent)
	}

	return c.JSON(http.StatusOK, responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": events}})
}
