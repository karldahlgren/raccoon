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
func CreateEarliestConsumer(bootstrap string, topic string, group string, tracker *progress.Tracker) *kafka.Consumer {
	return createConsumer(bootstrap, topic, group, "earliest", tracker)
}

// Create a new earliest Kafka consumer
func CreateLatestConsumer(bootstrap string, topic string, group string, tracker *progress.Tracker) *kafka.Consumer {
	return createConsumer(bootstrap, topic, group, "latest", tracker)
}

// Stop a Kafka consumer
func StopConsumer(consumer *kafka.Consumer, tracker *progress.Tracker) {
	err := consumer.Close()

	if err != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}
	tracker.MarkAsDone()
}

// Get all partitions for a particular topic
func GetPartitions(consumer *kafka.Consumer, topic string, tracker *progress.Tracker) []Partition {
	metaData, err := consumer.GetMetadata(&topic, false, -1)

	if err != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}

	topicMetaData := metaData.Topics[topic]

	if &topicMetaData == nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}

	tracker.Total = int64(len(topicMetaData.Partitions))
	var partitions []Partition
	for index := 0; index < len(topicMetaData.Partitions); index++ {
		partition := topicMetaData.Partitions[index]
		lowOffset, highOffset, err := consumer.QueryWatermarkOffsets(topic, partition.ID, -1)

		if err != nil {
			tracker.MarkAsDone()
			utility.ExitOnError(err)
		}

		partitions = append(partitions, Partition{
			id: partition.ID,
			lowOffset: lowOffset,
			highOffset: highOffset,
		})
		tracker.Increment(1)
	}

	tracker.MarkAsDone()
	return partitions
}

// Consume messages from a Kafka consumer
func Consume(consumer *kafka.Consumer, partitions *[]Partition, keyQuery string, valueQuery string, limit int64, tracker *progress.Tracker) Result {
	// List of messages
	messages := list.New()

	// Calculate the limit for each partition
	limitByPartition := getMessageLimitByPartition(partitions, limit)

	// Initialize a counter for each partition
	counterByPartition := getCounterByPartition(partitions)

	// Set the tracker length to limit + 1 since we otherwise get
	// invalid formatting for the tracker
	tracker.Total = utility.Sum(limitByPartition) + 1

	matchedMessages := int64(0)

	startTime := time.Now()
	for {
		msg, err := consumer.ReadMessage(-1)

		if err != nil {
			tracker.MarkAsDone()
			utility.ExitOnError(err)
		}

		partitionId := msg.TopicPartition.Partition
		if counterByPartition[partitionId] < limitByPartition[partitionId] {
			message := parseMessage(msg, keyQuery, valueQuery)
			if message != nil {
				matchedMessages++
				tracker.Message = "Reading messages (" + strconv.FormatInt(matchedMessages, 10) + " matches)"
				messages.PushFront(message)
			}

			counterByPartition[partitionId] = counterByPartition[partitionId] + 1
			tracker.Increment(1)
		}

		if isLimitReached(limitByPartition, counterByPartition) {
			break
		}

	}
	stopTime := time.Now()
	elapsedTime := stopTime.Sub(startTime)

	// Sleep for the progress to catch up
	time.Sleep(100 * time.Millisecond)

	readMessage := utility.Sum(counterByPartition)
	tracker.MarkAsDone()
	return Result{
		Messages: *messages,
		MatchedMessages:    matchedMessages,
		ReadMessages: readMessage,
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

	_, err := fmt.Scanln()

	if err != nil {
		running = false
		tracker.MarkAsDone()
		utility.ExitOnError(err)
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

func createConsumer(bootstrap string, topic string, group string, offset string, tracker *progress.Tracker) (*kafka.Consumer) {
	if group == "" {
		group = "raccoon-" + strconv.Itoa(rand.Int())
	}
	consumer, consumerError := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrap,
		"group.id":           group,
		"auto.offset.reset":  offset,
		"enable.auto.commit": "false",
	})
	if consumerError != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(consumerError)
	}

	subscribeError := consumer.SubscribeTopics([]string{topic}, nil)

	if subscribeError != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(consumerError)
	}

	tracker.MarkAsDone()
	return consumer
}

func isLimitReached(limitByPartition map[int32]int64, counterByPartition map[int32]int64) bool {
	for id, limit := range limitByPartition {
		if counterByPartition[id] < limit {
			return false
		}
	}

	return true
}

func getMessageLimitByPartition(partitions *[]Partition, limit int64) map[int32]int64  {
	limits := make(map[int32]int64)

	for _, partition := range *partitions {
		max := partition.highOffset - partition.lowOffset
		if limit > max {
			limits[partition.id] = max
		} else {
			limits[partition.id] = limit
		}
	}

	return limits
}

func getCounterByPartition(partitions *[]Partition) map[int32]int64 {
	counters := make(map[int32]int64)

	for _, partition := range *partitions {
		counters[partition.id] = 0
	}

	return counters
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


