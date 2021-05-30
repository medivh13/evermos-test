package repository

import (
	"github.com/medivh13/evermos-test/internal/models"
	"github.com/medivh13/evermos-test/pkg/dto"
)

type Repository interface {
	Register(data *models.Costumer) error
	Login(data *models.Costumer) (*dto.TokenRespDTO, error)
	GetAllProducts(tokens string) ([]*models.Product, error)
	CheckOut(data *dto.CheckOutReqDTO) (*models.CheckoutOrder, error)
	CancelExpiredOrder() error
	CancelByID(data *dto.CancelByIDReqDTO) error
}
