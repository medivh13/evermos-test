package dto

import (
	"errors"

	"github.com/medivh13/evermos-test/pkg/common/crypto"
	"github.com/medivh13/evermos-test/pkg/common/env"
	"github.com/medivh13/evermos-test/pkg/common/validator"
	util "github.com/medivh13/evermos-test/pkg/utils"
)

type CheckOutReqDTO struct {
	Tokens    string  `json:"tokens" valid:"required" validname="Authorization"`
	Signature string  `json:"signature" valid:"required" validname="signature"`
	ProductID []int64 `json:"product_id"`
	QTY       []int64 `json:"quantity"`
}

func (dto *CheckOutReqDTO) Validate() error {
	v := validator.NewValidate(dto)
	v.SetCustomValidation(true, func() error {
		return dto.customValidation()
	})
	return v.Validate()
}

func (dto *CheckOutReqDTO) customValidation() error {

	signature := crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), dto.Tokens)
	if signature != dto.Signature {
		if env.IsProduction() {
			return errors.New("invalid signature")
		}
		return errors.New("invalid signature" + " --> " + signature)
	}

	return nil
}

type CheckOutRespDTO struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

type OrderRespDTO struct {
	OrderID int64                  `json:"order_id"`
	Details []*OrderDetailsRespDTO `json:"details"`
}

type OrderDetailsRespDTO struct {
	Product string `json:"product_name"`
	QTY     int64  `json:"quantity"`
}

type CancelByIDReqDTO struct {
	Tokens    string `json:"tokens" valid:"required" validname="Authorization"`
	Signature string `json:"signature" valid:"required" validname="signature"`
	OrderID   int64  `json:"order_id"`
}

func (dto *CancelByIDReqDTO) Validate() error {
	v := validator.NewValidate(dto)
	v.SetCustomValidation(true, func() error {
		return dto.customValidation()
	})
	return v.Validate()
}

func (dto *CancelByIDReqDTO) customValidation() error {

	signature := crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), dto.Tokens)
	if signature != dto.Signature {
		if env.IsProduction() {
			return errors.New("invalid signature")
		}
		return errors.New("invalid signature" + " --> " + signature)
	}

	return nil
}
