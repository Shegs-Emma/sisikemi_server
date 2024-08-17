package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Intn(100)
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomUser() string {
	return RandomString(6)
}

func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

func RandomProductStatus() string {
	product_statuses := []string{AVAILABLE, OUT_OF_STOCK, DISCONTINUED}
	n := len(product_statuses)
	return product_statuses[rand.Intn(n)]
}

func RandomOrderStatus() string {
	order_statuses := []string{PENDING, SHIPPED, DELIVERED, CANCELLED}
	n := len(order_statuses)
	return order_statuses[rand.Intn(n)]
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}