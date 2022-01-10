package refund

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/sync/errgroup"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/floathelper"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type postgresRepo struct {
	client     *gorm.DB
	logger     logger.Logger
	requestID  string
	payService Payment
}

const (
	Refund = "REFUND"
)

// NewPostgresRepository Public.
func NewPostgresRepository(c *gorm.DB, logger logger.Logger, payService Payment) Repository {
	return &postgresRepo{client: c, logger: logger, requestID: "", payService: payService}
}

func (p *postgresRepo) SetRequestID(requestID string) {
	p.requestID = requestID
}

// RefundApprove Public
func (p *postgresRepo) RefundApprove(inputApproveKey map[string]interface{}) (*models.Approve, error) {
	methodName := "RefundApprove"
	p.logger.Debug(p.requestID, "M:%v start", methodName)
	start := time.Now()

	approveKey := ""
	if approveKeyValue, ok := inputApproveKey["approve_key"]; ok {
		approveKey = approveKeyValue.(string)
	}

	amount := ""
	if amountValue, ok := inputApproveKey["amount"]; ok {
		amount = amountValue.(string)
	}

	err := p.transactionCreateLedger(p.client, approveKey, amount)
	if err != nil {
		return nil, err
	}

	approveObj := models.Approve{}
	if err := p.client.Table("approve").
		Where("approve_key = ?", approveKey).
		Find(&approveObj).Error; err != nil {
		return nil, errors.New("invalid approve_key")
	}

	p.logger.Debug(p.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return &approveObj, nil
}

// transactionCreateLedger
func (p *postgresRepo) transactionCreateLedger(db *gorm.DB,
	approveKey, amount string) error {
	methodName := "transactionCreateLedger"
	p.logger.Debug(p.requestID, "M:%v start", methodName)
	start := time.Now()

	return db.Transaction(func(tx *gorm.DB) error {

		approveObj := models.Approve{}
		if err := p.client.Set("gorm:query_option", "FOR UPDATE").Table("approve").
			Where("approve_key = ?", approveKey).
			Where("status = ?", 1).
			Find(&approveObj).Error; err != nil {
			return errors.New("invalid approve_key")
		}

		var amountBalance, amountRefund float64
		var err error

		if amountBalance, err = strconv.ParseFloat(approveObj.AmountBalance, 64); err != nil {
			return errors.New("invalid amount balance")
		}

		if amountRefund, err = strconv.ParseFloat(amount, 64); err != nil {
			return errors.New("invalid amount balance")
		}

		f := floathelper.Floater{Accuracy: 0.01}
		if f.AGreaterThanB(amountRefund, amountBalance) == 1 {
			return errors.New("invalid refund amount")
		}

		amountRefund := fmt.Sprintf("%f", amountRefund)

		err = p.payService.RefundPayment(&approveObj, amountRefundStr)
		if err != nil {
			return errors.New("error from payment service")
		}

		g := new(errgroup.Group)
		g.Go(func() error {
			newLedger := &models.Ledger{}
			newLedger.MerchantID = approveObj.MerchantID
			newLedger.ApproveID = approveObj.ID
			newLedger.Amount = amountRefundStr
			newLedger.ActionType = Refund
			newLedger.CreatedAt = time.Now().Format(time.RFC3339)
			if err := p.client.Create(newLedger).Error; err != nil {
				return err
			}
			return nil
		})

		g.Go(func() error {
			inputApproveKey := make(map[string]interface{})
			inputApproveKey["amount_balance"] = amountBalance - amountRefund
			if err := p.client.
				Table("approve").
				Where("approve_key = ?", approveKey).
				Updates(inputApproveKey).Error; err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return err
		}

		p.logger.Debug(p.requestID, "M:%v ts %+v", methodName, time.Since(start))
		return nil
	})
}
