package models

type Costumer struct {
	ID       int64  `db:"c_id"`
	Email    string `db:"c_email"`
	Password string `db:"c_password"`
	Tokens   string `db:"c_tokens"`
}
type ExistingCustomer struct {
	ID int64 `db:"c_id"`
}

type Tokens struct {
	Token string `db:"g_tokens"`
}
