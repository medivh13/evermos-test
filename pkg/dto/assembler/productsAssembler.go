package assembler

import (
	models "github.com/medivh13/evermos-test/internal/models"
	"github.com/medivh13/evermos-test/pkg/dto"
)

func ToGetProducts(m *models.Product) *dto.GetAllProductRespDTO {
	return &dto.GetAllProductRespDTO{
		ID:   m.ID,
		Name: m.Name,
		Desc: m.Desc,
		QTY:  m.QTY,
	}
}

func ToProducts(datas []*models.Product) []*dto.GetAllProductRespDTO {
	var ds []*dto.GetAllProductRespDTO
	for _, m := range datas {
		ds = append(ds, ToGetProducts(m))
	}
	return ds
}
