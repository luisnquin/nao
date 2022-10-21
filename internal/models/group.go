package models

import "time"

type Group struct {
	// The name assigned by the user
	Name string `json:"name"`
	// Contains the keys of all the children
	Children []string `json:"children"`
	// Date when the group was created
	CreatedAt time.Time `json:"createdAt"`
}
