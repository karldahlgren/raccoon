# Raccoon: A CLI Kafka Search tool

[![Go Report Card](https://goreportcard.com/badge/github.com/karldahlgren/raccoon)](https://goreportcard.com/report/github.com/karldahlgren/raccoon)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/karldahlgren/raccoon/blob/master/LICENSE)

Raccoon is a command line search tool for Apache Kafka written in Go. 
The tool enables users to search and grep particular messages in a Kafka topic by providing a search query. 
Additionally, the tool also allows the user to tail a Kafka topic and filter messages based on a provided search query.
Once finished, all matched Kafka messages could either be presented in the terminal or exported to a CSV file.

<img src="https://raw.githubusercontent.com/karldahlgren/raccoon/main/img/screenshot.gif" width="90%"/>

Table Of Content
----

- [Installation](#installation)
- [Features](#features)
- [Running Raccoon](#running-raccoon)
    * [Grep](#grep)
    * [Tail](#tail)
- [Example](#example)
- [License](#license)

## Installation
Coming soon

## Features

- **Grep**: Search a Kafka topic and grep messages that matches a provided search query.
- **Tail**: Tail a Kafka topic and filter messages that matches a provided filter criteria.

## Running Raccoon

### Grep
The grep command will search through a Kafka topic from the earliest offset and read
messages from the beginning until the limit has been reached, or until the end of the topic has been reached.
All matched messages can be printed to the terminal and/or exported to a CSV file.

    Usage:
      raccoon grep [flags]
    
    Flags:
      -b, --bootstrap-server string   Bootstrap server address (Required)
      -g, --group string              Group name (Optional)
      -h, --help                      help for grep
      -k, --key-query string          Key query (Optional)
      -l, --limit int                 Limit message consumption per partition (Optional) (default 1000)
      -o, --output string             Output file name (Optional)
      -t, --topic string              Topic name (Required)
      -q, --value-query string        Value query (Optional)
      -v, --verbose                   Print output in terminal (Optional)
    
### Tail
The tail command will tail a Kafka topic from the latest offset and match all newly published 
messages on the subscribed topic with a provided filter query.
All matched messages can be printed to the terminal and/or exported to a CSV file.

    Usage:
      raccoon tail [flags]
    
    Flags:
      -b, --bootstrap-server string   Bootstrap server address (Required)
      -g, --group string              Group name (Optional)
      -h, --help                      help for tail
      -k, --key-query string          Key query (Optional)
      -l, --limit int                 Limit message consumption per partition. -1 is no limit (Optional) (default -1)
      -o, --output string             Output file name (Optional)
      -t, --topic string              Topic name (Required)
      -q, --value-query string        Value query (Optional)
      -v, --verbose                   Print output in terminal (Optional)

## Example

    raccoon grep -b localhost:9092 -q MyQuery -t MyTopic -o result.csv -l 1000000
    
    ____
    |  _ \ __ _  ___ ___ ___   ___  _ __
    | |_) / _` |/ __/ __/ _ \ / _ \| '_ \
    |  _ < (_| | (_| (_| (_) | (_) | | | |
    |_| \_\__,_|\___\___\___/ \___/|_| |_|
    Raccoon: Kafka search tool (v1.0.0)

    Connecting to Kafka                      ... done! [0 in 109ms]
    Reading topic partition metadata         ... done! [10 in 110ms]
    Reading messages (229 matches)           ... done! [999.01K in 3.042s]
    Disconnecting from Kafka                 ... done! [0 in 111ms]
    Writing to file                          ... done! [229 in 102ms]
    
    Summary:
      Read messages.......................:  999009
      Matched messages....................:  229
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
