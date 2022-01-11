package refund

import (
	"errors"
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type PaymentService struct {
	logger    logger.Logger
	requestID string
}

const (
	RefundFailureCCNumber = "4000000000003238"
)

// NewPaymentService Public.
func NewPaymentService(logger logger.Logger) Payment {
	return &PaymentService{logger: logger}
}

// RefundPayment Public.
func (payWrap *PaymentService) RefundPayment(approveObj *models.Approve, refundAmount string) error {
	methodName := "RefundPayment"
	payWrap.logger.Debug(payWrap.requestID, "M:%v start", methodName)
	start := time.Now()

	if approveObj.CCNumber == RefundFailureCCNumber {
		return errors.New("refund failure")
	}

	payWrap.logger.Debug(payWrap.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return nil
}
