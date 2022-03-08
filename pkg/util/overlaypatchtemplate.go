// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

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
