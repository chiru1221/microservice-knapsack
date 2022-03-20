package main

import (
	"testing"

	"github.com/chiru1221/microservice-knapsack/client"
	pb "github.com/chiru1221/microservice-knapsack/knapsack"
	"github.com/golang/mock/gomock"
)

func TestServiceRun(t *testing.T) {
	type testCase struct {
		input      *pb.Input
		mockFn     func(*client.MockClient)
		wantResult int32
	}
	testCases := []testCase{
		{
			input: &pb.Input{
				N:      1,
				W:      1,
				Weight: []int32{1},
				Value:  []int32{1},
			},
			mockFn: func(m *client.MockClient) {
				for range [1]int{} {
					for range [1 + 1]int{} {
						m.EXPECT().CallUpdate(gomock.Any()).Return(&pb.Result{Result: 0}, nil)
					}
				}
			},
			wantResult: 0,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := client.NewMockClient(ctrl)
	service := &ServiceImpl{client: mockClient}
	for _, tc := range testCases {
		tc.mockFn(mockClient)
		result, _ := service.run(tc.input)
		if tc.wantResult != result.Result {
			t.Errorf("Expected: %d, but got: %d", tc.wantResult, result.Result)
		}
	}
}
