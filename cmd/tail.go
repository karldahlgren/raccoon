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
	"github.com/karldahlgren/raccoon/kafka"
	"github.com/spf13/cobra"
)

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail a Kafka topic and match all new messages",
	Long:  `The tail command will subscribe to a topic with the latest offset and listen for all new message
			published on the topic. All matched messages can be displayed in the terminal and/or exported to a CSV file.`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap := getStringFlag(cmd,"bootstrap-server")
		group := getStringFlag(cmd,"group")
		topic := getStringFlag(cmd,"topic")
		keyQuery := getStringFlag(cmd,"key-query")
		valueQuery := getStringFlag(cmd,"value-query")
		output := getStringFlag(cmd,"output")
		limit := getInt64Flag(cmd, "limit")
		verbose := getBoolFlag(cmd, "verbose")

		fmt.Println("Press enter to stop reading messages")
		fmt.Println()
		
		// Create progress and trackers
		writer := CreateProgress()

		InitiateProgress(writer)

		// Create Kafka consumer
		createConsumerTracker := CreateTracker("Connecting to Kafka", 2, writer)
		consumer := kafka.CreateLatestConsumer(bootstrap, topic, group, createConsumerTracker)

		// Consumer from Kafka topic
		consumeTracker := CreateTracker("Reading messages (0 matches)", limit, writer)
		result := kafka.Tail(consumer, keyQuery, valueQuery, limit, consumeTracker)

		// Stop Kafka consumer
		stopConsumerTracker := CreateTracker("Disconnecting from Kafka", 1, writer)
		kafka.StopConsumer(consumer, stopConsumerTracker)

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
	tailCmd.Flags().StringP("bootstrap-server", "b", "", "Bootstrap server address (Required)")
	tailCmd.Flags().StringP("topic", "t", "", "Topic name (Required)")
	tailCmd.Flags().StringP("group", "g", "", "Group name (Optional)")
	tailCmd.Flags().StringP( "value-query", "q", "", "Value query (Optional)")
	tailCmd.Flags().StringP( "key-query", "k", "", "Key query (Optional)")
	tailCmd.Flags().StringP("output", "o", "", "Output file name (Optional)")
	tailCmd.Flags().Int64P("limit", "l", -1, "Limit message consumption per partition. -1 is no limit (Optional)")
	tailCmd.Flags().BoolP("verbose", "v", false, "Print output in terminal (Optional)")

	_ = tailCmd.MarkFlagRequired("bootstrap-server")
	_ = tailCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(tailCmd)
}
