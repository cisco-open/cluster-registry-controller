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
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"

	"github.com/banzaicloud/operator-tools/pkg/resources"
	"github.com/banzaicloud/operator-tools/pkg/utils"
)

func K8SResourceOverlayPatchExecuteTemplate(patch resources.K8SResourceOverlayPatch, data interface{}) (resources.K8SResourceOverlayPatch, error) {
	p := patch.DeepCopy()

	value := utils.PointerToString(p.Value)
	t, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(value)
	if err != nil {
		return resources.K8SResourceOverlayPatch{}, err
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		return resources.K8SResourceOverlayPatch{}, err
	}

	p.Value = utils.StringPointer(tpl.String())

	return *p, nil
}

func K8SResourceOverlayPatchExecuteTemplates(patches []resources.K8SResourceOverlayPatch, data interface{}) ([]resources.K8SResourceOverlayPatch, error) {
	result := make([]resources.K8SResourceOverlayPatch, 0)
	for _, p := range patches {
		p, err := K8SResourceOverlayPatchExecuteTemplate(p, data)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}
