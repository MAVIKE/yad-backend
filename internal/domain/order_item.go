package domain

type OrderItem struct {
	Id         int `json:"id" db:"id"`
	OrderId    int `json:"order_id" db:"order_id"`
	MenuItemId int `json:"menu_item_id" db:"menu_item_id"`
	Count      int `json:"count" db:"count"`
}
