package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidatorError menyimpan pesan error field-specific
type ValidatorError map[string]string

// ValidateStruct memvalidasi struct berdasarkan tag `validate:"..."`.
// Mengembalikan map[field]pesan_error jika ada error, atau nil jika valid.
func ValidateStruct(s any) ValidatorError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(ValidatorError)
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := strings.ToLower(err.Field())
		errors[fieldName] = msgForTag(err)
	}
	return errors
}

// msgForTag mengubah error mesin menjadi pesan manusiawi
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("Minimal %s karakter", fe.Param())
	case "max":
		return fmt.Sprintf("Maksimal %s karakter", fe.Param())
	case "numeric":
		return "Harus berupa angka"
	case "alphanum":
		return "Hanya boleh huruf dan angka"
	default:
		return "Data tidak valid"
	}
}
