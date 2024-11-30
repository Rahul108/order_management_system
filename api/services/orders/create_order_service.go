package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/rahul108/order_management_system/api/models"
	"github.com/rahul108/order_management_system/api/utils/formaterror"
)

var validate = validator.New()
var hardcodedData = map[string]interface{}{
	"store_id":          131172,
	"recipient_city":    1,
	"recipient_zone":    1,
	"recipient_area":    1,
	"delivery_type":     48,
	"item_type":         2,
	"item_quantity":     1,
	"item_weight":       0.5,
	"recipient_address": "banani, gulshan 2, dhaka, bangladesh",
}

type ResponseFromCreateOrder struct {
	Data *models.Orders
	Err  map[string][]string
}

func HardCodedValidation(validationErrors map[string][]string, data models.Orders) map[string][]string {
	reflectValue := reflect.ValueOf(data)
	reflectType := reflectValue.Type()
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		fieldValue := field.Interface()
		fieldTag := reflectType.Field(i).Tag.Get("json")
		if hardcodedData[fieldTag] != nil && fmt.Sprint(hardcodedData[fieldTag]) != fmt.Sprint(fieldValue) {
			validationErrors[fieldTag] = append(validationErrors[fieldTag], fmt.Sprintf("Wrong %s field Selected", fieldTag))
		}
	}
	return validationErrors
}

func ToSnakeCase(input string) string {
	regex := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := regex.ReplaceAllString(input, "${1}_${2}")
	return strings.ToLower(snake)
}

func FieldValidatorForOrders(data models.Orders) map[string][]string {
	validationErrors := make(map[string][]string)
	// will be removed once the hardcoded validation is not required
	validationErrors = HardCodedValidation(validationErrors, data)
	err := validate.Struct(data)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := ToSnakeCase(err.Field())
			validationErrors[field] = append(validationErrors[field], fmt.Sprintf("The %s field is required.", field))
		}
	}

	return validationErrors
}

func CreateOrderService(data []byte, db *gorm.DB) ResponseFromCreateOrder {
	validationErrors := make(map[string][]string)
	order := models.Orders{}
	err := json.Unmarshal(data, &order)
	if err != nil {
		validationErrors["validationError"] = append(validationErrors["validationError"], err.Error())
	}
	order.Prepare()
	validationErrors = FieldValidatorForOrders(order)
	err = order.Validate()
	if err != nil {
		validationErrors["validationError"] = append(validationErrors["validationError"], err.Error())
	}

	if len(validationErrors) != 0 {
		return ResponseFromCreateOrder{
			Data: nil,
			Err:  validationErrors,
		}
	}

	orderCreated, err := order.SaveOrder(db)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		validationErrors["formatError"] = append(validationErrors["formatError"], formattedError.Error())
		return ResponseFromCreateOrder{
			Data: nil,
			Err:  validationErrors,
		}
	}

	return ResponseFromCreateOrder{
		Data: orderCreated,
		Err:  validationErrors,
	}
}
