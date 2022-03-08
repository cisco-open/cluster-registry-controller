// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clustermeta_test

import (
	"io/ioutil"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var scheme = runtime.NewScheme()

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
}

// ReadFile reads the content of the given file or fails the test if an error is encountered.
func ReadFile(tb testing.TB, file string) []byte {
	tb.Helper()
	content, err := ioutil.ReadFile(file)
	if err != nil {
		tb.Fatalf(err.Error())
	}

	return content
}
