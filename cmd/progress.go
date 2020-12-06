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
	"github.com/jedib0t/go-pretty/v6/text"
	"time"
)

// CreateProgress creates a new progress, which is used to create progress bars and trackers
func CreateProgress() progress.Writer {
	pw := progress.NewWriter()
	pw.SetTrackerLength(30)
	pw.ShowOverallTracker(false)
	pw.ShowTime(true)
	pw.ShowTracker(true)
	pw.ShowValue(true)
	pw.SetMessageWidth(55)
	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleCircle)
	pw.Style().Colors.Tracker = text.Colors{text.FgYellow}
	pw.Style().Colors.Percent = text.Colors{text.FgRed}
	pw.Style().Colors.Stats = text.Colors{text.FgBlue}
	pw.Style().Colors.Time = text.Colors{text.FgGreen}
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.SetNumTrackersExpected(3)
	return pw
}

// CreateTracker Create new tracker for a progress
func CreateTracker(message string, total int64, writer progress.Writer) *progress.Tracker {
	tracker := progress.Tracker{Message: message, Total: total, Units: progress.UnitsDefault}
	writer.AppendTracker(&tracker)
	time.Sleep(100 * time.Millisecond)
	return &tracker
}

// InitiateProgress will render the progress bars and sleep 100 ms for rendering reasons
func InitiateProgress(writer progress.Writer) {
	go writer.Render()
	time.Sleep(100 * time.Millisecond)
}

// FinishProgress finish the progress bars
func FinishProgress(writer progress.Writer) {
	time.Sleep(time.Second)
	if writer.LengthActive() == 0 {
		writer.Stop()
	}
}
