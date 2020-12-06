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

package anothername

import (
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/karldahlgren/raccoon/utility"
	confluent "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"math/rand"
	"strconv"
	"time"
)

// CreateEarliestConsumer Creates a new Kafka consumer with the earliest offset
func CreateEarliestConsumer(bootstrap string, topic string, group string, tracker *progress.Tracker) *confluent.Consumer {
	return createConsumer(bootstrap, topic, group, "earliest", tracker)
}

// CreateLatestConsumer Creates a new Kafka consumer with the latest offset
func CreateLatestConsumer(bootstrap string, topic string, group string, tracker *progress.Tracker) *confluent.Consumer {
	return createConsumer(bootstrap, topic, group, "latest", tracker)
}

// StopConsumer will stop and disconnect a consumer from Kafka
func StopConsumer(consumer *confluent.Consumer, tracker *progress.Tracker) {
	err := consumer.Close()

	if err != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}
	tracker.MarkAsDone()
}

// GetPartitions retrieves information regarding all partitions for a provided topic
func GetPartitions(consumer *confluent.Consumer, topic string, tracker *progress.Tracker) map[int32]Partition {
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

	// Set the tracker length to limit + 1 since we otherwise get
	// invalid formatting for the tracker
	tracker.Total = int64(len(topicMetaData.Partitions)) + 1
	partitions := make(map[int32]Partition)
	for index, partition := range topicMetaData.Partitions {
		lowOffset, highOffset, err := consumer.QueryWatermarkOffsets(topic, partition.ID, -1)

		if err != nil {
			tracker.MarkAsDone()
			utility.ExitOnError(err)
		}

		partitions[partition.ID] = Partition {
			id: partition.ID,
			lowOffset: lowOffset,
			highOffset: highOffset,
		}
		tracker.Message = "Reading topic partition metadata (" + strconv.Itoa(index + 1) + " partitions)"
		tracker.Increment(1)
	}
	// Sleep for the progress to catch up
	time.Sleep(100 * time.Millisecond)
	tracker.MarkAsDone()
	return partitions
}

func createConsumer(bootstrap string, topic string, group string, offset string, tracker *progress.Tracker) (*confluent.Consumer) {
	if group == "" {
		group = "raccoon-" + strconv.Itoa(rand.Int())
	}
	consumer, consumerError := confluent.NewConsumer(&confluent.ConfigMap{
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


