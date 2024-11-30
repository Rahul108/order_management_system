package helpers

import (
	"github.com/rahul108/order_management_system/api/models"
	services "github.com/rahul108/order_management_system/api/services/orders"
)

func CreateOrderSuccessResponse(data models.Orders) map[string]interface{} {
	return map[string]interface{}{
		"message": "Order Created Successfully",
		"type":    "success",
		"code":    201,
		"data":    postProcessDataToReturn(data),
	}
}

func CreateOrderCreationFailedResponse(err map[string][]string) map[string]interface{} {
	return map[string]interface{}{
		"message": "Please fix the given errors",
		"type":    "error",
		"code":    422,
		"errors":  err,
	}
}

func postProcessDataToReturn(data models.Orders) map[string]interface{} {
	response := make(map[string]interface{})
	response["consignment_id"] = data.ConsignmentID
	response["merchant_order_id"] = data.MerchantOrderid
	response["order_status"] = "PENDING"
	response["delivery_fee"] = data.DeliveryCost

	return response
}

func OrderListGetResponse(data services.PaginatedOrderResponse) map[string]interface{} {
	return map[string]interface{}{
		"message": "Orders successfully fetched.",
		"type":    "success",
		"code":    200,
		"data":    data,
	}
}

func OrderCancelSuccess() map[string]interface{} {
	return map[string]interface{}{
		"message": "Order Cancelled Successfully",
		"type":    "success",
		"code":    200,
	}
}

func OrderCancelFailure(err error) map[string]interface{} {
	return map[string]interface{}{
		"message": "Order cancel failed",
		"type":    "failed",
		"code":    200,
		"error":   err.Error(),
	}
}
