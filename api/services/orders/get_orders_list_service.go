package services

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/rahul108/order_management_system/api/models"
)

type OrderQueryParams struct {
	TransferStatus int    `json:"transfer_status"`
	Archive        bool   `json:"archive"`
	Limit          int    `json:"limit"`
	Page           int    `json:"page"`
	Sort           string `json:"sort"`
}

type PaginatedOrderResponse struct {
	Total      int64           `json:"total"`
	Page       int             `json:"current_page"`
	Limit      int             `json:"per_page"`
	TotalPages int             `json:"total_in_pages"`
	Data       []models.Orders `json:"data"`
}

func ExtractOrderQueryParams(r *http.Request) OrderQueryParams {
	query := r.URL.Query()

	// Parse transfer_status
	transferStatus, _ := strconv.Atoi(query.Get("transfer_status"))

	// Parse archive (convert to bool)
	archive, _ := strconv.ParseBool(query.Get("archive"))

	// Parse limit
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10 // default limit
	}

	// Parse page
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1 // default page
	}

	// Get sort parameter
	sort := query.Get("sort")

	return OrderQueryParams{
		TransferStatus: transferStatus,
		Archive:        archive,
		Limit:          limit,
		Page:           page,
		Sort:           sort,
	}
}

func GetOrdersList(params OrderQueryParams, db *gorm.DB) PaginatedOrderResponse {
	query := db.Model(&models.Orders{})

	// Apply filters
	if params.TransferStatus > 1 {
		query = query.Where("transfer_status = ?", params.TransferStatus)
	}

	// Set default limit and page
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Page == 0 {
		params.Page = 1
	}

	// Calculate offset
	offset := (params.Page - 1) * params.Limit

	// Count total records (for pagination)
	var totalCount int64
	query.Count(&totalCount)

	// Calculate total pages
	totalPages := int(totalCount) / params.Limit
	if int(totalCount)%params.Limit != 0 {
		totalPages++
	}

	// Apply sorting
	if params.Sort != "" {
		query = query.Order(params.Sort)
	} else {
		// Default sorting
		query = query.Order("created_at DESC")
	}

	// Fetch orders with pagination
	var orders []models.Orders
	query.Limit(params.Limit).Offset(offset).Find(&orders)

	// Prepare paginated response
	response := PaginatedOrderResponse{
		Total:      totalCount,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
		Data:       orders,
	}

	return response
}
