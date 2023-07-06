package mongo

type CostParams struct {
	LowCost  int
	HighCost int
}

type AddToCartParams struct {
	UserID     string
	PurchaseID string
}
