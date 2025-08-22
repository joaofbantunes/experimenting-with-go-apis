package shared

const globalPrefix = "tag:example.com,2025:problems:sample-api/"

const (
	generalPrefix = globalPrefix + "general/"
	NotFound      = generalPrefix + "not-found"
	Validation    = generalPrefix + "validation"
)

const (
	ordersPrefix             = globalPrefix + "orders/"
	UnknownDishes            = ordersPrefix + "unknown-dishes"
	OrderNoLongerCancellable = ordersPrefix + "order-no-longer-cancellable"
)
