package refund

import (
	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// Repository interface
type Repository interface {
	SetRequestID(requestID string)
	RefundApprove(inputApproveKey map[string]interface{}) (*models.Approve, error)
}
