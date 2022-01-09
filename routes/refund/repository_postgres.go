package refund

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type postgresRepo struct {
	client    *gorm.DB
	logger    logger.Logger
	requestID string
}

// NewPostgresRepository Public.
func NewPostgresRepository(c *gorm.DB, logger logger.Logger) Repository {
	return &postgresRepo{client: c, logger: logger, requestID: ""}
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

	approveObj := models.Approve{}
	if err := p.client.Table("approve").
		Where("approve_key = ?", approveKey).
		Find(&approveObj).Error; err != nil {
		return nil, errors.New("invalid approve_key")
	}

	inputApproveKey["status"] = 0

	if err := p.client.
		Table("approve").
		Where("approve_key = ?", approveKey).
		Updates(inputApproveKey).Error; err != nil {
		return nil, err
	}

	p.logger.Debug(p.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return &approveObj, nil
}
