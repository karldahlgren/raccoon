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
	"github.com/karldahlgren/raccoon/utility"
	"github.com/spf13/cobra"
)

func getStringFlag(cmd *cobra.Command, name string) string  {
	value, err := cmd.Flags().GetString(name)

	if err != nil {
		utility.ExitOnError(err)
	}

	return value
}

func getInt64Flag(cmd *cobra.Command, name string) int64  {
	value, err := cmd.Flags().GetInt64(name)

	if err != nil {
		utility.ExitOnError(err)
	}

	return value
}

func getBoolFlag(cmd *cobra.Command, name string) bool  {
	value, err := cmd.Flags().GetBool(name)

	if err != nil {
		utility.ExitOnError(err)
	}

	return value
}
