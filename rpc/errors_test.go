// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewError(t *testing.T) {
	tests := []struct {
		Code    int
		Message string
	}{
		{
			110,
			"oneonezero",
		},
		{
			111,
			"oneoneone",
		},
	}

	for _, test := range tests {
		got := NewError(test.Code, test.Message)
		assert.Equal(t, got.Code, test.Code)
		assert.Equal(t, got.Message, test.Message)
	}
}

func Test_newError(t *testing.T) {
	tests := []struct {
		Message string
	}{
		{
			"rpc: service/method request ill-formed",
		},
		{
			"rpc: can't find service",
		},
		{
			"rpc: can't find method",
		},
	}

	for _, test := range tests {
		got := newError(test.Message)
		assert.Equal(t, got.Message, test.Message)
	}

	got1 := newError("sdsd")
	assert.Equal(t, got1.Message, "sdsd")
}

func Test_ServerError(t *testing.T) {
	err := ServerError(nil)
	assert.Equal(t, err == nil, true)
	rpcerr := errors.New(`{"code":-32603,"message":"Internal error","Data":"[]int{1, 2, 3}"}`)
	err = ServerError(rpcerr)
	assert.Equal(t, err.Data != nil, true)
	_, ok := err.Data.(*Error)
	assert.Equal(t, ok, false)
	assert.Equal(t, err.Message, "Internal error")
	assert.Equal(t, err.Code, -32603)
}

func Test_Error(t *testing.T) {
	e := &Error{
		110,
		"oneonezero",
		[]string{"Math", "English", "Chinese"},
	}

	nameOut1 := e.Error()
	buf, err := json.Marshal(e)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(buf), nameOut1)
}
