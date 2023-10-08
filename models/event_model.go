package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Id          primitive.ObjectID `json:"_id"`
	Title       string             `json:"Title" validate:"required"`
	Start_Date  primitive.DateTime `json:"Start_Date"`
	End_Date    primitive.DateTime `json:"End_Date"`
	Place       string             `json:"Place"`
	City        string             `json:"City"`
	Province    string             `json:"Province"`
	Country     string             `json:"Country"`
	Image_Url   string             `json:"Image_Url"`
	Description string             `json:"Description"`
	Url_Page    string             `json:"Url_Page"`
	Is_Free     bool               `json:"Is_Free"`
	Promo_Code  string             `json:"Promo_Code"`
	Organizer   string             `json:"Organizer"`
	Is_Featured bool               `json:"Is_Featured"`
}
