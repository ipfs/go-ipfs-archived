// +build testrunmain

package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestRunMain(t *testing.T) {
	args := flag.Args()
	os.Args = append([]string{os.Args[0]}, args...)
	ret := mainRet()
	fmt.Println(ret)
}
