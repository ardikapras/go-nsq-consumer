package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type myMessageHandler struct{}

// HandleMessage implements the Handler interface.
func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	fmt.Println("Start Execute m.Body::", string(m.Body))
	time.Sleep(10 * time.Second)
	fmt.Println("Finish Execute m.Body::", string(m.Body))

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func main() {
	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	config.MaxInFlight = 10
	consumer, err := nsq.NewConsumer("my_topic", "my_channel", config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	consumer.AddConcurrentHandlers(&myMessageHandler{}, 1)

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer.Stop()
}
