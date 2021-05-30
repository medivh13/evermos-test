package assembler

import (
	models "github.com/medivh13/evermos-test/internal/models"
	"github.com/medivh13/evermos-test/pkg/dto"
)

func ToSaveRegister(d *dto.RegisterReqDTO) *models.Costumer {
	return &models.Costumer{
		Email:    d.Email,
		Password: d.Password,
		Tokens:   d.Tokens,
	}
}

func ToTokens(m *models.Tokens) *dto.TokenRespDTO {
	return &dto.TokenRespDTO{
		Token: m.Token,
	}
}
