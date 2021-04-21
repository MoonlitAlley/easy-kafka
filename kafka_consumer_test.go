package imq

import (
	"fmt"
	"testing"
	"time"
)

func TestImq(t *testing.T) {
	conf := &MQConfig{
		BrokerList: []string{"10.141.10.66:9090", "10.141.0.234:9090", "10.141.39.56:9090"},
		Topic:      "flume_ae_topic",
		GroupId:    "imq-new-version",
		HandleFunc: func(v string) {
			fmt.Println("Call:", v)
		}}
	err := NewMQ(conf)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 10)
	// sleep 10 seconds set qps to 1
	SetImqReadQPSLimit(1)
	// sleep 100 seconds to see the QPS limit update
	time.Sleep(time.Second * 100)
}
