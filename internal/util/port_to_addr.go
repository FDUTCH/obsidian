package util

import "fmt"

func Addr(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}
