package repository

import (
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestOrderPostgres_AddOrder(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Ошибка при тестированном подключении к БД: %s", err)
	}
	defer db.Close()

	r := NewOrderPostgres(db)

	type args struct {
		orderUID string
		data     []byte
	}
	testTable := []struct {
		name      string
		mock      func()
		input     args
		expectErr bool
	}{
		{
			name: "OK",
			input: args{
				orderUID: "b563feb7b2b84b6test",
				data:     []byte(`{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"","phone":"","zip":"","city":"","address":"","region":"","email":""},"payment":{"transaction":"","request_id":"","currency":"","provider":"","amount":0,"payment_dt":0,"bank":"","delivery_cost":0,"goods_total":0,"custom_fee":0},"items":null,"locale":"","internal_signature":"","customer_id":"","delivery_service":"","shardkey":"","sm_id":0,"date_created":"0001-01-01T00:00:00Z","oof_shard":""}`),
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO "+itemsTable).WithArgs("b563feb7b2b84b6test", []byte(`{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"","phone":"","zip":"","city":"","address":"","region":"","email":""},"payment":{"transaction":"","request_id":"","currency":"","provider":"","amount":0,"payment_dt":0,"bank":"","delivery_cost":0,"goods_total":0,"custom_fee":0},"items":null,"locale":"","internal_signature":"","customer_id":"","delivery_service":"","shardkey":"","sm_id":0,"date_created":"0001-01-01T00:00:00Z","oof_shard":""}`)).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			err := r.AddOrder(testCase.input.orderUID, testCase.input.data)
			if testCase.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
