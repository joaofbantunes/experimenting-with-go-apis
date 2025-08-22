package shared

const globalPrefix = "tag:example.com,2025:problems:sample-api/"

const (
	generalPrefix            = globalPrefix + "general/"
	ProblemGeneralNotFound   = generalPrefix + "not-found"
	ProblemGeneralValidation = generalPrefix + "validation"
)

const (
	ordersPrefix                     = globalPrefix + "orders/"
	ProblemOrdersUnknownDishes       = ordersPrefix + "unknown-dishes"
	ProblemOrdersNoLongerCancellable = ordersPrefix + "order-no-longer-cancellable"
)
