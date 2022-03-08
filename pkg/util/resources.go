// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package util

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GVKToString(gvk schema.GroupVersionKind) string {
	return fmt.Sprintf("%s.%s/%s", gvk.Kind, gvk.Group, gvk.Version)
}

func ParseGVKFromString(str string) *schema.GroupVersionKind {
	fp := strings.SplitN(str, ".", 2)
	if len(fp) != 2 {
		return nil
	}

	sp := strings.SplitN(fp[1], "/", 2)
	if len(sp) != 2 {
		return nil
	}

	return &schema.GroupVersionKind{
		Group:   sp[0],
		Kind:    fp[0],
		Version: sp[1],
	}
}
