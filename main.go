package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chiru1221/microservice-knapsack/client"
	pb "github.com/chiru1221/microservice-knapsack/knapsack"
)

type Service interface {
	run(input *pb.Input) (*pb.Result, error)
}

type ServiceImpl struct {
	client client.Client
}

type ServerImpl struct {
	service Service
}

func (service *ServiceImpl) run(input *pb.Input) (*pb.Result, error) {
	var (
		dp   []*pb.DP
		i, w int32
	)
	dp = make([]*pb.DP, input.GetN()+1)
	for idx := range dp {
		dp[idx] = &pb.DP{Row: make([]int32, input.GetW()+1)}
	}
	for i = 0; i < input.GetN(); i++ {
		for w = 0; w <= input.GetW(); w++ {
			result, _ := service.client.CallUpdate(
				&pb.KnapsackDP{
					I:      i,
					W:      w,
					Weight: input.GetWeight()[i],
					Value:  input.GetValue()[i],
					Dp:     dp,
				},
			)
			dp[i+1].Row[w] = result.Result
		}
	}
	return &pb.Result{Result: dp[input.N].Row[input.W]}, nil
}

func (server *ServerImpl) index(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	var input *pb.Input = &pb.Input{}
	if err = json.Unmarshal(data, input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	result, err := server.service.run(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	log.Println("Start server")
	service := &ServiceImpl{
		client: &client.ClientImpl{Addr: os.Getenv("SERVER_ADDR")},
	}
	if err := service.client.Construct(); err != nil {
		log.Fatal(err)
	}
	defer service.client.Destruct()

	server := &ServerImpl{
		service: service,
	}
	http.HandleFunc("/", server.index)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
