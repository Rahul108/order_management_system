package controllers

import (
	"net/http"

	"github.com/rahul108/order_management_system/api/models"
	utils "github.com/rahul108/order_management_system/api/utils/jwt"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := user.SetPassword(req.Password); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	result := s.DB.Create(&user)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusConflict, "User already exists")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user models.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.Username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.Username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	user.RefreshToken = refreshToken
	s.DB.Save(&user)

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := utils.ParseJSON(r, &refreshRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user models.User
	if err := s.DB.Where("refresh_token = ?", refreshRequest.RefreshToken).First(&user).Error; err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.Username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate new access token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}

func ProtectedProfile(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Protected route",
		"username": username,
	})
}
