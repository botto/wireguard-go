/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2020 WireGuard LLC. All Rights Reserved.
 */

package device

import (
	"log"
	"os"
	"unsafe"
)

// A Logger provides logging for a Device.
// The functions are Printf-style functions.
// They must be safe for concurrent use.
// They do not require a trailing newline in the format.
// If nil, that level of logging will be silent.
type Logger struct {
	Verbosef func(format string, args ...interface{})
	Errorf   func(format string, args ...interface{})
}

// Log levels for use with NewLogger.
const (
	LogLevelSilent = iota
	LogLevelError
	LogLevelVerbose
)

// Function for use in Logger for discarding logged lines.
func DiscardLogf(format string, args ...interface{}) {}

// NewLogger constructs a Logger that writes to stdout.
// It logs at the specified log level and above.
// It decorates log lines with the log level, date, time, and prepend.
func NewLogger(level int, prepend string) *Logger {
	logger := &Logger{DiscardLogf, DiscardLogf}
	logf := func(prefix string) func(string, ...interface{}) {
		return log.New(os.Stdout, prefix+": "+prepend, log.Ldate|log.Ltime).Printf
	}
	if level >= LogLevelVerbose {
		logger.Verbosef = logf("DEBUG")
	}
	if level >= LogLevelError {
		logger.Errorf = logf("ERROR")
	}
	return logger
}

var discardLogfFptr uintptr

func init() {
	iface := (interface{})(DiscardLogf)
	discardLogfFptr = (*struct{ t, d uintptr })(unsafe.Pointer(&iface)).d
}

func isDiscardf(fn interface{}) bool {
	return (*struct{ t, d uintptr })(unsafe.Pointer(&fn)).d == discardLogfFptr
}

func (device *Device) verbosef(format string, args ...interface{}) {
	if !isDiscardf(device.log.Verbosef) {
		device.log.Verbosef(format, args...)
	}
}

func (device *Device) errorf(format string, args ...interface{}) {
	if !isDiscardf(device.log.Errorf) {
		device.log.Errorf(format, args...)
	}
}
