package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"gorm.io/gorm"
	"reflect"
	"time"
)

func GetAllInvestors(_ *gin.Context) (*[]models.Investor, error) {
	var investors []models.Investor
	if result := initializers.DB.Find(&investors); result.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &investors, nil
}

func GetInvestor(_ *gin.Context, in *routes.GetInvestorRequest) (*models.Investor, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var investor models.Investor
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&investor); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Investor not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	return &investor, nil
}

func CreateNewInvestor(c *gin.Context, in *routes.InvestorCreationRequest) (*models.Investor, error) {
	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return nil, errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return nil, errors.New("Internal server error")
	}

	if authUser.Role != "admin" {
		return nil, errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

	var investorFound models.Investor
	if rst := initializers.DB.Where("email=?", in.Email).Find(&investorFound); rst.Error == nil {
		return nil, errors.NewAlreadyExists(nil, "Email already used")
	} else {
		if !errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Internal server error")
		}
	}

	currentDate := time.Now().Format("2006-01-02")
	investor := models.Investor{
		UUID:            uuid.New().String(),
		ID:              nil,
		Name:            in.Name,
		LegalStatus:     nil,
		Address:         nil,
		Email:           in.Email,
		Phone:           nil,
		CreatedAt:       &currentDate,
		Description:     nil,
		InvestorType:    nil,
		InvestmentFocus: nil,
	}

	if err := initializers.DB.Create(&investor); err.Error != nil {
		return nil, errors.New("Internal server error")
	}

	return &investor, nil
}

func DeleteInvestor(c *gin.Context, in *routes.DeleteInvestorRequest) error {
	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return errors.New("Internal server error")
	}

	if authUser.Role != "admin" && authUser.UUID != in.UUID {
		return errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var investor models.Investor
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&investor); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewUserNotFound(nil, "Investor not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if rst := initializers.DB.Delete(&investor); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func UpdateInvestor(_ *gin.Context, in *routes.InvestorUpdateRequest) (*models.Investor, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var investor models.Investor
	if rst := initializers.DB.Where("uuid=?", in.UUID).First(&investor); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Investor not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	updates := make(map[string]interface{})

	val := reflect.ValueOf(*in)
	typ := reflect.TypeOf(*in)
	hasUpdate := false

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" || jsonTag == "uuid" {
			continue
		}
		fieldValue := val.Field(i)
		if fieldValue.Kind() == reflect.String && fieldValue.String() != "" {
			hasUpdate = true
			updates[jsonTag] = fieldValue.String()
		}
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			strVal, ok := fieldValue.Elem().Interface().(string)
			if ok && strVal != "" {
				hasUpdate = true
				updates[jsonTag] = strVal
			}
		}
	}

	if !hasUpdate {
		return nil, errors.NewNotValid(nil, "Invalid body")
	}

	if err := initializers.DB.Model(&investor).Updates(updates).Error; err != nil {
		return nil, errors.New("Internal server error")
	}

	return &investor, nil
}
