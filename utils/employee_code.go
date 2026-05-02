package utils

import (
	"fmt"
	"time"
)

func GenerateEmployeeCode() string {
	return fmt.Sprintf("KHM-%x", time.Now().UnixNano()>>20)
}
