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

package cmd

import (
	"fmt"
	"github.com/karldahlgren/raccoon/anothername"
	"github.com/spf13/cobra"
)

var grepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Search through a Kafka topics and grep all matching messages",
	Long:  `The tail command will subscribe to a topic with the earliest offset and read the 
			specified amount of messages, or to the end of the topic, from each partition. 
			All matched messages can be displayed in the terminal and/or exported to a CSV file.`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap := getStringFlag(cmd,"bootstrap-server")
		group := getStringFlag(cmd,"group")
		topic := getStringFlag(cmd,"topic")
		keyQuery := getStringFlag(cmd,"key-query")
		valueQuery := getStringFlag(cmd,"value-query")
		output := getStringFlag(cmd,"output")
		seekTimestamp := getStringFlag(cmd,"seek")
		limit := getInt64Flag(cmd, "limit")
		verbose := getBoolFlag(cmd, "verbose")
		earliest := getBoolFlag(cmd, "earliest")
		latest := getBoolFlag(cmd, "latest")

		if earliest == true && latest == true {
			fmt.Printf("Not allowed to combine earliest flag with latest flag")
			return
		} else if seekTimestamp != "" && earliest == true {
			fmt.Printf("Not allowed to combine seek timestamp flag with earliest flag")
			return
		} else if seekTimestamp != "" && latest == true {
			fmt.Printf("Not allowed to combine seek timestamp flag with latest flag")
			return
		} else if limit < 0 {
			fmt.Printf("Limit cannot be less than zero")
			return
		}

		// Create progress and trackers
		writer := CreateProgress()
		InitiateProgress(writer)

		// Create Kafka consumer
		createConsumerTracker := CreateTracker("Connecting to Kafka", 2, writer)
		consumer := anothername.CreateEarliestConsumer(bootstrap, topic, group, createConsumerTracker)

		// Retrieve Partition metadata
		getPartitionsTracker := CreateTracker("Reading topic partition metadata", 100, writer)
		partitions := anothername.GetPartitions(consumer, topic, getPartitionsTracker)

		// Consumer from Kafka topic
		consumeTracker := CreateTracker("Reading messages (0 matches)", limit, writer)
		result := anothername.Consume(consumer, partitions, topic, keyQuery, valueQuery, limit, seekTimestamp, latest, consumeTracker)

		// Stop Kafka consumer
		stopConsumerTracker := CreateTracker("Disconnecting from Kafka", 1, writer)
		anothername.StopConsumer(consumer, stopConsumerTracker)

		if output != "" {
			// File output has been provided. Writing to file
			writeToFileTracker := CreateTracker("Writing to file", limit, writer)
			writeResultToFile(result, output, writeToFileTracker)
		}

		printSummaryToPrompt(result)

		if verbose {
			// Print all the matched messages in the terminal
			printResultToPrompt(result)
		}

		FinishProgress(writer)
	},
}

func init() {
	grepCmd.Flags().StringP("bootstrap-server", "b", "", "Bootstrap server address (Required)")
	grepCmd.Flags().StringP("topic", "t", "", "Topic name (Required)")
	grepCmd.Flags().StringP("group", "g", "", "Group name (Optional)")
	grepCmd.Flags().StringP( "value-query", "q", "", "Value query (Optional)")
	grepCmd.Flags().StringP( "key-query", "k", "", "Key query (Optional)")
	grepCmd.Flags().StringP("output", "o", "", "Output file name (Optional)")
	grepCmd.Flags().String("seek", "", "Seek and set offset to a timestamp. RFC3339 time format (Optional)")
	grepCmd.Flags().Int64P("limit", "l", 1000, "Limit message consumption per partition (Optional)")
	grepCmd.Flags().BoolP("verbose", "v", false, "Print output in terminal (Optional)")
	grepCmd.Flags().Bool("earliest", false, "Start at the earliest offset (Optional)")
	grepCmd.Flags().Bool("latest", false, "Start at the latest offset minus the limit (Optional)")


	_ = grepCmd.MarkFlagRequired("bootstrap-server")
	_ = grepCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(grepCmd)
}
