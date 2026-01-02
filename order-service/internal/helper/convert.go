package helper

import (
	"errors"
	"strconv"

	"github.com/gofrs/uuid"
)

func StringToUUID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid format")
	}

	return parsedID, nil
}

func StringToIntDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil || i <= 0 {
		return defaultValue
	}

	return i
}