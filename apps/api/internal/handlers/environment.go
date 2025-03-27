package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"monji/internal/database"
	"monji/internal/models"
	"monji/internal/utils"

	"github.com/gin-gonic/gin"
)

// encryptionKey must be 16, 24, or 32 bytes long.
// For demonstration purposes, we use a hardcoded 32-byte key.
// In production, load this securely from an environment variable or secrets manager.
var encryptionKey = []byte("01234567890123456789012345678901")

// encrypt encrypts plainText using AES-GCM and returns a base64-encoded ciphertext.
func encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherBytes := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// decrypt decrypts a base64-encoded ciphertext using AES-GCM.
func decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plainBytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

// maskConnectionString masks the password in a connection string.
// For example, it converts:
// "mongodb://root:strongpassword@mongodb:27017/..."
// to "mongodb://root:s***************@mongodb:27017/..."
func maskConnectionString(conn string) string {
	schemeIndex := strings.Index(conn, "://")
	if schemeIndex == -1 {
		return conn
	}
	remainder := conn[schemeIndex+3:]
	atIndex := strings.Index(remainder, "@")
	if atIndex == -1 {
		return conn
	}
	credentials := remainder[:atIndex] // expected "username:password"
	colonIndex := strings.Index(credentials, ":")
	if colonIndex == -1 {
		return conn
	}
	username := credentials[:colonIndex]
	password := credentials[colonIndex+1:]
	if len(password) == 0 {
		return conn
	}
	// Only reveal the first character of the password
	maskedPassword := string(password[0])
	for i := 1; i < len(password); i++ {
		maskedPassword += "*"
	}
	maskedCreds := username + ":" + maskedPassword
	return conn[:schemeIndex+3] + maskedCreds + conn[schemeIndex+3+atIndex:]
}

// getEnvPermissionString returns the environment-level permission for the given user.
// If the user is admin/superadmin, it returns "readAndWrite".
func getEnvPermissionString(user models.User, envID int) string {
	if utils.IsAdmin(user) {
		return "readAndWrite"
	}
	row := database.DB.QueryRow(
		`SELECT permission FROM user_env_permissions WHERE user_id = ? AND environment_id = ?`,
		user.ID, envID,
	)
	var perm string
	if err := row.Scan(&perm); err != nil {
		return "none"
	}
	return perm
}

// CreateEnvironment encrypts the connection string before storing it.
func CreateEnvironment(c *gin.Context) {
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	if !utils.IsAdmin(currentUser) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin or superadmin can create environments"})
		return
	}
	var req struct {
		Name             string `json:"name"`
		ConnectionString string `json:"connection_string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" || req.ConnectionString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing environment name or connection string"})
		return
	}
	encryptedConn, err := encrypt(req.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt connection string: " + err.Error()})
		return
	}
	stmt, err := database.DB.Prepare(`INSERT INTO environments (name, connection_string, created_by) VALUES (?,?,?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	res, err := stmt.Exec(req.Name, encryptedConn, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last insert id"})
		return
	}
	env := models.Environment{
		ID:               int(id),
		Name:             req.Name,
		ConnectionString: encryptedConn,
		CreatedBy:        currentUser.ID,
	}
	// POST returns the stored (encrypted) value.
	c.JSON(http.StatusOK, gin.H{"environment": env})
}

// ListEnvironments returns all environments with the connection string decrypted and masked,
// and includes the current user's permission ("myPermission").
func ListEnvironments(c *gin.Context) {
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	if utils.IsAdmin(currentUser) {
		rows, err := database.DB.Query(`SELECT id, name, connection_string, created_by FROM environments`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var envs []gin.H
		for rows.Next() {
			var e models.Environment
			if err := rows.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy); err != nil {
				if err == sql.ErrNoRows {
					break
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			decryptedConn, err := decrypt(e.ConnectionString)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
				return
			}
			maskedConn := maskConnectionString(decryptedConn)
			envs = append(envs, gin.H{
				"id":                e.ID,
				"name":              e.Name,
				"connection_string": maskedConn,
				"created_by":        e.CreatedBy,
				"myPermission":      "readAndWrite",
			})
		}
		c.JSON(http.StatusOK, gin.H{"environments": envs})
		return
	}
	query := `
	SELECT e.id, e.name, e.connection_string, e.created_by, p.permission
	  FROM environments e
	  JOIN user_env_permissions p ON e.id = p.environment_id
	 WHERE p.user_id = ?
	   AND (p.permission = 'readOnly' OR p.permission = 'readAndWrite');
	`
	rows, err := database.DB.Query(query, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var envs []gin.H
	for rows.Next() {
		var e models.Environment
		var perm string
		if err := rows.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy, &perm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		decryptedConn, err := decrypt(e.ConnectionString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
			return
		}
		maskedConn := maskConnectionString(decryptedConn)
		envs = append(envs, gin.H{
			"id":                e.ID,
			"name":              e.Name,
			"connection_string": maskedConn,
			"created_by":        e.CreatedBy,
			"myPermission":      perm,
		})
	}
	c.JSON(http.StatusOK, gin.H{"environments": envs})
}

// GetEnvironment returns a single environment with the connection string decrypted and masked,
// along with the current user's permission.
func GetEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	var e models.Environment
	row := database.DB.QueryRow("SELECT id, name, connection_string, created_by FROM environments WHERE id = ?", id)
	if err := row.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}
	myPerm := getEnvPermissionString(currentUser, e.ID)
	if myPerm == "none" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to read this environment"})
		return
	}
	decryptedConn, err := decrypt(e.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}
	maskedConn := maskConnectionString(decryptedConn)
	c.JSON(http.StatusOK, gin.H{
		"environment": gin.H{
			"id":                e.ID,
			"name":              e.Name,
			"connection_string": maskedConn,
			"created_by":        e.CreatedBy,
		},
		"myPermission": myPerm,
	})
}

// UpdateEnvironment updates an environment configuration.
// If a new connection string is provided, it is encrypted before storage.
func UpdateEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	hasWrite, err := utils.HasEnvPermission(currentUser, id, "write")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to update this environment"})
		return
	}
	var req struct {
		Name             string `json:"name"`
		ConnectionString string `json:"connection_string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" && req.ConnectionString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update parameters provided"})
		return
	}
	query := "UPDATE environments SET "
	var params []interface{}
	var updates []string
	if req.Name != "" {
		updates = append(updates, "name = ?")
		params = append(params, req.Name)
	}
	if req.ConnectionString != "" {
		encryptedConn, err := encrypt(req.ConnectionString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt connection string: " + err.Error()})
			return
		}
		updates = append(updates, "connection_string = ?")
		params = append(params, encryptedConn)
	}
	query += joinUpdates(updates, ", ") + " WHERE id = ?"
	params = append(params, id)
	res, err := database.DB.Exec(query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}
	var e models.Environment
	row := database.DB.QueryRow("SELECT id, name, connection_string, created_by FROM environments WHERE id = ?", id)
	if err := row.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	decryptedConn, err := decrypt(e.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}
	maskedConn := maskConnectionString(decryptedConn)
	c.JSON(http.StatusOK, gin.H{"environment": gin.H{
		"id":                e.ID,
		"name":              e.Name,
		"connection_string": maskedConn,
		"created_by":        e.CreatedBy,
	}})
}

// DeleteEnvironment removes an environment configuration.
func DeleteEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	hasWrite, err := utils.HasEnvPermission(currentUser, id, "write")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to delete this environment"})
		return
	}
	res, err := database.DB.Exec("DELETE FROM environments WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Environment deleted successfully"})
}
