package cart

import (
	"fmt"

	"github.com/robsonpaixao/ecom-api/types"
)

func getCartItemsIDs(items []types.CartCheckoutItem) ([]int, error) {
	productsIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productsIds[i] = item.ProductID
	}

	return productsIds, nil
}

func (h *Handler) createOder(ps []types.Product, items []types.CartCheckoutItem, userId int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(items, productMap)

	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.store.UpdateProduct(product)
	}

	orderId, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderId,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderId, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartCheckoutItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartCheckoutItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}
