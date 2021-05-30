package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/medivh13/evermos-test/internal/models"
	"github.com/medivh13/evermos-test/internal/repository"
	"github.com/medivh13/evermos-test/pkg/dto"
	"github.com/medivh13/evermos-test/pkg/dto/assembler"
	btbErrors "github.com/medivh13/evermos-test/pkg/errors"
	util "github.com/medivh13/evermos-test/pkg/utils"
)

const (
	Register          = `INSERT INTO customers (c_email, c_password) VALUES ($1, $2)`
	CekEmail          = `SELECT c_id from customers where c_email = $1`
	CekPass           = `SELECT c_id from customers where c_password = $1`
	UpdateToken       = `UPDATE customers set c_tokens = '%s' where c_id = %d returning c_tokens`
	CekToken          = `SELECT c_id from customers where c_tokens = $1`
	GetAllProduct     = `SELECT p_id, p_name, p_desc, p_stocks from products`
	GetProductByID    = `SELECT p_id, p_name, p_desc, p_stocks from products where p_id = $1`
	CheckoutStock     = `UPDATE products SET p_stocks = p_stocks - $1 where p_id = $2`
	CheckinStock      = `UPDATE products SET p_stocks = p_stocks + $1 where p_id = $2`
	InsertOrder       = `INSERT INTO orders (o_c_id, o_status, o_due_date) VALUES ($1, 1, $2) RETURNING o_id, o_status`
	InsertOrderDetail = `INSERT INTO order_details (od_o_id, od_p_id, od_qty) VALUES ($1, $2, $3)`
	CancelAllOrder    = `UPDATE orders set o_status = 0 where (now() - o_due_date >= interval '1 day')`
	CancelOrderByID   = `UPDATE orders set o_status = 0 where o_id = $1`
	GetExpiredOrder   = `SELECT a.o_id, b.od_p_id, b.od_qty from orders a INNER JOIN order_details b ON a.o_id = b.od_o_id where (now() - a.o_due_date >= interval '1 day') AND a.o_status = 1`
	GetOrderToProcess = `SELECT a.o_id, b.od_p_id, b.od_qty from orders a INNER JOIN order_details b ON a.o_id = b.od_o_id where a.o_status = 1 AND a.o_id = $1`
)

// order status -> 0 = CANCEL, 1 = CHECKOUT / WAITING FOR PAYMENT, 2 = PAID
var statement PreparedStatement

type PreparedStatement struct {
	register          *sqlx.Stmt
	cekemail          *sqlx.Stmt
	cekpass           *sqlx.Stmt
	updatetoken       *sqlx.Stmt
	cektoken          *sqlx.Stmt
	getAllProduct     *sqlx.Stmt
	getProductByID    *sqlx.Stmt
	cancelAllOrder    *sqlx.Stmt
	getOrderToProcess *sqlx.Stmt
	getExpiredOrder   *sqlx.Stmt
}

type PostgresRepo struct {
	Conn *sqlx.DB
}

func NewPostgresRepo(Conn *sqlx.DB) repository.Repository {

	repo := &PostgresRepo{Conn}
	InitPreparedStatement(repo)
	return repo
}

func (m *PostgresRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := m.Conn.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *PostgresRepo) {
	statement = PreparedStatement{
		register:          m.Preparex(Register),
		cekemail:          m.Preparex(CekEmail),
		cekpass:           m.Preparex(CekPass),
		cektoken:          m.Preparex(CekToken),
		getAllProduct:     m.Preparex(GetAllProduct),
		getProductByID:    m.Preparex(GetProductByID),
		cancelAllOrder:    m.Preparex(CancelAllOrder),
		getOrderToProcess: m.Preparex(GetOrderToProcess),
		getExpiredOrder:   m.Preparex(GetExpiredOrder),
	}
}

