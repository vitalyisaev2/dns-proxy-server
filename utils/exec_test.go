package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExecSuccess(t *testing.T) {

    out, err, code := Exec("echo", "hi")

    assert.Equal(t, "hi\n", string(out))
    assert.Equal(t, nil, err)
    assert.Equal(t, 0, code)

}

func TestExecExitSuccess(t *testing.T) {

    out, err, code := Exec("mkdir", "/")

    assert.Equal(t, []byte("mkdir: cannot create directory '/': File exists\n"), out)
    assert.Equal(t, "exit status 1", err.Error())
    assert.Equal(t, 1, code)

}

func TestExecCommandThatNotExistsSuccess(t *testing.T) {

    out, err, code := Exec("notExists")

    assert.Equal(t, "", string(out))
    assert.Equal(t, "exec: \"notExists\": executable file not found in $PATH", err.Error())
    assert.Equal(t, -255, code)

}

func TestExistsSuccess(t *testing.T) {
    exists := Exists("which")
    assert.Equal(t, true, exists)
}

func TestExistsFalse(t *testing.T) {
    exists := Exists("notExists")
    assert.Equal(t, false, exists)
}