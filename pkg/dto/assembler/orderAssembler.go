package assembler

import (
	models "github.com/medivh13/evermos-test/internal/models"
	"github.com/medivh13/evermos-test/pkg/dto"
)

func ToCheckoutOrders(m *models.CheckoutOrder) *dto.CheckOutRespDTO {
	var status string
	if m.Status == 1 {
		status = "WAITING FOR PAYMENT"
	}
	return &dto.CheckOutRespDTO{
		OrderID: m.ID,
		Status:  status,
	}
}
