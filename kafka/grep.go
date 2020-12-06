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
	confluent "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"strconv"
	"time"
)

// Consume messages from a Kafka consumer
func Consume(consumer *confluent.Consumer, partitions map[int32]Partition, topic string, keyQuery string, valueQuery string,
	limit int64, seekTimestamp string, latest bool, tracker *progress.Tracker) Result {
	if seekTimestamp != "" {
		partitions = seekToTimestamp(consumer, partitions, topic, seekTimestamp)
	} else if latest {
		partitions = seekToLatest(consumer, partitions, topic, limit)
	}

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
		if isLimitReached(limitByPartition, counterByPartition) {
			break
		}

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


func isLimitReached(limitByPartition map[int32]int64, counterByPartition map[int32]int64) bool {
	for id, limit := range limitByPartition {
		if counterByPartition[id] < limit {
			return false
		}
	}

	return true
}

func getMessageLimitByPartition(partitions map[int32]Partition, limit int64) map[int32]int64  {
	limits := make(map[int32]int64)

	for _, partition := range partitions {
		max := partition.highOffset - partition.lowOffset
		if limit > max {
			limits[partition.id] = max
		} else {
			limits[partition.id] = limit
		}
	}

	return limits
}

func getCounterByPartition(partitions map[int32]Partition) map[int32]int64 {
	counters := make(map[int32]int64)

	for _, partition := range partitions {
		counters[partition.id] = 0
	}

	return counters
}
