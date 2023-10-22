package controllers

import (
	"evt-be-go/configs"
	"evt-be-go/models"
	"evt-be-go/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	// "github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	return c.JSON(http.StatusCreated, responses.EventResponse{Status: http.StatusCreated, Message: "Event created successfully", Data: &echo.Map{"data": result}})
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
	defer cancel()

	// Get search parameters from query parameters
	eventName := c.QueryParam("name")
	place := c.QueryParam("place")
	city := c.QueryParam("city")

	// Define the filter based on the search parameters
	filter := bson.M{}

	// Create an array to hold the OR conditions
	orConditions := []bson.M{}

	if eventName != "" {
		// Add Title filter condition to the OR array
		orConditions = append(orConditions, bson.M{"Title": bson.M{"$regex": eventName, "$options": "i"}})
	}
	if place != "" {
		// Add Place filter condition to the OR array
		orConditions = append(orConditions, bson.M{"Place": bson.M{"$regex": place, "$options": "i"}})
	}
	if city != "" {
		// Add City filter condition to the OR array
		orConditions = append(orConditions, bson.M{"City": bson.M{"$regex": city, "$options": "i"}})
	}

	// Use the $or operator to combine the OR conditions
	if len(orConditions) > 0 {
		filter["$or"] = orConditions
	}

	// Define a sort filter to order the events by ascending Is_Featured and Start_Date
	sortFilter := bson.D{
		{Key: "Is_Featured", Value: 1}, // 1 for ascending order
		{Key: "Start_Date", Value: 1},  // 1 for ascending order
	}

	// Create a MongoDB cursor with the filters and sorting criteria
	findOptions := options.Find().SetSort(sortFilter)
	results, err := eventCollection.Find(ctx, filter, findOptions)

	if err != nil {
		// Handle the error
		return c.JSON(http.StatusInternalServerError, responses.EventResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": nil},
		})
	}

	defer results.Close(ctx)

	// Create a slice to store the retrieved event data
	data := &echo.Map{"data": []models.Event{}}

	for results.Next(ctx) {
		var singleEvent models.Event
		if err = results.Decode(&singleEvent); err != nil {
			// Handle the error
			return c.JSON(http.StatusInternalServerError, responses.EventResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": nil},
			})
		}
		// Cast and dereference data for indexing
		eventData := (*data)["data"].([]models.Event)

		// Append the event to the eventData slice
		eventData = append(eventData, singleEvent)

		// Update the data map with the modified eventData
		(*data)["data"] = eventData
	}

	// Check if any events were found
	if len((*data)["data"].([]models.Event)) == 0 {
		return c.JSON(http.StatusNotFound, responses.EventResponse{
			Status:  http.StatusNotFound,
			Message: "no events found",
			Data:    data,
		})
	}

	return c.JSON(http.StatusOK, responses.EventResponse{
		Status:  http.StatusOK,
		Message: "Events retrieved successfully",
		Data:    data,
	})
}
