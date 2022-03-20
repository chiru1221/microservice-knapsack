package main

import (
	"context"
	"testing"

	pb "github.com/chiru1221/microservice-knapsack/knapsack"
)

func TestUpdate(t *testing.T) {
	type testCase struct {
		input      *pb.KnapsackDP
		wantResult int32
	}
	testCases := []testCase{
		{
			input: &pb.KnapsackDP{
				I:      0,
				W:      0,
				Weight: 0,
				Value:  0,
				Dp: []*pb.DP{
					{Row: []int32{0}},
					{Row: []int32{0}},
				},
			},
			wantResult: 0,
		},
	}
	srv := serverImpl{}
	for _, tc := range testCases {
		result, _ := srv.Update(context.Background(), tc.input)
		if tc.wantResult != result.Result {
			t.Errorf("Expected: %d, but got: %d", tc.wantResult, result.Result)
		}
	}
}
