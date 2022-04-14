// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
