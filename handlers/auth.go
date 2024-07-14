// handlers/auth.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
	"myapp/models"
	"myapp/utils"
	"myapp/middleware"
)

var jwtKey = []byte("JWT_SECRET")

type Credentials struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func RegisterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start register")
		var newUser models.User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Hash the password before storing it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		newUser.Password = string(hashedPassword)

		// Create the user in the database
		if err := db.Create(&newUser).Error; err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Respond with success message or user object if needed
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("username = ? OR email = ?", creds.UsernameOrEmail, creds.UsernameOrEmail).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Create JWT tokens
		expirationTime := time.Now().Add(15 * time.Minute)
		tokenString, err := utils.GenerateJWT(user.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		refreshTokenString, err := utils.GenerateRefreshToken(user.Username)
		if err != nil {
			http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
			return
		}

		// Set JWT token as cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token":         tokenString,
			"refresh_token": refreshTokenString,
		})
	}
}

func ChangePasswordHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var changePasswordRequest ChangePasswordRequest
		err := json.NewDecoder(r.Body).Decode(&changePasswordRequest)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Get current user
		username := middleware.GetUsername(r)
		var user models.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			log.Fatalf("failed to retrieve user: %v", err)
		}

		// Verify old password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordRequest.OldPassword)); err != nil {
			http.Error(w, "Invalid old password", http.StatusUnauthorized)
			return
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
			return
		}

		// Update the password in the database
		user.Password = string(hashedPassword)
		if err := db.Save(&user).Error; err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}