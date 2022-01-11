package refund

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/mockserver"
)

const (
	ValidApproveKey   = "06F3BCC1C3B836B1AA6D"
	InvalidApproveKey = "1D754E20948F3EB8589A9"
)

func TestRefundAction(t *testing.T) {
	c := Controller{
		Logger: logger.New("DEBUG", "TEST#", os.Stdout),
		Repo:   &mockDB{},
		Pay:    &mockPay{}}

	method := "POST"
	url := "/refund"

	type TestResponse struct {
		Data   string
		Result string
	}

	//Invalid approve_key
	res, req := mockserver.MockTestServer(method, url, []byte(`{"amount":"2", "approve_key":"`+InvalidApproveKey+`"}`))
	c.RefundAction(res, req)
	response := new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected := "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Valid approve_key
	res, req = mockserver.MockTestServer(method, url, []byte(`{"amount":"10", "approve_key":"`+ValidApproveKey+`"}`))
	c.RefundAction(res, req)
	response = new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected = "success"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
}

type mockPay struct {
}

func (mdb *mockPay) RefundPayment(approveObj *models.Approve, refundAmount string) error {
	if approveObj.CCNumber == RefundFailureCCNumber {
		return errors.New("refund failure")
	}
	return nil
}

type mockDB struct {
}

func (mdb *mockDB) SetRequestID(requestID string) {
}

func (mdb *mockDB) RefundApprove(inputApproveKey map[string]interface{}) (*models.Approve, error) {
	approveKey := ""
	if approveKeyValue, ok := inputApproveKey["approve_key"]; ok {
		approveKey = approveKeyValue.(string)
	}
	if approveKey != ValidApproveKey {
		return nil, errors.New("invalid approve_key")
	}
	return &models.Approve{}, nil
}
