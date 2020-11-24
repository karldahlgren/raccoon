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
	"encoding/csv"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/karldahlgren/raccoon/kafka"
	"github.com/karldahlgren/raccoon/utility"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"time"
)

func writeResultToFile(result kafka.Result, output string, tracker *progress.Tracker) {
	if result.Messages.Len() == 0 {
		tracker.MarkAsDone()
		return
	}

	tracker.Total = int64(result.Messages.Len())
	file, err := os.Create(output)
	if err != nil {
		tracker.MarkAsDone()
		utility.ExitOnError(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writeHeader(*writer)
	for element := result.Messages.Front(); element != nil; element = element.Next() {
		message := element.Value.(kafka.Message)
		row := getData(message)
		writeRow(*writer, row)
		tracker.Increment(1)
	}
	tracker.MarkAsDone()
}

func writeHeader(writer csv.Writer) {
	header := []string{"partition", "offset", "timestamp", "keu", "value"}
	err := writer.Write(header)
	if err != nil {
		utility.ExitOnError(err)
	}
}

func writeRow(writer csv.Writer, row []string) {
	err := writer.Write(row)
	if err != nil {
		utility.ExitOnError(err)
	}
}

func printSummaryToPrompt(result kafka.Result) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
	fmt.Println("Matched messages: " + strconv.FormatInt(result.Count, 10))
	fmt.Println("Search time: " + result.Duration.Round(time.Second).String())
	fmt.Println()
}

func printResultToPrompt(result kafka.Result) {
	if result.Messages.Len() == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Partition", "Offset", "Timestamp", "Key", "Value"})

	for element := result.Messages.Front(); element != nil; element = element.Next() {
		message := element.Value.(kafka.Message)
		row := getData(message)
		table.Append(row)
	}
	table.Render()
}

func getData(message kafka.Message) []string {
	return []string{
		strconv.FormatInt(int64(message.Partition), 10),
		message.Offset,
		message.Timestamp.String(),
		message.Key,
		message.Value}
}
