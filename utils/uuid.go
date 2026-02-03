package utils

import "github.com/google/uuid"

func UuidToBytes(id uuid.UUID) ([]byte, error) {
	return id.MarshalBinary()
}

func BytesToUUID(b []byte) (uuid.UUID, error) {
	return uuid.FromBytes(b)
}
