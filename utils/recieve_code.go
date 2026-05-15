package utils

import (
	"fmt"
	"time"
)

func GenerateRecieveCode() string {
	return fmt.Sprintf("RE-%x", time.Now().UnixNano()>>20)
}