func (m *PostgresRepo) Register(data *models.Costumer) error {
	gID := []models.ExistingCustomer{}

	err := statement.cekemail.Select(&gID, data.Email)
	if err != nil {
		log.Println("Failed Query to get Existing Email : ", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(gID) > 0 {
		return fmt.Errorf(btbErrors.ErrorExistingData)
	}

	_, err = statement.register.Exec(data.Email, data.Password)

	if err != nil {
		log.Println("Failed Query Register : ", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	return nil
}

func (m *PostgresRepo) Login(data *models.Costumer) (*dto.TokenRespDTO, error) {
	gID := []models.ExistingCustomer{}
	tokens := models.Tokens{}
	err := statement.cekpass.Select(&gID, data.Password)
	if err != nil {
		log.Println("Failed Query to get Existing Customer : ", err.Error())
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(gID) < 1 {

		return nil, fmt.Errorf(btbErrors.ErrorUserNotFound)
	}

	data.Tokens, _ = util.CreateToken(uint32(gID[0].ID))
	query := fmt.Sprintf(UpdateToken, data.Tokens, gID[0].ID)
	err = m.Conn.QueryRow(query).Scan(&tokens.Token)

	if err != nil {
		log.Println("Failed Query Login : ", err.Error())
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}

	return assembler.ToTokens(&tokens), nil
}

func (m *PostgresRepo) GetAllProducts(tokens string) ([]*models.Product, error) {
	gID := []models.ExistingCustomer{}

	err := statement.cektoken.Select(&gID, tokens)
	if err != nil {
		log.Println("Failed Query to get Existing Customer : check - tokens", err.Error())
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(gID) < 1 {

		return nil, fmt.Errorf(btbErrors.ErrorUserNotFound)
	}

	var data []*models.Product

	err = statement.getAllProduct.Select(&data)
	if err != nil {
		log.Println("Failed Query Products - getProducts: ", err.Error())
		return data, fmt.Errorf(btbErrors.ErrorDB)
	}
	if len(data) == 0 {
		return data, fmt.Errorf(btbErrors.ErrorDataNotFound)
	}
	return data, nil
}

func (m *PostgresRepo) CheckOut(data *dto.CheckOutReqDTO) (*models.CheckoutOrder, error) {
	gID := []models.ExistingCustomer{}

	err := statement.cektoken.Select(&gID, data.Tokens)
	if err != nil {
		log.Println("Failed Query to get Existing Customer : check - tokens", err.Error())
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(gID) < 1 {

		return nil, fmt.Errorf(btbErrors.ErrorUserNotFound)
	}

	var dataOrder models.CheckoutOrder

	tx, err := m.Conn.Beginx()
	if err != nil {
		log.Println("Failed to start database transaction. Error: %s", err.Error())
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}
	now := time.Now()
	err = tx.QueryRow(InsertOrder, gID[0].ID, now.AddDate(0, 0, 1)).Scan(&dataOrder.ID, &dataOrder.Status)

	if err != nil {
		log.Println("Failed to start database transaction. Error: %s", err.Error())
		tx.Rollback()
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}

	for key, val := range data.ProductID {
		_, err = tx.Exec(InsertOrderDetail, dataOrder.ID, val, data.QTY[key])
		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return nil, err
		}

		_, err = tx.Exec(CheckoutStock, data.QTY[key], val)
		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Failed to commit database transaction. Error: %s", err.Error())
	}

	return &dataOrder, nil
}

func (m *PostgresRepo) CancelExpiredOrder() error {

	dataOrder := []models.OrderToProcess{}

	err := statement.getExpiredOrder.Select(&dataOrder)
	if err != nil {
		log.Println("Failed Query to get Expired Order", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(dataOrder) < 1 {

		return fmt.Errorf(btbErrors.ErrorDataNotFound)
	}

	tx, err := m.Conn.Beginx()
	if err != nil {
		log.Println("Failed to start database transaction. Error: %s", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	for _, val := range dataOrder {

		_, err = tx.Exec(CancelOrderByID, val.ID)

		if err != nil {
			log.Println("Failed to start database transaction. Error: %s", err.Error())
			tx.Rollback()
			return fmt.Errorf(btbErrors.ErrorDB)
		}

		_, err = tx.Exec(CheckinStock, val.Qty, val.ProductID)
		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit database transaction. Error: %s", err.Error())
	}

	return nil
}

func (m *PostgresRepo) CancelByID(data *dto.CancelByIDReqDTO) error {
	gID := []models.ExistingCustomer{}

	err := statement.cektoken.Select(&gID, data.Tokens)
	if err != nil {
		log.Println("Failed Query to get Existing Customer : check - tokens", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(gID) < 1 {

		return fmt.Errorf(btbErrors.ErrorUserNotFound)
	}

	dataOrder := []models.OrderToProcess{}

	err = statement.getOrderToProcess.Select(&dataOrder, data.OrderID)
	if err != nil {
		log.Println("Failed Query to get Existing Order", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(dataOrder) < 1 {

		return fmt.Errorf(btbErrors.ErrorDataNotFound)
	}

	tx, err := m.Conn.Beginx()
	if err != nil {
		log.Println("Failed to start database transaction. Error: %s", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	_, err = tx.Exec(CancelOrderByID, data.OrderID)

	if err != nil {
		log.Println("Failed to start database transaction. Error: %s", err.Error())
		tx.Rollback()
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	for _, val := range dataOrder {
		_, err = tx.Exec(CheckinStock, val.Qty, val.ProductID)
		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit database transaction. Error: %s", err.Error())
	}

	return nil
}
