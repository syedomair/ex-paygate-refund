package refund

import (
	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// Payment Interface
type Payment interface {
	RefundPayment(approveObj *models.Approve, amount string) error
}
