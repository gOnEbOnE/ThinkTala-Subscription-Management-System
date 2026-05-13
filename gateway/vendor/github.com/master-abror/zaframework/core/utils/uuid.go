package utils

import (
	"github.com/google/uuid"
)

func CreateUUID() (uuid.UUID, error) {
	val, err := uuid.NewV7()
	if err != nil {
		return val, err
	}
	return val, nil
}
