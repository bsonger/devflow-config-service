package app

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func placeholderClause(column string, index int) string {
	return fmt.Sprintf("%s=$%d", column, index)
}

func marshalJSON(value any) []byte {
	if value == nil {
		return []byte("[]")
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return []byte("[]")
	}
	return payload
}

func nullableUUIDPtr(id *uuid.UUID) any {
	if id == nil || *id == uuid.Nil {
		return nil
	}
	return *id
}
