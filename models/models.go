package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title,omitempty"`
	Categoty       string             `json:"category,omitempty"`
	StartDate      primitive.DateTime `json:"startDate,omitempty"`
	End_Date       primitive.DateTime `json:"endDate,omitempty"`
	Place          string             `json:"place,omitempty"`
	City           string             `json:"city,omitempty"`
	Province       string             `json:"province,omitempty"`
	Country        string             `json:"country,omitempty"`
	Image          string             `json:"image,omitempty"`
	Description    string             `json:"description,omitempty"`
	Url            string             `json:"url,omitempty"`
	IsFree         bool               `json:"isFree,omitempty"`
	TicketPrice    int                `json:"ticketPrice,omitempty"`
	TicketCategory string             `json:"ticketCategory,omitempty"`
	PromoCode      string             `json:"promoCode,omitempty"`
}
