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
	"github.com/karldahlgren/raccoon/config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "raccoon",
	Short: "Raccoon is a Kafka-tool to search and find messages in Kafka topics",
}

// Execute method prints the logo and starts all commands
func Execute() {
	fmt.Println("____")
	fmt.Println("|  _ \\ __ _  ___ ___ ___   ___  _ __")
	fmt.Println("| |_) / _` |/ __/ __/ _ \\ / _ \\| '_ \\")
	fmt.Println("|  _ < (_| | (_| (_| (_) | (_) | | | |")
	fmt.Println("|_| \\_\\__,_|\\___\\___\\___/ \\___/|_| |_|")
	fmt.Println("Raccoon: Kafka search tool " + config.Version)
	fmt.Println()
	config.CheckForUpdates()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
