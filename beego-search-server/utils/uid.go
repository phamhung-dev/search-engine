package utils

import (
	"encoding/base64"
	"strconv"
)

func EncodeUID(id uint64, tableID int) string {
	combinedID := strconv.FormatUint(id<<32|uint64(tableID), 10)
	fakeID := base64.RawStdEncoding.EncodeToString([]byte(combinedID))
	return fakeID
}

func DecodeUID(encoded string) (uint64, int, error) {
	idBytes, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return 0, 0, err
	}
	combinedID, err := strconv.ParseUint(string(idBytes), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	id := combinedID >> 32
	tableID := int(combinedID & 0xFFFFFFFF)

	return id, tableID, nil
}
