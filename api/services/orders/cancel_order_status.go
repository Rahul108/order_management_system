package services

import (
	"github.com/jinzhu/gorm"
	"github.com/rahul108/order_management_system/api/models"
)

func CancelOrder(consignment_id string, db *gorm.DB) error {
	order := models.Orders{}
	err := order.CancelOrder(db, consignment_id)
	if err != nil {
		return err
	}
	return nil
}
