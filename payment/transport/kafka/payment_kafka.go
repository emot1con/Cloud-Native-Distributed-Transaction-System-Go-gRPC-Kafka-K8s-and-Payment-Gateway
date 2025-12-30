package kafka

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"payment/proto"
	"payment/service"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ConsumerHandler struct {
	service *service.PaymentService
}

func connectKafka(addr []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	var CGError error

	for i := 0; i < 5; i++ {
		worker, err := sarama.NewConsumerGroup(addr, groupID, config)
		if err == nil {
			logrus.Info("Connected to kafka")
			return worker, nil
		}

		CGError = err
		logrus.Warnf("failed connect to kafka, retrying...(%d/5)", i+1)
		time.Sleep(5 * time.Second)
		continue
	}

	return nil, CGError
}

func ProcessMessage(addr []string, topic []string, groupID string, service *service.PaymentService) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumerGroup, err := connectKafka(addr, groupID)
	if err != nil {
		logrus.Fatalf("error when connect to kafka: %v", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := &ConsumerHandler{service: service}

	go func() {
		logrus.Infof("addr: %s topic: %s groupID: %s", addr, topic, groupID)
		for {
			select {
			case <-ctx.Done():
				logrus.Info("Consumer stopping...")
				return
			default:
				if err := consumerGroup.Consume(ctx, topic, handler); err != nil {
					logrus.Errorf("failed when consume partition, retrying: %v", err)
					time.Sleep(2 * time.Second)
				}
			}
		}
	}()

	<-sigchan
	logrus.Info("shutting down consumer...")
	cancel()
}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		order := new(proto.Order)
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			logrus.Errorf("Error parsing message: %s. Raw data: %s", err, string(msg.Value))
			continue
		}
		logrus.Infof("Received message, UserID: %d with OrderId: %d \n", order.UserId, order.Id)

		// Create payment record with new schema (amount instead of total_price)
		response, err := h.service.CreatePayment(&proto.CreatePaymentRequest{
			OrderId: order.Id,
			Amount:  order.TotalPrice,
		})
		if err != nil {
			logrus.Errorf("Error creating payment: %v", err)
			continue
		}

		logrus.Infof("Payment created with ID: %d for order: %d", response.Id, order.Id)

		sess.MarkMessage(msg, "")
	}

	return nil
}
