package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rahul108/order_management_system/api/utils/customvalidator"
	"github.com/rahul108/order_management_system/api/utils/generator"
)

type Orders struct {
	ID                 uint64    `gorm:"primary_key;auto_increment" json:"id"`
	ConsignmentID      string    `gorm:"size:255;not null;unique" json:"consignment_id"`
	MerchantOrderid    string    `gorm:"size:255" json:"merchant_order_id"`
	StoreID            uint64    `gorm:"size:255;not null" json:"store_id" validate:"required"`
	RecipientName      string    `gorm:"size:255;not null" json:"recipient_name" validate:"required"`
	RecipientPhone     string    `gorm:"size:255;not null" json:"recipient_phone" validate:"required"`
	RecipientAddress   string    `gorm:"type:text;not null" json:"recipient_address" validate:"required"`
	ItemDescription    string    `gorm:"type:text;not null" json:"item_description"`
	RecipientCity      uint32    `gorm:"not null" json:"recipient_city" validate:"required"`
	RecipientZone      uint32    `gorm:"not null" json:"recipient_zone" validate:"required"`
	RecipientArea      uint32    `gorm:"not null" json:"recipient_area" validate:"required"`
	DeliveryType       uint32    `gorm:"not null" json:"delivery_type" validate:"required"`
	ItemType           uint32    `gorm:"not null" json:"item_type" validate:"required"`
	Status             uint32    `gorm:"not null" json:"status"`
	SpecialInstruction string    `gorm:"type:text" json:"special_instruction"`
	ItemQuantity       uint64    `gorm:"not null" json:"item_quantity" validate:"required"`
	ItemWeight         float64   `gorm:"not null" json:"item_weight" validate:"required"`
	AmountToCollect    float64   `gorm:"not null" json:"amount_to_collect" validate:"required"`
	DeliveryCost       int       `gorm:"not null" json:"delivery_cost"`
	CodFee             float64   `gorm:"not null" json:"cod_fee"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (o *Orders) Prepare() {
	o.ID = 0
	generated_string, _ := generator.GenerateRandomString(10)
	o.ConsignmentID = strconv.FormatUint(o.StoreID, 10) + strconv.FormatUint(o.ID, 10) + generated_string
	o.DeliveryCost = o.CalculateDeliveryCost()
	o.Status = 1
	o.CodFee = o.CalculateCodFee()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
}

func (o *Orders) CalculateDeliveryCost() int {
	calculatedDeliveryCost := 60
	totalWeight := float64(o.ItemQuantity) * float64(o.ItemWeight)
	if totalWeight > 0.5 && totalWeight <= 1.0 {
		calculatedDeliveryCost += 10
	} else {
		extra_weight := totalWeight - 1
		if (extra_weight - float64(int64(extra_weight))) != 0 {
			calculatedDeliveryCost += (15 * int(extra_weight+1))
		}
	}

	if o.RecipientCity != 1 {
		calculatedDeliveryCost += 40
	}
	return calculatedDeliveryCost
}

func (o *Orders) CalculateCodFee() float64 {
	percentage := 1
	return (o.AmountToCollect * float64(percentage)) / 100
}

func (o *Orders) Validate() error {
	err := customvalidator.ValidateBdPhoneNumber(o.RecipientPhone)
	if err != nil {
		return errors.New("please provide a valid recipient_phone")
	}

	return nil
}

func (o *Orders) SaveOrder(db *gorm.DB) (*Orders, error) {
	var err error = db.Debug().Model(&Orders{}).Create(&o).Error
	if err != nil {
		return &Orders{}, err
	}

	return o, nil
}

func (o *Orders) CancelOrder(db *gorm.DB, consignment_id string) error {
	err := db.Debug().Model(&Orders{}).Where("consignment_id = ?", consignment_id).Updates(map[string]interface{}{
		"status":     2,
		"updated_at": time.Now(),
	}).Error

	if err != nil {
		return err
	}

	return nil
}
