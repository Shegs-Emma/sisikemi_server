package util

const (
	AVAILABLE = "available"
	OUT_OF_STOCK = "out_of_stock"
	DISCONTINUED = "discontinued"
	PENDING = "pending"
	SHIPPED = "shipped"
	DELIVERED = "delivered"
	CANCELLED = "cancelled"
)

func IsSupportedProductStatus(product_status string) bool {
	switch product_status {
	case AVAILABLE, OUT_OF_STOCK, DISCONTINUED:
		return true
	}
	return false
}

func IsSupportedOrderStatus(order_status string) bool {
	switch order_status {
	case PENDING, SHIPPED, DELIVERED, CANCELLED:
		return true
	}
	return false
}