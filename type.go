package backren

import "time"

type User struct {
	Name           string `bson:"name,omitempty" json:"name,omitempty"`
	Username       string `bson:"username,omitempty" json:"username,omitempty"`
	Email          string `bson:"email,omitempty" json:"email,omitempty"`
	PhoneNumber    string `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	License        string `bson:"license,omitempty" json:"license,omitempty"`
	Password       string `bson:"password,omitempty" json:"password,omitempty"`
}

type Response struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Username string    `json:"username"`
	Exp      time.Time `json:"exp"`
	Iat      time.Time `json:"iat"`
	Nbf      time.Time `json:"nbf"`
}
