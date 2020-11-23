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

package config

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const versionURL = "https://github.com/karldahlgren/raccoon/releases/latest"
const updateURL = "https://github.com/karldahlgren/raccoon#install"

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func CheckForUpdates() {
	currentVersion := Version
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100)

	if r == 0 {
		req, err := http.NewRequest("HEAD", versionURL, nil)
		if err != nil {
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		loc, err := resp.Location()
		if err != nil {
			return
		}

		s := strings.Split(loc.Path, "/")
		newVersion := s[len(s)-1]

		if currentVersion != newVersion {
			fmt.Printf("New version avaiable %s, check %s\n\n", newVersion, updateURL)
		}
	}
}

