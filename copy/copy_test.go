package copy

import "testing"

func Test_CopyFile(t *testing.T) {
	if err := Copy("./copy.go", "./copybak"); err != nil {
		t.Error("copy error")
	} else {
		t.Log("copy pass")
	}
}

func Test_CopyDir(t *testing.T) {
	if err := Copy("../cmd", "../cmde"); err != nil {
		t.Error("copy dir error")
	} else {
		t.Log("copy dir pass")
	}
}
