package utils

import "fmt"

func GenerateLoanCode(id int) string {
	return fmt.Sprintf("002LN%06d", id)
}
