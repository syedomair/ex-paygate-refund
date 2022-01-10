package refund

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/mockserver"
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
	res, req := mockserver.MockTestServer(method, url, []byte(`{"approve_key":"1D754E20948F3EB8589A9"}`))
	c.RefundAction(res, req)
	response := new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected := "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Valid approve_key
	res, req = mockserver.MockTestServer(method, url, []byte(`{"approve_key":"D754E20948F3EB8589A9"}`))
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

func (mdb *mockPay) RefundPayment(approveObj *models.Approve) (string, error) {
	return "", nil
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
	if approveKey != "D754E20948F3EB8589A9" {
		return nil, errors.New("invalid approve_key")
	}
	return &models.Approve{}, nil
}
