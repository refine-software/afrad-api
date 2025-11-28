package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator

func RegisterTranslations(val *v.Validate) {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(val, trans)

	val.RegisterTagNameFunc(GetJSONTag)
}

func GetJSONTag(fld reflect.StructField) string {
	tag := fld.Tag.Get("json")
	if tag == "-" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}

func ParseValidationErrors(err error) map[string]string {
	out := make(map[string]string)

	if errs, ok := err.(v.ValidationErrors); ok {
		for _, e := range errs {
			out[e.Field()] = e.Translate(trans)
		}
	} else if err != nil {
		out["error"] = err.Error()
	}

	return out
}
