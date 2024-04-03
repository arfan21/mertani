package entity

import "time"

type Device struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Sensors     []Sensor  `json:"sensors"`
}

func (Device) TableName() string {
	return "devices"
}
