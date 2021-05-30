package models

type CheckoutOrder struct {
	ID     int64 `db:"o_id"`
	Status int64 `db:"o_status"`
}

type OrderToProcess struct {
	ID        int64 `db:"o_id"`
	ProductID int64 `db:"od_p_id"`
	Qty       int64 `db:"od_qty"`
}
