package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func SetupValidation() {
	validate = validator.New()
}

func ValidateInput(c *gin.Context, input interface{}) error {
	if err := c.ShouldBindJSON(input); err != nil {
		return err
	}

	if err := validate.Struct(input); err != nil {
		return err
	}

	return nil
}
