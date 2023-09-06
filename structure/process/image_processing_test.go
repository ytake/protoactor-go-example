package process

import (
	"testing"
	"time"
)

func TestImageProcessing_GetSpeed(t *testing.T) {
	i := &ImageProcessing{}
	image := "01012022 12:34:56.789|55|XYZ123"
	r, err := i.GetSpeed(image)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if r != 55 {
		t.Errorf("expected %v, got %v", 55, r)
	}
}

func TestImageProcessing_GetTime(t *testing.T) {
	i := &ImageProcessing{}
	image := "01012022 12:34:56.789|55|XYZ123"
	r, err := i.GetTime(image)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if r.String() != "2022-01-01 12:34:56.789 +0000 UTC" {
		t.Errorf("expected %v, got %v", "2022-01-01 12:34:56.789 +0000 UTC", r.String())
	}
}

func TestImageProcessing_GetLicense(t *testing.T) {
	i := &ImageProcessing{}
	image := "01012022 12:34:56.789|55|XYZ123"
	r, err := i.GetLicense(image)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if r != "XYZ123" {
		t.Errorf("expected %v, got %v", "XYZ123", r)
	}
}

func TestImageProcessing_CreatePhotoString(t *testing.T) {
	i := &ImageProcessing{}
	ti := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
	r := i.CreatePhotoString(ti, 60, "")
	if r != "03022001 04:05:06|60|" {
		t.Errorf("expected %v, got %v", "03022001 04:05:06|60|", r)
	}
	r = i.CreatePhotoString(ti, 60, "testing")
	if r != "03022001 04:05:06|60|testing" {
		t.Errorf("expected %v, got %v", "03022001 04:05:06|60|testing", r)
	}
}
