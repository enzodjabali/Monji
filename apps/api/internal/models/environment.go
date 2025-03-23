package models

// Environment represents a MongoDB environment configuration.
type Environment struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ConnectionString string `json:"connection_string"`
	CreatedBy        int    `json:"created_by"`
}
