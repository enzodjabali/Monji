package utils

import (
	"database/sql"
	"errors"
	"fmt"

	"monji/internal/database"
	"monji/internal/models"
)

// IsAdmin returns true if user's role is "admin" or "superadmin".
func IsAdmin(user models.User) bool {
	return user.Role == "admin" || user.Role == "superadmin"
}

// IsSuperAdmin returns true if user's role is "superadmin".
func IsSuperAdmin(user models.User) bool {
	return user.Role == "superadmin"
}

// HasEnvPermission checks if the given user has the required environment permission.
//
// required can be "read" or "write".
// - If user is admin/superadmin, return true immediately.
// - Otherwise look at user_env_permissions table.
//
//	If required == "read", we accept both "readOnly" or "readAndWrite" stored in DB.
//	If required == "write", we accept only "readAndWrite" stored in DB.
func HasEnvPermission(user models.User, envID int, required string) (bool, error) {
	// admin or superadmin => automatically pass
	if IsAdmin(user) {
		return true, nil
	}

	// normal user => check user_env_permissions
	row := database.DB.QueryRow(
		`SELECT permission FROM user_env_permissions WHERE user_id = ? AND environment_id = ?`,
		user.ID, envID,
	)

	var perm string
	err := row.Scan(&perm)
	if err != nil {
		if err == sql.ErrNoRows {
			// no permission row => no access
			return false, nil
		}
		return false, err
	}

	switch required {
	case "read":
		// "readOnly" or "readAndWrite" are acceptable
		if perm == "readOnly" || perm == "readAndWrite" {
			return true, nil
		}
		return false, nil
	case "write":
		// must be "readAndWrite" to proceed
		if perm == "readAndWrite" {
			return true, nil
		}
		return false, nil
	default:
		return false, errors.New("invalid required permission type")
	}
}

// HasDBPermission checks if the user has the required database permission.
//
// required can be "read" or "write".
//   - If user is admin/superadmin, return true immediately.
//   - Otherwise, user must have at least read permission on the environment AND
//     must have at least the required permission in user_db_permissions for that db.
func HasDBPermission(user models.User, envID int, dbName string, required string) (bool, error) {
	if IsAdmin(user) {
		return true, nil
	}

	// must have environment read permission as a baseline
	hasEnvRead, err := HasEnvPermission(user, envID, "read")
	if err != nil {
		return false, err
	}
	if !hasEnvRead {
		return false, nil
	}

	// Now check user_db_permissions
	row := database.DB.QueryRow(
		`SELECT permission FROM user_db_permissions WHERE user_id = ? AND environment_id = ? AND db_name = ?`,
		user.ID, envID, dbName,
	)

	var perm string
	err = row.Scan(&perm)
	if err != nil {
		if err == sql.ErrNoRows {
			// no permission => no access
			return false, nil
		}
		return false, err
	}

	switch required {
	case "read":
		if perm == "readOnly" || perm == "readAndWrite" {
			return true, nil
		}
		return false, nil
	case "write":
		if perm == "readAndWrite" {
			return true, nil
		}
		return false, nil
	default:
		return false, fmt.Errorf("invalid required DB permission type: %s", required)
	}
}
