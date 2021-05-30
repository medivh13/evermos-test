package dto

import (
	"errors"

	"github.com/medivh13/evermos-test/pkg/common/crypto"
	"github.com/medivh13/evermos-test/pkg/common/env"
	"github.com/medivh13/evermos-test/pkg/common/validator"
	util "github.com/medivh13/evermos-test/pkg/utils"
)

type RegisterReqDTO struct {
	Email     string `json:"email" valid:"required,email" validname="email"`
	Tokens    string `json:"tokens, omitempty"`
	Password  string `json:"password" valid:"required" validname="password"`
	Signature string `json:"signature" valid:"required" validname="signature"`
}

func (dto *RegisterReqDTO) Validate() error {
	v := validator.NewValidate(dto)
	v.SetCustomValidation(true, func() error {
		return dto.customValidation()
	})
	return v.Validate()
}

func (dto *RegisterReqDTO) customValidation() error {

	signature := crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), dto.Email)
	if signature != dto.Signature {
		if env.IsProduction() {
			return errors.New("invalid signature")
		}
		return errors.New("invalid signature" + " --> " + signature)
	}

	return nil
}

type TokenRespDTO struct {
	Token string `json:"token"`
}
