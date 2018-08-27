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
	for i, test := range tests {
		got := NewError(test.Code, test.Message)
		if got.Code != test.Code && got.Message != test.Message {
			t.Errorf(" case #%d: Code %v; Message: %v", i+1, test.Code, test.Message)
		}
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

	rpcerr := errors.New(`{"code":12,"message":"2323"}`)
	err = ServerError(rpcerr)
	assert.Equal(t, err.Message, "2323")
	assert.Equal(t, err.Code, 12)
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
