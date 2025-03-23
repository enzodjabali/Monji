package models

// User represents an application user.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Company   string `json:"company,omitempty"`
	Password  string `json:"password,omitempty"`
	Role      string `json:"role"`
}
