package goUtils

import (
	"fmt"
	"testing"
)

func TestWalkDir(t *testing.T) {
	dirs := []string{"/Users/atlasv/GolandProjects/LHQ"}
	fmt.Println(WalkDir(dirs...))
}
