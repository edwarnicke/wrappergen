// Copyright (c) 2022 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wrappergen

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Template struct {
	templates map[string]*template.Template
	input     fs.FS
}

func ParseFS(input fs.FS, patterns ...string) (*Template, error) {
	rv := &Template{
		input:     input,
		templates: make(map[string]*template.Template),
	}
	var err error
	fs.WalkDir(rv.input, ".", rv.addAllToTemplateWalkFn)
	if err != nil {
		return nil, err
	}
	return rv, nil
}

func (t *Template) addAllToTemplateWalkFn(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
		return nil
	}
	if !strings.HasSuffix(d.Name(), ".tmpl") {
		return nil
	}
	tmpl, err := template.ParseFS(t.input, path)
	if err != nil {
		return err
	}
	t.templates[path] = tmpl
	return nil
}

func (t *Template) createExecuteWalkFn(outputDir string, data interface{}) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".tmpl") {
			return nil
		}

		outputBuffer := bytes.NewBuffer([]byte{})
		tmpl, ok := t.templates[path]
		if !ok {
			return nil
		}
		if err := tmpl.Execute(outputBuffer, data); err != nil {
			return err
		}

		if strings.TrimSpace(outputBuffer.String()) == "" {
			return nil
		}

		outputPath := filepath.Join(outputDir, path)

		if err := os.MkdirAll(filepath.Dir(outputPath), 0700); err != nil {
			return err
		}

		output, err := os.Create(strings.TrimSuffix(outputPath, ".tmpl"))
		defer func() { _ = output.Close() }()

		_, err = io.Copy(output, outputBuffer)
		return err
	}
}

func (t *Template) ExecuteAll(outputDir string, data interface{}) error {
	if err := fs.WalkDir(t.input, ".", t.createExecuteWalkFn(outputDir, data)); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	return nil
}
