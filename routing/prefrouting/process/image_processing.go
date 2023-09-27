package process

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ImageProcessing struct{}

func splitImageAttributes(image string) []string {
	return strings.Split(image, "|")
}

// GetSpeed is a ImageProcessing
func (ip *ImageProcessing) GetSpeed(image string) (int, error) {
	attributes := splitImageAttributes(image)
	if len(attributes) == 3 {
		speed, err := strconv.Atoi(attributes[1])
		if err != nil {
			return 0, err
		}
		return speed, nil
	}
	return 0, fmt.Errorf("invalid image format")
}

// GetTime is a ImageProcessing
func (ip *ImageProcessing) GetTime(image string) (time.Time, error) {
	attributes := splitImageAttributes(image)
	if len(attributes) == 3 {
		t, err := time.Parse("02012006 15:04:05.999", attributes[0])
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid image format")
}

// GetLicense is a ImageProcessing
func (ip *ImageProcessing) GetLicense(image string) (string, error) {
	attributes := splitImageAttributes(image)
	if len(attributes) == 3 {
		return attributes[2], nil
	}
	return "", fmt.Errorf("invalid image format")
}

// CreatePhotoString is a ImageProcessing
func (ip *ImageProcessing) CreatePhotoString(date time.Time, speed int, license string) string {
	// SimpleDateFormat("ddMMyyyy HH:mm:ss.SSS") に対応したもの
	return fmt.Sprintf("%s|%d|%s", date.Format("02012006 15:04:05.999"), speed, license)
}
