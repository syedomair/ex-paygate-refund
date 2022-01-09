package refund

import (
	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// Payment Interface
type Payment interface {
	RefundPayment(*models.Approve) (string, error)
}
