package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"os"

	"github.com/IBM/sarama"
)

func main() {
	keypair, err := tls.LoadX509KeyPair("./certs/service.cert", "./certs/service.key")
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.Open("./certs/ca.pem")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	caCert, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{keypair},
		RootCAs:      caCertPool,
	}

	// init config, enable error and notifications
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = tlsConfig
	config.Version = sarama.V2_5_0_0

	brokers := []string{"kafka-go-demo-johnleoclaudio-af60.d.aivencloud.com:17379"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Println("Failed to start Sarama producer:", err)
		return
	}

	defer producer.Close()

	producerMessage := &sarama.ProducerMessage{
		Topic: "demo",
		Key:   sarama.StringEncoder("key"),
		Value: sarama.StringEncoder("Hello, World!"),
	}

	partition, offset, err := producer.SendMessage(producerMessage)
	if err != nil {
		log.Println("Failed to send message:", err)
		return
	}
	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Println("Failed to start Sarama consumer:", err)
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("demo", 0, sarama.OffsetOldest)
	if err != nil {
		log.Println("Failed to start partition consumer:", err)
		return
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		log.Printf("Consumed message offset %d: %s\n", msg.Offset, string(msg.Value))
	}
}
