// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package errors

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type IndexOutOfBoundsError struct {
	index  int
	length int
}

func (e IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("Index %d out of bounds for length %d", e.index, e.length)
}

func (e IndexOutOfBoundsError) GetIndex() int {
	return e.index
}

func (e IndexOutOfBoundsError) GetLength() int {
	return e.length
}

func NewIndexOutOfBoundsError(index, length int) *IndexOutOfBoundsError {
	return &IndexOutOfBoundsError{
		index:  index,
		length: length,
	}
}

type IllegalArgumentError struct {
	argument any
}

func (e IllegalArgumentError) Error() string {
	return fmt.Sprintf("Illegal argument: %v", e.argument)
}

func (e IllegalArgumentError) GetArgument() any {
	return e.argument
}

func NewIllegalArgumentError(argument any) *IllegalArgumentError {
	return &IllegalArgumentError{
		argument: argument,
	}
}

type ErrorWithStackTrace struct {
	err   error
	stack []uintptr
}

func (e *ErrorWithStackTrace) Error() string {
	if e.stack != nil {
		return fmt.Sprintf("%s%s", e.err.Error(), e.StackTrace())
	}
	return e.err.Error()
}

func DecorateWithStackTrace(err error) error {
	if os.Getenv("BAL_BACKTRACE") != "true" {
		return err
	}

	if err == nil {
		return nil
	}

	const (
		maxDepth = 32
		skip     = 2 // skip DecorateWithStackTrace and its caller
	)

	stack := make([]uintptr, maxDepth)
	length := runtime.Callers(skip, stack[:])
	return &ErrorWithStackTrace{
		err:   err,
		stack: stack[:length],
	}
}

func (e *ErrorWithStackTrace) StackTrace() []byte {
	if e == nil || len(e.stack) == 0 {
		return nil
	}

	var buf bytes.Buffer
	frames := runtime.CallersFrames(e.stack)
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&buf, "\n\tat %s(%s:%d)", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	return buf.Bytes()
}

const (
	logSource            = "ballerina"
	logLevel             = "SEVERE"
	internalErrorMessage = `ballerina: Oh no, something really went wrong. Bad. Sad.

We appreciate it if you can report the code that broke Ballerina in
https://github.com/ballerina-platform/ballerina-lang/issues with the
log you get below and your sample code.

We thank you for helping make us better.
`
)

// LogBadSad logs unhandled errors with an internal error message.
// These are unexpected errors in the runtime that should be reported.
func LogBadSad(err error) {
	fmt.Fprint(os.Stderr, internalErrorMessage)
	PrintCrashLog(err)
}

// PrintCrashLog logs error messages to stderr in a structured format.
// Format: [timestamp] LEVEL {source} - error message
func PrintCrashLog(err error) {
	if err == nil {
		return
	}

	now := time.Now()
	timestamp := now.Format("2006-01-02 15:04:05") + fmt.Sprintf(",%03d", now.Nanosecond()/1e6)

	msg := fmt.Sprintf("[%s] %-5s {%s} - %s",
		timestamp,
		logLevel,
		logSource,
		err.Error(),
	)

	logger := log.New(os.Stderr, "", 0)
	logger.Println(msg)
}
