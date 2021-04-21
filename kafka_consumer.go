package imq

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MQConfig struct {
	BrokerList      []string // "10.141.10.66:9090,10.141.0.234:9090,10.141.39.56:9090"
	TimeoutMS       int      //  50
	Topic           string   // base-stats_topic
	GroupId         string   // group_1
	HandleFunc      func(string)
	ReadThreadLimit int64
}

func NewMQ(conf *MQConfig) error {

	if imqReadQPSLimit == 0 {
		imqReadQPSLimit = defaultImqReadQPSLimit
	}

	if conf.ReadThreadLimit == 0 {
		conf.ReadThreadLimit = defaultImqReadThreadLimit
	}

	for index := int64(0); index < conf.ReadThreadLimit; index++ {
		go func(handle func(string)) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("read handle func panic", err)
				}
			}()
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:        conf.BrokerList,
				GroupID:        conf.GroupId,
				Topic:          conf.Topic,
				MinBytes:       10e3, // 10KB
				MaxBytes:       10e6, // 10MB
				StartOffset:    kafka.FirstOffset,
				CommitInterval: time.Millisecond * 10,
			})

			signalChan := make(chan os.Signal)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
			tc := time.NewTimer(time.Microsecond * time.Duration(1000000/(imqReadQPSLimit/conf.ReadThreadLimit)))
			for {
				select {
				case <-signalChan:
					r.Close()
					log.Println("loop read stop by syscall.SIGTERM")
					return
				case <-tc.C:
					tc.Reset(time.Microsecond * time.Duration(1000000/(imqReadQPSLimit/conf.ReadThreadLimit)))
					m, err := r.ReadMessage(context.Background())
					if err != nil {
						log.Println("loop read message err", err)
						continue
					}
					go handle(string(m.Value))
				}
			}
		}(conf.HandleFunc)
	}

	log.Println("IMQ consumer up and running!...")
	return nil
}

const (
	defaultImqReadQPSLimit    int64 = 100
	defaultImqReadThreadLimit int64 = 10
)

var imqReadQPSLimit int64

// SetImqReadQPSLimit set qps limit in case of config update
func SetImqReadQPSLimit(limit int64) {
	if limit > 0 {
		imqReadQPSLimit = limit
	}
}
