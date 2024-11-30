package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	helpers "github.com/rahul108/order_management_system/api/helper/orders"
	"github.com/rahul108/order_management_system/api/responses"
	services "github.com/rahul108/order_management_system/api/services/orders"
	utils "github.com/rahul108/order_management_system/api/utils/jwt"
)

func (server *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdOrderResponse := services.CreateOrderService(body, server.DB)

	if len(createdOrderResponse.Err) != 0 {
		responses.JSON(w, http.StatusUnprocessableEntity, helpers.CreateOrderCreationFailedResponse(createdOrderResponse.Err))
		return
	}

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, createdOrderResponse.Data.ID))

	responses.JSON(w, http.StatusCreated, helpers.CreateOrderSuccessResponse(*createdOrderResponse.Data))
}

func (server *Server) ListOrders(w http.ResponseWriter, r *http.Request) {

	// Parse query parameters
	queryParams := services.ExtractOrderQueryParams(r)

	response := services.GetOrdersList(queryParams, server.DB)

	utils.RespondWithJSON(w, http.StatusOK, helpers.OrderListGetResponse(response))
}

func (server *Server) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	consignmentID := vars["CONSIGNMENT_ID"]

	if consignmentID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing consignment ID")
		return
	}

	err := services.CancelOrder(consignmentID, server.DB)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, helpers.OrderCancelFailure(err))
	}
	utils.RespondWithJSON(w, http.StatusOK, helpers.OrderCancelSuccess())
}
