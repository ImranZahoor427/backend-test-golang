package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func ParseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID '%s': %w", s, err)
	}
	return id, nil
}
