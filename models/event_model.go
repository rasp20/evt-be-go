package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Title       string             `json:"title,omitempty"`
	Start_Date  primitive.DateTime `json:"start_date,omitempty"`
	End_Date    primitive.DateTime `json:"end_date,omitempty"`
	Place       string             `json:"place,omitempty"`
	City        string             `json:"city,omitempty"`
	Province    string             `json:"province,omitempty"`
	Country     string             `json:"country,omitempty"`
	Image_Url   string             `json:"image_url,omitempty"`
	Description string             `json:"description,omitempty"`
	Url_Page    string             `json:"url_page,omitempty"`
	Is_Free     bool               `json:"is_free,omitempty"`
	Promo_Code  string             `json:"promo_code,omitempty"`
	Organizer   string             `json:"organizer,omitempty"`
	Is_Featured bool               `json:"is_featured,omitempty"`
}
