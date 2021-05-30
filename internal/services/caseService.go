package services

import (
	"github.com/medivh13/evermos-test/internal/repository"
	"github.com/medivh13/evermos-test/pkg/dto"
	"github.com/medivh13/evermos-test/pkg/dto/assembler"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Services {
	return &service{repo}
}

func (s *service) Register(req *dto.RegisterReqDTO) error {

	err := req.Validate()
	if err != nil {
		return err
	}

	err = s.repo.Register(assembler.ToSaveRegister(req))
	if err != nil {
		return err
	}

	return err
}

func (s *service) Login(req *dto.RegisterReqDTO) (*dto.TokenRespDTO, error) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	data, err := s.repo.Login(assembler.ToSaveRegister(req))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) GetAllProducts(req *dto.GetAllProductReqDTO) ([]*dto.GetAllProductRespDTO, error) {

	err := req.Validate()

	if err != nil {
		return nil, err
	}

	datas, err := s.repo.GetAllProducts(req.Tokens)
	if err != nil {
		return nil, err
	}

	return assembler.ToProducts(datas), nil

}

func (s *service) ChekOut(req *dto.CheckOutReqDTO) (*dto.CheckOutRespDTO, error) {

	err := req.Validate()

	if err != nil {
		return nil, err
	}

	datas, err := s.repo.CheckOut(req)
	if err != nil {
		return nil, err
	}

	return assembler.ToCheckoutOrders(datas), nil

}

func (s *service) CancelExpiredOrder() error {

	err := s.repo.CancelExpiredOrder()
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CancelByID(req *dto.CancelByIDReqDTO) error {

	err := req.Validate()

	if err != nil {
		return err
	}

	err = s.repo.CancelByID(req)
	if err != nil {
		return err
	}
	return nil
}
