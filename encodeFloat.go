package graphql

import (
	"bytes"
	"math"
	"strconv"
)

// Modified copy of https://golang.org/src/encoding/json/encode.go > floatEncoder.encode(..)
// Copyright for function below:
//
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// IMPORTANT the full license can be found in this repo: https://github.com/golang/go
func floatToJson(bits int, f float64, e *bytes.Buffer) {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		e.WriteString("0.0")
		return
	}

	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}

	b := strconv.FormatFloat(f, fmt, -1, bits)
	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(b)
		if n >= 4 && b[n-4:n-1] == "e-0" {
			e.WriteString(b[:n-2] + b[n-1:])
			return
		}
	}
	e.WriteString(b)
}