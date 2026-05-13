package utils

import "fmt"

func GenerateLoanCode(id int) string {
	return fmt.Sprintf("001LN%06d", id)
}
