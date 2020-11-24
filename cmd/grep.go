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
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/karldahlgren/raccoon/kafka"
	"github.com/spf13/cobra"
)

var grepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Search through a Kafka topics that and grep matches",
	Long:  `Search through a Kafka topics that and grep matches`,
	Run: func(cmd *cobra.Command, args []string) {
		group := getStringFlag(cmd,"group")
		topic := getStringFlag(cmd,"topic")
		keyQuery := getStringFlag(cmd,"key-query")
		valueQuery := getStringFlag(cmd,"value-query")
		output := getStringFlag(cmd,"output")
		limit := getInt64Flag(cmd, "limit")
		verbose := getBoolFlag(cmd, "verbose")
		
		// Create progress and trackers
		writer := CreateProgress()
		createConsumerTracker := CreateTracker("Connecting to Kafka              # 1", 3, writer)
		consumeTracker := CreateTracker("Reading messages                 # 2", limit, writer)
		stopConsumerTracker := CreateTracker("Disconnecting from Kafka         # 3", 2, writer)
		var writeToFileTracker *progress.Tracker = nil
		if output != "" {
			writeToFileTracker = CreateTracker("Writing to file                  # 4", 1000000, writer)
		}

		InitiateProgress(writer)

		// Create Kafka consumer
		consumer := kafka.CreateConsumer(topic, group, createConsumerTracker)
		result := kafka.Consume(consumer, keyQuery, valueQuery, limit, consumeTracker)
		kafka.StopConsumer(consumer, stopConsumerTracker)

		if output != "" {
			writeResultToFile(result, output, writeToFileTracker)
		}

		printSummaryToPrompt(result)

		if verbose {
			printResultToPrompt(result)
		}

		FinishProgress(writer)
	},
}

func init() {
	grepCmd.Flags().StringP("broker", "b", "", "Broker address (Required)")
	grepCmd.Flags().StringP("topic", "t", "", "Topic name (Required)")
	grepCmd.Flags().StringP("group", "g", "", "Group name (Optional)")
	grepCmd.Flags().StringP( "value-query", "q", "", "Value query (Optional)")
	grepCmd.Flags().StringP( "key-query", "k", "", "Key query (Optional)")
	grepCmd.Flags().StringP("output", "o", "", "Output file name (Optional)")
	grepCmd.Flags().Int64P("limit", "l", 1000, "Limit message consumption per partition (Optional)")
	grepCmd.Flags().BoolP("verbose", "v", false, "Print output in terminal (Optional)")

	grepCmd.MarkFlagRequired("broker")
	grepCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(grepCmd)
}
