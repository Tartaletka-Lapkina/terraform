// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
)

func TestModules_noJsonFlag(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	ui := new(cli.MockUi)
	view, _ := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{}
	code := cmd.Run(args)
	if code == 0 {
		t.Fatal("expected an non zero exit status\n")
	}

	output := ui.ErrorWriter.String()
	if !strings.Contains(output, "The `terraform modules` command requires the `-json` flag.\n") {
		t.Fatal("expected an error message about requiring -json flag.\n")
	}

	if !strings.Contains(output, modulesCommandHelp) {
		t.Fatal("expected the modules command help to be displayed\n")
	}
}

func TestModules_fullCmd(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	testCopyDir(t, testFixturePath("modules"), dir)

	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{"-json"}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("Got a non-zero exit code: %d\n", code)
	}

	output := done(t).All()
	compareJSONOutput(t, output, expected_fullCmd)
}

func TestModules_fullCmd_unreferencedEntries(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	testCopyDir(t, testFixturePath("modules-unreferenced-entries"), dir)

	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{"-json"}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("Got a non-zero exit code: %d\n", code)
	}
	output := done(t).All()
	compareJSONOutput(t, output, expected_fullCmd_unreferencedEntries)
}

func TestModules_uninstalledModules(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	testCopyDir(t, testFixturePath("modules-uninstalled-entries"), dir)

	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{"-json"}
	code := cmd.Run(args)
	if code == 0 {
		t.Fatal("Expected a non-zero exit code\n")
	}
	output := done(t).All()
	if !strings.Contains(output, "Module not installed") {
		t.Fatalf("expected to see a `not installed` error message: %s\n", output)
	}

	if !strings.Contains(output, `Run "terraform init"`) {
		t.Fatalf("expected error message to ask user to run terraform init: %s\n", output)
	}
}

func compareJSONOutput(t *testing.T, got string, want string) {
	var expected, actual map[string]interface{}

	if err := json.Unmarshal([]byte(got), &actual); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("unexpected output, got: %s\n, want:%s\n", got, want)
	}
}

var expected_fullCmd = `{"Modules":[{"Source":"./child","Version":"","Key":"child","Dir":"child","Referenced":true},{"Source":"./child","Version":"","Key":"count_child","Dir":"child","Referenced":true}]}`

var expected_fullCmd_unreferencedEntries = `{"Modules":[{"Source":"./child","Version":"","Key":"child","Dir":"child","Referenced":true},{"Source":"./child","Version":"","Key":"count_child","Dir":"child","Referenced":true},{"Source":"./child","Version":"","Key":"old_count_child","Dir":"child","Referenced":false}]}`
