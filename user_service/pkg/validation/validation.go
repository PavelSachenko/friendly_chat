package validation

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type Validator struct {
	*validator.Validate
}

func InitValidator() *Validator {
	return &Validator{
		validator.New(),
	}
}

func (v *Validator) Struct(s interface{}, tagName string) error {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tagName), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return v.StructCtx(context.Background(), s)
}

func (v *Validator) in_array(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
func (v *Validator) ValidateRequest(ctx *gin.Context, obj any) []*IError {
	var tagName string
	if !v.in_array([]string{"GET", "DELETE"}, ctx.Request.Method) {
		tagName = "json"
		err := ctx.ShouldBindJSON(&obj)
		if err != nil {
			err := v.Struct(obj, tagName)
			if err != nil {
				return v.buildJsonErrors(err)
			}
			return []*IError{{Tag: "unknown error", Value: err.Error()}}
		}
	} else {
		tagName = "form"
		err := ctx.ShouldBindQuery(&obj)
		if err != nil {
			return []*IError{{Tag: "unknown error", Value: err.Error()}}
		}
	}

	err := v.Struct(obj, tagName)
	if err != nil {
		return v.buildJsonErrors(err)
	}

	return nil
}

type IError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (v *Validator) buildJsonErrors(err error) []*IError {

	var Errors []*IError
	for _, err := range err.(validator.ValidationErrors) {
		var el IError
		el.Field = strings.ToLower(err.Field())
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		Errors = append(Errors, &el)
	}
	return Errors
}
