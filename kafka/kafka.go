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
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/karldahlgren/raccoon/utility"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Create a new Kafka consumer
func CreateConsumer(topic string, group string, tracker *progress.Tracker) *kafka.Consumer {
	if group == "" {
		group = "raccoon-" + strconv.Itoa(rand.Int())
	}
	tracker.Increment(1)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost",
		"group.id":           group,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	})
	if err != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}
	tracker.Increment(1)
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
	messages := list.New()
	var count int64 = 0
	var index int64 = 0
	startTime := time.Now()
	for ; index < limit; index++ {
		msg, err := consumer.ReadMessage(-1)

		if err == nil {
			key := strings.ToLower(string(msg.Key))
			value := strings.ToLower(string(msg.Value))
			if (keyQuery != "" && strings.Contains(key, keyQuery)) ||
				(valueQuery != "" && strings.Contains(value, valueQuery)) {
				count++
				message := Message{
					Key:       string(msg.Key),
					Value:     value,
					Timestamp: msg.Timestamp,
					Partition: msg.TopicPartition.Partition,
					Offset:    msg.TopicPartition.Offset.String(),
				}

				messages.PushFront(message)
			}
		}
		tracker.Increment(1)
	}
	stopTime := time.Now()
	elapsedTime := stopTime.Sub(startTime)
	tracker.MarkAsDone()
	return Result{
		Messages: *messages,
		Count:    count,
		Duration: elapsedTime,
	}
}
