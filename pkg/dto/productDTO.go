package dto

import (
	"errors"

	"github.com/medivh13/evermos-test/pkg/common/crypto"
	"github.com/medivh13/evermos-test/pkg/common/env"
	"github.com/medivh13/evermos-test/pkg/common/validator"
	util "github.com/medivh13/evermos-test/pkg/utils"
)

type GetAllProductReqDTO struct {
	Tokens    string `json:"tokens" valid:"required" validname="Authorization"`
	Signature string `json:"signature" valid:"required" validname="signature"`
}

func (dto *GetAllProductReqDTO) Validate() error {
	v := validator.NewValidate(dto)
	v.SetCustomValidation(true, func() error {
		return dto.customValidation()
	})
	return v.Validate()
}

func (dto *GetAllProductReqDTO) customValidation() error {

	signature := crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), dto.Tokens)
	if signature != dto.Signature {
		if env.IsProduction() {
			return errors.New("invalid signature")
		}
		return errors.New("invalid signature" + " --> " + signature)
	}

	return nil
}

type GetAllProductRespDTO struct {
	ID   int64  `json:"product_id"`
	Name string `json:"product_name"`
	Desc string `json:"product_description"`
	QTY  int64  `json:"product_qty"`
}
