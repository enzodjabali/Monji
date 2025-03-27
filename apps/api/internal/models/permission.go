package models

// UserEnvPermission represents a user's permission on a specific environment.
type UserEnvPermission struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	EnvironmentID int    `json:"environment_id"`
	Permission    string `json:"permission"` // "none", "readOnly", "readAndWrite"
}

// UserDBPermission represents a user's permission on a specific database within an environment.
type UserDBPermission struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	EnvironmentID int    `json:"environment_id"`
	DBName        string `json:"db_name"`
	Permission    string `json:"permission"` // "none", "readOnly", "readAndWrite"
}
