## 简介
此库是基于segmentio/kafka-go库对consumer的封装，使其在生产环境中简单易用

## 配置介绍

1. BrokerList                // kafka机器的IP+端口
2. TimeoutMS                 // 超时时间
3. Topic                     // 订阅的主题
4. GroupId                   // 组ID
5. ReadThreadLimit           // 开启consumer的数量
6. HandleFunc   func(string) // 消费数据的handler

## 使用方法
go get code.aliyun.com/module-go/imq

在代码中构造配置,可以指定QPS,并传入消费函数,关闭消费队列: kill -TERM PID

代码示例：
```
import (
	"fmt"
	"github.com/MoonlitAlley/easy-kafka"
	"time"
)

func main() {

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
```
