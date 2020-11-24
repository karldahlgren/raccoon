/*
 * The MIT License
 *
 * Copyright (c) 2020-, Karl A. Dahlgren
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package kafka

import (
	"container/list"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/karldahlgren/raccoon/utility"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Create a new earliest Kafka consumer
func CreateEarliestConsumer(broker string, topic string, group string, tracker *progress.Tracker) *kafka.Consumer {
	return createConsumer(broker, topic, group, "earliest", tracker)
}


// Create a new earliest Kafka consumer
func CreateLatestConsumer(broker string, topic string, group string, tracker *progress.Tracker) *kafka.Consumer {
	return createConsumer(broker, topic, group, "latest", tracker)
}

func createConsumer(broker string, topic string, group string, offset string, tracker *progress.Tracker) (*kafka.Consumer) {
	tracker.Increment(1)
	if group == "" {
		group = "raccoon-" + strconv.Itoa(rand.Int())
	}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           group,
		"auto.offset.reset":  offset,
		"enable.auto.commit": "false",
	})
	tracker.Increment(1)
	if consumer != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}

	consumer.SubscribeTopics([]string{topic}, nil)

	tracker.Increment(1)
	return consumer
}


// Stop a Kafka consumer
func StopConsumer(consumer *kafka.Consumer, tracker *progress.Tracker) {
	tracker.Increment(1)
	err := consumer.Close()

	if err != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}
	tracker.Increment(1)
}

// Consume messages from a Kafka consumer
func Consume(consumer *kafka.Consumer, keyQuery string, valueQuery string, limit int64, tracker *progress.Tracker) Result {
	// Set the tracker length to limit + 1 since we otherwise get
	// invalid formatting for the tracker
	tracker.Total = limit + 1

	// List of messages
	messages := list.New()

	var matchedMessages int64 = 0
	var readMessages int64 = 0
	var index int64 = 0
	startTime := time.Now()
	for ; index < limit; index++ {
		msg, err := consumer.ReadMessage(-1)

		if err == nil {
			readMessages++
			message := parseMessage(msg, keyQuery, valueQuery)
			if message != nil {
				matchedMessages++
				tracker.Message = "Reading messages (" + strconv.FormatInt(matchedMessages, 10) + " matches)"
				messages.PushFront(message)
			}
		}
		tracker.Increment(1)
	}
	stopTime := time.Now()
	elapsedTime := stopTime.Sub(startTime)

	// Sleep for the progress to catch up
	time.Sleep(100 * time.Millisecond)
	tracker.MarkAsDone()

	return Result{
		Messages: *messages,
		MatchedMessages:    matchedMessages,
		ReadMessages: readMessages,
		Duration: elapsedTime,
	}
}

// Tail messages from a Kafka consumer
func Tail(consumer *kafka.Consumer, keyQuery string, valueQuery string, limit int64, tracker *progress.Tracker) Result {
	messages := list.New()
	var matchedMessages int64 = 0
	var readMessages int64 = 0
	var index int64 = 0
	startTime := time.Now()
	running := true
	go func() {
	if limit != -1 {
		for ; index < limit && running; index++ {
			msg, err := consumer.ReadMessage(-1)

			if err == nil {
				readMessages++
				message := parseMessage(msg, keyQuery, valueQuery)
				if message != nil {
					matchedMessages++
					tracker.Message = "Reading messages (" + strconv.FormatInt(matchedMessages, 10) + " matches)"
					messages.PushFront(message)
				}
			}
			tracker.Increment(1)
		}
	} else {
		// No limit
			for running {
				msg, err := consumer.ReadMessage(1000)

				if err == nil {
					readMessages++
					message := parseMessage(msg, keyQuery, valueQuery)
					if message != nil {
						matchedMessages++
						tracker.Message = "Reading messages (" + strconv.FormatInt(matchedMessages, 10) + " matches)"
						messages.PushFront(message)
					}
				}
			}
		}
	}()

	fmt.Scanln()
	running = false

	stopTime := time.Now()
	elapsedTime := stopTime.Sub(startTime)

	// Sleep for the progress to catch up
	time.Sleep(100 * time.Millisecond)
	tracker.MarkAsDone()

	return Result{
		Messages: *messages,
		MatchedMessages:    matchedMessages,
		ReadMessages: readMessages,
		Duration: elapsedTime,
	}
}

func parseMessage(message *kafka.Message, keyQuery string, valueQuery string) *Message  {
	key := strings.ToLower(string(message.Key))
	value := strings.ToLower(string(message.Value))
	if (keyQuery != "" && strings.Contains(key, keyQuery)) ||
		(valueQuery != "" && strings.Contains(value, valueQuery)) {
		return &Message{
			Key:       string(message.Key),
			Value:     value,
			Timestamp: message.Timestamp,
			Partition: message.TopicPartition.Partition,
			Offset:    message.TopicPartition.Offset.String(),
		}
	}

	return nil
}


