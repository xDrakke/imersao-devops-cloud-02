package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"os"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/KubeDev/imersao-devops-cloud-02/home-broker/01-worker/internal/infra/kafka"
	"github.com/KubeDev/imersao-devops-cloud-02/home-broker/01-worker/internal/market/dto"
	"github.com/KubeDev/imersao-devops-cloud-02/home-broker/01-worker/internal/market/entity"
	"github.com/KubeDev/imersao-devops-cloud-02/home-broker/01-worker/internal/market/transformer"
)

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")
	kafkaAutoOffsetReset := os.Getenv("KAFKA_AUTO_OFFSET_RESET")

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": kafkaServers,
		"group.id":          kafkaGroupId,
		"auto.offset.reset": kafkaAutoOffsetReset,
	}
	
	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMsgChan) // T2

	// recebe do canal do kafka, joga no input, processa joga no output e depois publica no kafka
	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() // T3

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
