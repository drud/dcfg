package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// expect takes a testing type and two values
// it raises an error if a does not equal b
func expect(t *testing.T, d string, a interface{}, b interface{}) {
	if a != b {
		t.Errorf(d)
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// refute takes a testing type and two values
// it raises and error if a equals b
func refute(t *testing.T, d string, a interface{}, b interface{}) {
	if a == b {
		t.Errorf(d)
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func RunCommand(command string, args []string) (string, error) {
	out, err := exec.Command(
		command,
		args...,
	).CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

func getTempFile(content string) string {
	contentBytes := []byte(content)
	tmpfile, err := ioutil.TempFile("", "dcfgtest")
	if err != nil {
		log.Fatal(err)
	}

	//defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(contentBytes); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

func TestRunOptions(t *testing.T) {
	cmdCFG := `
- name: create
  tasks:
  - action: command
    cmd: touch testdcfg.txt
- name: delete
  tasks:
  - action: command
    cmd: rm testdcfg.txt
    `

	tmpFile := getTempFile(cmdCFG)
	defer os.Remove(tmpFile)

	out, err := RunCommand("dcfg", []string{"run", "all", "--config", tmpFile})
	if err != nil {
		log.Fatal(err)
	}

	expect(t, "Should saying 'Running [all]'", strings.Contains(out, "Running [all]"), true)
	expect(t, "Should run the 'create' group", strings.Contains(out, "create"), true)
	expect(t, "Should run the 'delete' group", strings.Contains(out, "delete"), true)

	if _, err := os.Stat(tmpFile); os.IsExist(err) {
		t.Errorf("testdcfg.txt should no longer exist!")
	}

	out, err = RunCommand("dcfg", []string{"run", "create", "--config", tmpFile})
	if err != nil {
		log.Fatal(err)
	}

	expect(t, "Should saying 'Running [create]'", strings.Contains(out, "Running [create]"), true)
	expect(t, "Should run the 'create' group", strings.Contains(out, "create"), true)
	refute(t, "Should not run the 'delete' group", strings.Contains(out, "delete"), true)

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Errorf("testdcfg.txt should exist!")
	}

	out, err = RunCommand("dcfg", []string{"run", "create", "delete", "--config", tmpFile})
	if err != nil {
		log.Fatal(err)
	}

	expect(t, "Should saying 'Running [create delete]'", strings.Contains(out, "Running [create delete]"), true)
	expect(t, "Should run the 'create' group", strings.Contains(out, "create"), true)
	expect(t, "Should run the 'delete' group", strings.Contains(out, "delete"), true)

	if _, err := os.Stat(tmpFile); os.IsExist(err) {
		t.Errorf("testdcfg.txt should no longer exist!")
	}

}
