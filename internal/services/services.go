package services

import "github.com/medivh13/evermos-test/pkg/dto"

type Services interface {
	Register(req *dto.RegisterReqDTO) error
	Login(req *dto.RegisterReqDTO) (*dto.TokenRespDTO, error)
	GetAllProducts(req *dto.GetAllProductReqDTO) ([]*dto.GetAllProductRespDTO, error)
	ChekOut(req *dto.CheckOutReqDTO) (*dto.CheckOutRespDTO, error)
	CancelExpiredOrder() error
	CancelByID(req *dto.CancelByIDReqDTO) error
}
