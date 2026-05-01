package utils

import (
	"fmt"
	"time"
)

func GenerateEmployeeCode() string {
	return fmt.Sprintf("KHM-%d", time.Now().UnixNano())
}
