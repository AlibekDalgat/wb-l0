package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	wb_l0 "wb-l0"
	"wb-l0/pkg/service"
	mock_service "wb-l0/pkg/service/mocks"
)

func TestHandler_AddOrder(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrder, order wb_l0.Order)
	testTable := []struct {
		name         string
		inputBody    string
		inputOrder   wb_l0.Order
		mockBehavior mockBehavior
		expectedErr  error
	}{
		{
			name:      "OK",
			inputBody: `{"order_uid": "b563feb7b2b84b6test", "track_number": "WBILMTESTTRACK", "entry": "WBIL"}`,
			inputOrder: wb_l0.Order{
				OrderUID:    "b563feb7b2b84b6test",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
			},
			mockBehavior: func(s *mock_service.MockOrder, order wb_l0.Order) {
				s.EXPECT().AddOrder(order, []byte(`{"order_uid": "b563feb7b2b84b6test", "track_number": "WBILMTESTTRACK", "entry": "WBIL"}`)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:      "Wrong primary key order_uid param",
			inputBody: `{"uid": "b563feb7b2b84b6test", "track_number": "WBILMTESTTRACK", "entry": "WBIL"}`,
			inputOrder: wb_l0.Order{
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
			},
			mockBehavior: func(s *mock_service.MockOrder, order wb_l0.Order) {},
			expectedErr:  errors.New("Невалидный json"),
		},
		{
			name:         "Wrong input",
			inputBody:    `no json`,
			inputOrder:   wb_l0.Order{},
			mockBehavior: func(s *mock_service.MockOrder, order wb_l0.Order) {},
			expectedErr:  errors.New("Некорректное содержание сообщения в канале"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			order := mock_service.NewMockOrder(c)
			testCase.mockBehavior(order, testCase.inputOrder)

			services := &service.Service{Order: order}
			handler := NewHandler(services)

			orderData := []byte(testCase.inputBody)
			err := handler.AddOrder(orderData)

			assert.Equal(t, err, testCase.expectedErr)
		})
	}
}

func TestHandler_getOrder(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrder, id string)
	testTable := []struct {
		name            string
		id              string
		mockBehavior    mockBehavior
		expStatusCode   int
		expResponseBody string
	}{
		{
			name: "OK",
			id:   "b563feb7b2b84b6test",
			mockBehavior: func(s *mock_service.MockOrder, id string) {
				s.EXPECT().GetOrder(id).Return(wb_l0.Order{
					OrderUID:    "b563feb7b2b84b6test",
					TrackNumber: "WBILMTESTTRACK",
					Entry:       "WBIL",
				}, nil)
			},
			expStatusCode:   http.StatusOK,
			expResponseBody: `{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"","phone":"","zip":"","city":"","address":"","region":"","email":""},"payment":{"transaction":"","request_id":"","currency":"","provider":"","amount":0,"payment_dt":0,"bank":"","delivery_cost":0,"goods_total":0,"custom_fee":0},"items":null,"locale":"","internal_signature":"","customer_id":"","delivery_service":"","shardkey":"","sm_id":0,"date_created":"0001-01-01T00:00:00Z","oof_shard":""}`,
		},
		{
			name: "Not found",
			id:   "non",
			mockBehavior: func(s *mock_service.MockOrder, id string) {
				s.EXPECT().GetOrder(id).Return(wb_l0.Order{}, errors.New("Нет заказа с данным id: "+id))
			},
			expStatusCode:   http.StatusInternalServerError,
			expResponseBody: `{"message":"Нет заказа с данным id: non"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			order := mock_service.NewMockOrder(c)
			testCase.mockBehavior(order, testCase.id)

			services := &service.Service{Order: order}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/:id", handler.getOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/"+testCase.id, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expResponseBody)
		})
	}
}
