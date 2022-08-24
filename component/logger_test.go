package component

import (
	"fmt"
	"os"
	"testing"
)

func TestXxx(t *testing.T) {
	logger := GetLogger()
	fmt.Println(logger)
}

func TestFoo(t *testing.T) {
	ids,_ := os.Getwd()
	fmt.Println(ids)
}
