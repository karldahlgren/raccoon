# Raccoon: A CLI Kafka Search tool

[![Go Report Card](https://goreportcard.com/badge/github.com/karldahlgren/raccoon)](https://goreportcard.com/report/github.com/karldahlgren/raccoon)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/karldahlgren/raccoon/blob/master/LICENSE)

Raccoon is a command line search tool for Apache Kafka written in Go. 
The tool allows users to search for particular messages in a Kafka topic by providing a search query. 
The tool will check whether a message matches the provided query. 
Once finished, all matched Kafka messages could either be presented in the terminal or exported to a CSV file.

<img src="https://raw.githubusercontent.com/karldahlgren/raccoon/main/img/screenshot.gif" width="90%"/>

Table Of Content
----

- [Install](#install)
- [How to use](#how-to-use)
- [Example](#example)
- [License](#license)

## Install
Coming soon

## How to use
    ____
    |  _ \ __ _  ___ ___ ___   ___  _ __
    | |_) / _` |/ __/ __/ _ \ / _ \| '_ \
    |  _ < (_| | (_| (_| (_) | (_) | | | |
    |_| \_\__,_|\___\___\___/ \___/|_| |_|
    Raccoon: Kafka search tool (dev-build)
    
    Usage:
      raccoon grep [flags]
    
    Flags:
      -b, --broker string        Broker address (Required)
      -g, --group string         Group name (Optional)
      -h, --help                 help for grep
      -k, --key-query string     Key query (Optional)
      -l, --limit int            Limit message consumption per partition (Optional) (default 1000)
      -o, --output string        Output file name (Optional)
      -t, --topic string         Topic name (Required)
      -q, --value-query string   Value query (Optional)
      -v, --verbose              Print output in terminal (Optional)
    
    required flag(s) "broker", "topic" not set

A summary will always be printed in the terminal when the grep command has been executed. The summary will show the matched message count, and the search time. The full result can either be presented in the terminal by using the --verbose flag or be exported to a CSV file by using the --output flag. The result will include each the partitioned id, offset, timestamp, key and value.

## Example

    raccoon grep -b localhost:9092 -q MyQuery -t MyTopic -o result.csv
    
    Connecting to Kafka                      ... done! [2 in 105ms]
    Reading messages (280 matches)           ... done! [1.00M in 3.003s]
    Disconnecting from Kafka                 ... done! [1 in 100ms]
    Writing to file                          ... done! [280 in 104ms]
    
    Summary:
      Read messages.......................:  1000000
      Matched messages....................:  280
      Search time.........................:  3.000000s
      Messages/s..........................:  0.000003


## License

    The MIT License
    
    Copyright (c) 2020-, Karl A. Dahlgren
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.
