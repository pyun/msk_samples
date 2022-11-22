/**
 * @Author: mcg
 * @Date: 2021/3/1 20:41
 * @Desc:
 */

package service

import (
	"context"
	"crypto/tls"
	"fmt"
	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam_v2"
	"time"
)

var (
	//paintext cluster
	bootstrapServersTxt = []string{"*.kafka.ap-southeast-1.amazonaws.com:9092", "*.kafka.ap-southeast-1.amazonaws.com:9092", "*.kafka.ap-southeast-1.amazonaws.com:9092"}
	//iam tls cluste
	bootstrapServersTls = []string{"*.kafka.ap-southeast-1.amazonaws.com:9098", "*.kafka.ap-southeast-1.amazonaws.com:9098"}
	topic               = "test-topic1"
	region              = "ap-southeast-1"
)

func Start(iam bool) {
	if iam {
		producerIamTLS()
		consumerTLS()
	} else {
		producerDirect()
		consumerDirect()
	}

}
func producerDirect() {
	fmt.Println("begin producing plaintext msg …… ")
	config := kafka.WriterConfig{
		Brokers: bootstrapServersTxt,
		Topic:   topic,
	}
	w := kafka.NewWriter(config)

	err := w.WriteMessages(context.Background(), kafka.Message{Value: []byte("this is test msg")})

	if err != nil {
		fmt.Println("producing msg err: ", err.Error())
	} else {
		fmt.Println("msg  is sent successfully! ")
	}
}
func producerIamTLS() {
	fmt.Println("begin producing tlstext msg …… ")
	cfg, _ := config.LoadDefaultConfig(context.Background())
	m := &aws_msk_iam_v2.Mechanism{
		Signer:      signer.NewSigner(),
		Credentials: cfg.Credentials,
		Region:      region,
		SignTime:    time.Now(),
		Expiry:      time.Minute * 15,
	}

	config := kafka.WriterConfig{
		Brokers: bootstrapServersTls,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			Timeout:       50 * time.Second,
			DualStack:     true,
			SASLMechanism: m,
			TLS: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	w := kafka.NewWriter(config)
	//fmt.Println("Consumer configuration: ", config)

	err := w.WriteMessages(context.TODO(), kafka.Message{Value: []byte("this is a tls msg")})

	if err != nil {
		fmt.Println("producing msg err: ", err.Error())
	} else {
		fmt.Println("msg is sent successfully! ")
	}

}

func consumerDirect() {
	readercfg := kafka.ReaderConfig{
		Brokers:     bootstrapServersTxt,
		GroupID:     "some-consumer-group",
		GroupTopics: []string{topic},
	}
	r := kafka.NewReader(readercfg)
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading message: %+v\n", err)
			break
		}
		fmt.Printf("Received message from %s-%d [%d]: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

func consumerTLS() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	m := &aws_msk_iam_v2.Mechanism{
		Signer:      signer.NewSigner(),
		Credentials: cfg.Credentials,
		Region:      region,
		SignTime:    time.Now(),
		Expiry:      time.Minute * 15,
	}

	config := kafka.ReaderConfig{
		Brokers: bootstrapServersTls,
		GroupID: "test-consumer-group-1",
		Topic:   topic,
		// Partition: 0,
		MaxWait: 50000 * time.Millisecond,
		Dialer: &kafka.Dialer{
			Timeout:       50 * time.Second,
			DualStack:     true,
			SASLMechanism: m,
			TLS: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	r := kafka.NewReader(config)
	//fmt.Println("Consumer configuration: ", config)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading message: %+v\n", err)
			break
		}
		fmt.Printf("Received message from %s-%d [%d]: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
