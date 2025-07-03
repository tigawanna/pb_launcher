package utils_test

import (
	"pb_launcher/utils"
	"testing"
	"time"
)

func TestStrPointer(t *testing.T) {
	ptr := utils.StrPointer("test")
	if ptr == nil || *ptr != "test" {
		t.Errorf("expected 'test', got %v", ptr)
	}

	ptr = utils.StrPointer("")
	if ptr != nil {
		t.Errorf("expected nil, got %v", ptr)
	}
}
func TestPtr(t *testing.T) {
	intPtr := utils.Ptr(10)
	if intPtr == nil || *intPtr != 10 {
		t.Errorf("expected 10, got %v", intPtr)
	}

	floatPtr := utils.Ptr(3.14)
	if floatPtr == nil || *floatPtr != 3.14 {
		t.Errorf("expected 3.14, got %v", floatPtr)
	}

	strPtr := utils.Ptr("example")
	if strPtr == nil || *strPtr != "example" {
		t.Errorf("expected 'example', got %v", strPtr)
	}

	boolPtr := utils.Ptr(true)
	if boolPtr == nil || *boolPtr != true {
		t.Errorf("expected true, got %v", boolPtr)
	}

	nilInt := utils.Ptr(0)
	if nilInt != nil {
		t.Errorf("expected nil, got %v", nilInt)
	}

	nilFloat := utils.Ptr(0.0)
	if nilFloat != nil {
		t.Errorf("expected nil, got %v", nilFloat)
	}

	nilStr := utils.Ptr("")
	if nilStr != nil {
		t.Errorf("expected nil, got %v", nilStr)
	}

	nilBool := utils.Ptr(false)
	if nilBool != nil {
		t.Errorf("expected nil, got %v", nilBool)
	}

	now := time.Now()
	timePtr := utils.Ptr(now)
	if timePtr == nil || !timePtr.Equal(now) {
		t.Errorf("expected %v, got %v", now, timePtr)
	}
}
