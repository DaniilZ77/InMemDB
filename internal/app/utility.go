package app

import (
	"errors"
	"strings"
)

func parseBytes(size string) (int, error) {
	sizeBytes := 0
	var unitMeasure string
	var i int
	for ; i < len(size) && '0' <= size[i] && size[i] <= '9'; i++ {
		sizeBytes = sizeBytes*10 + int(size[i]-'0')
	}
	unitMeasure = strings.TrimSpace(size[i:])

	switch strings.ToUpper(unitMeasure) {
	case "B":
		return sizeBytes, nil
	case "KB":
		return sizeBytes << 10, nil
	case "MB":
		return sizeBytes << 20, nil
	case "GB":
		return sizeBytes << 30, nil
	default:
		return 0, errors.New("invalid unit of measure")
	}
}
