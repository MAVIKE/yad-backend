package consts

const (
	CourierUnable  = 0
	CourierWaiting = 1
	CourierWorking = 2

	RestaurantUnable  = 1
	RestaurantWorking = 2

	OrderCreated           = 0
	OrderPaid              = 1
	OrderPreparing         = 2
	OrderWaitingForCourier = 3
	OrderEnRoute           = 4
	OrderDelivered         = 5
)
