package refund

import (
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type PaymentService struct {
	logger    logger.Logger
	requestID string
}

// NewPaymentService Public.
func NewPaymentService(logger logger.Logger) Payment {
	return &PaymentService{logger: logger}
}

// RefundPayment Public.
func (payWrap *PaymentService) RefundPayment(approveObj *models.Approve, refundAmount string) error {
	methodName := "RefundPayment"
	payWrap.logger.Debug(payWrap.requestID, "M:%v start", methodName)
	start := time.Now()

	/*
			4000 0000 0000 0119: authorisation failure
			4000 0000 0000 0259: capture failure
		        4000 0000 0000 3238: refund failure
	*/

	payWrap.logger.Debug(payWrap.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return nil
}
