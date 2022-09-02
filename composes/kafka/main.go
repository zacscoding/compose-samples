package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	brokers = []string{"localhost:9092"}
	topic   = "sample-message"
)

func main() {
	if err := setupTopic(); err != nil {
		panic(err)
	}

	msgWg := sync.WaitGroup{}
	consumerReady := make(chan struct{})
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.Partitioner = sarama.NewManualPartitioner
	conf.Consumer.Offsets.AutoCommit.Enable = true
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	cli, err := sarama.NewClient(brokers, conf)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	go func() {
		consumer, err := sarama.NewConsumerGroupFromClient("group1", cli)
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		topics := []string{topic}
		handler := &MessageConsumer{wg: &msgWg, readyCh: consumerReady}
		for {
			err := consumer.Consume(ctx, topics, handler)
			if err == sarama.ErrClosedClient || err == sarama.ErrClosedConsumerGroup {
				log.Print("[Consumer] Stopping")
				return
			}
			log.Printf("[Consumer] failed to consume. err: %v", err)
		}
	}()

	<-consumerReady

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("message-%d", i+1)),
		})
		if err != nil {
			log.Printf("failed to produce message. err: %v", err)
			continue
		}
		msgWg.Add(1)
	}
	msgWg.Wait()
	cli.Close()
}

func setupTopic() error {
	admin, err := sarama.NewClusterAdmin(brokers, sarama.NewConfig())
	if err != nil {
		return fmt.Errorf("new cluster admin: %v", err)
	}
	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		terr, ok := err.(*sarama.TopicError)
		if !ok {
			return err
		}
		if terr.Err != sarama.ErrTopicAlreadyExists {
			return err
		}
		log.Printf("Skip to create a new topic %s because already exists", topic)
	} else {
		log.Printf("Success to create a new topic: %s", topic)
	}
	return nil
}

type MessageConsumer struct {
	wg      *sync.WaitGroup
	readyCh chan struct{}
}

func (c *MessageConsumer) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("[Consumer] Setup Session. memberid:%s", session.MemberID())
	close(c.readyCh)
	return nil
}

func (c *MessageConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Printf("[Consumer] Cleanup Session. memberid:%s", session.MemberID())
	return nil
}

func (c *MessageConsumer) ConsumeClaim(_ sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("consume message: %s", string(msg.Value))
		c.wg.Done()
	}
	return nil
}
