package controllers

import (
	"net/http"

	"github.com/rahul108/order_management_system/api/responses"
)

func Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Order Management Backend")

}
