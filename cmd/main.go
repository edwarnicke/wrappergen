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

package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
)

//go:embed _templates/main.go.impl
var mainGo string

// TemplateInput - Configuration for imports-gen
type TemplateInput struct {
	// Package - package for generated imports.go
	Package string
}

func main() {

	filename := "main.go"
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		log.Fatalf("unable to remove %s because %+v", filename, err)
	}

	// Load input information from environment
	// go generate passes a number of envs we will pick up this way
	input := &TemplateInput{}
	input.Package = os.Getenv("GOPACKAGE")

	if input.Package == "" {
		log.Fatal("error did not find GOPACKAGE env")
	}

	// Create the template
	tmpl := template.Must(template.New(filename).Parse(mainGo))

	// Create the main.go file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error creating file: %q: %+v", filename, err)
	}
	defer func() { _ = f.Close() }()
	if err := tmpl.Execute(f, input); err != nil {
		fmt.Fprintf(os.Stderr, "error processing template: %+v", err)
	}
}
