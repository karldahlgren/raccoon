# Raccoon: A CLI Kafka Search tool

[![Build Status](https://github.com/karldahlgren/raccoon/workflows/Go/badge.svg)](https://github.com/karldahlgren/raccoon/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/karldahlgren/raccoon)](https://goreportcard.com/report/github.com/karldahlgren/raccoon)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/karldahlgren/raccoon/blob/master/LICENSE)

Raccoon is a command line search tool for Apache Kafka. 
The tool allows the users to search for particular messages in a Kafka topic. 
The result could either be presented in the prompt or exported to a file.

<img src="https://github.com/karldahlgren/raccoon/blob/main/img/screenshot.gif" width="100%"/>

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
    Raccoon: Kafka search tool
    
    required flag(s) "broker", "group", "query", "topic" not set
    Usage:
      raccoon search [flags]
    
    Flags:
      -b, --broker string   Broker address (Required)
      -g, --group string    Group name (Optional)
      -h, --help            help for search
      -o, --output string   Output file name (Optional)
      -q, --query string    Query (Required)
      -t, --topic string    Topic name (Required)
      -v, --verbose         Print output in terminal (Optional)
    
    required flag(s) "broker", "group", "query", "topic" not set

##Example

    raccoon search -b localhost:9092 -q MyQuery -t MyTopic -o result.csv
    
    Establishing connection to Kafka   # 1   ... done! [3 in 115ms]
    Reading message                    # 2   ... done! [1.00M in 2.871s]
    Disconnecting from Kafka           # 3   ... done! [2 in 2.876s]
    Writing to file                    # 4   ... done! [280 in 2.878s]
    
    Matched messages: 280
    Search time: 3s

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
