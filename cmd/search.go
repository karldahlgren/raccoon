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

var broker, group, topic, query, output string
var verbose bool

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search through a Kafka topics that matches",
	Long:  `Search through a Kafka topics that matches`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create progress and trackers
		writer := CreateProgress()
		createConsumerTracker := CreateTracker("Establishing connection to Kafka   # 1", 3, writer)
		consumeTracker := CreateTracker("Reading message                    # 2", 1000000, writer)
		stopConsumerTracker := CreateTracker("Disconnecting from Kafka           # 3", 2, writer)
		var writeToFileTracker *progress.Tracker = nil
		if output != "" {
			writeToFileTracker = CreateTracker("Writing to file                    # 4", 1000000, writer)
		}

		InitiateProgress(writer)

		// Create Kafka consumer
		consumer := kafka.CreateConsumer(topic, group, createConsumerTracker)
		result := kafka.Consume(consumer, consumeTracker)
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
	searchCmd.PersistentFlags().StringVarP(&broker, "broker", "b", "", "Broker address (Required)")
	searchCmd.PersistentFlags().StringVarP(&topic, "topic", "t", "", "Topic name (Required)")
	searchCmd.PersistentFlags().StringVarP(&group, "group", "g", "", "Group name (Optional)")
	searchCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Query (Required)")
	searchCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output file name (Optional)")
	searchCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,"Print output in terminal (Optional)")

	searchCmd.MarkPersistentFlagRequired("broker")
	searchCmd.MarkPersistentFlagRequired("topic")
	searchCmd.MarkPersistentFlagRequired("query")
	rootCmd.AddCommand(searchCmd)
}
