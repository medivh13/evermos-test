package models

type Product struct {
	ID   int64  `db:"p_id"`
	Name string `db:"p_name"`
	Desc string `db:"p_desc"`
	QTY  int64  `db:"p_stocks"`
}
