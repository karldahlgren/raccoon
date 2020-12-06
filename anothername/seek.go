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
	"github.com/karldahlgren/raccoon/utility"
	confluent "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"strconv"
	"time"
)

func seekToTimestamp(consumer *confluent.Consumer, partitions map[int32]Partition, topic string, seekTimestamp string) map[int32]Partition  {
	timestamp, err := time.Parse(time.RFC3339, seekTimestamp)

	if err != nil {
		utility.ExitOnError(err)
	}

	var topicPartitions []confluent.TopicPartition
	for _, partition := range partitions {
		topicPartition := confluent.TopicPartition{
			Topic: &topic,
			Partition: partition.id,
			Offset: confluent.Offset(timestamp.Round(time.Millisecond).Unix() * 1000),
		}

		topicPartitions = append(topicPartitions, topicPartition)
	}

	offsetTopicPartitions, err := consumer.OffsetsForTimes(topicPartitions, -1)
	var newOffsetPartitions []confluent.TopicPartition
	for _, partition := range offsetTopicPartitions {
		newOffset := partition.Offset
		if partition.Offset == -1 {
			newOffset = confluent.Offset(partitions[partition.Partition].highOffset)
			partition.Offset = newOffset
		}
		newOffsetPartition := confluent.TopicPartition{
			Topic: &topic,
			Partition: partition.Partition,
			Offset: newOffset,
		}

		newOffsetPartitions = append(newOffsetPartitions, newOffsetPartition)
	}

	return seek(consumer, newOffsetPartitions, partitions)
}

func seekToLatest(consumer *confluent.Consumer, partitions map[int32]Partition, topic string, limit int64) map[int32]Partition {
	var newOffsetPartitions []confluent.TopicPartition
	for _, partition := range partitions {
		offset := partition.highOffset - limit

		if offset < partition.lowOffset {
			offset = partition.lowOffset
		}

		topicPartition := confluent.TopicPartition{
			Topic: &topic,
			Partition: partition.id,
			Offset: confluent.Offset(offset),
		}

		newOffsetPartitions = append(newOffsetPartitions, topicPartition)
	}

	return seek(consumer, newOffsetPartitions, partitions)
}

func seek(consumer *confluent.Consumer, topicPartitions []confluent.TopicPartition,
	partitions map[int32]Partition) map[int32]Partition  {

	err := consumer.Assign(topicPartitions)
	if err != nil {
		utility.ExitOnError(err)
	}

	_,err = consumer.CommitOffsets(topicPartitions)
	if err != nil {
		utility.ExitOnError(err)
	}

	updatedPartitions := make(map[int32]Partition)
	for _,newOffsetPartition := range topicPartitions {
		lowOffset, err := strconv.ParseInt(newOffsetPartition.Offset.String(), 10, 64)

		if err != nil {
			utility.ExitOnError(err)
		}

		updatedPartitions[newOffsetPartition.Partition] = Partition {
			id: newOffsetPartition.Partition,
			lowOffset: lowOffset,
			highOffset: partitions[newOffsetPartition.Partition].highOffset,
		}
	}

	return updatedPartitions
}