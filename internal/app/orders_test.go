package app

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tonydmorris/tranched/internal/app/mocks"
	"github.com/tonydmorris/tranched/internal/models"
)

// sample table test with mocks
func TestApp_createOrder(t *testing.T) {

	type args struct {
		order models.Order
	}
	tests := []struct {
		name             string
		mockExpectations func(mockUserRepository *mocks.MockUserRepository, mockOrderRepository *mocks.MockOrderRepository, mockLogger *mocks.MockLogger)
		args             args
		want             *models.Order
		wantErr          bool
	}{
		{
			name: "should return error if quantity is less than or equal to 0",
			args: args{
				order: models.Order{
					Quantity: 0,
				},
			},
			wantErr: true,
			mockExpectations: func(mockUserRepository *mocks.MockUserRepository, mockOrderRepository *mocks.MockOrderRepository, mockLogger *mocks.MockLogger) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			mockOrderRepository := mocks.NewMockOrderRepository(ctrl)
			mockLogger := mocks.NewMockLogger(ctrl)

			tt.mockExpectations(mockUserRepository, mockOrderRepository, mockLogger)

			a := &App{
				logger:          mockLogger,
				userRepository:  mockUserRepository,
				orderRepository: mockOrderRepository,
			}
			got, err := a.CreateOrder(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.createOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.createOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
