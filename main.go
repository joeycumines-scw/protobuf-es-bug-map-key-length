package main

import (
	"encoding/base64"
	"fmt"
	"github.com/joeycumines-scw/protobuf-es-bug-map-key-length/go/example"
	"google.golang.org/protobuf/proto"
	"os"
	"os/exec"
	"strings"
)

func main() {
	maxWorkingTargetSize := -1
	minBreakingTargetSize := -1
	for _, v := range []int{122, 123, 124, 125} {
		fmt.Printf("Map key size: %d\n", v)
		msg := &example.Message{
			FieldOne: []*example.RepeatedMessage{
				{FieldTwo: map[string]uint32{
					strings.Repeat("a", v): 5,
				}},
			},
		}
		if check(msg) == nil {
			maxWorkingTargetSize = max(maxWorkingTargetSize, v)
		} else {
			minBreakingTargetSize = max(minBreakingTargetSize, v)
		}
	}
	fmt.Printf("Max working target size: %d\n", maxWorkingTargetSize)
	fmt.Printf("Min breaking target size: %d\n", minBreakingTargetSize)
}

func check(msg *example.Message) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Encoded message size: %d bytes\n", len(b))
	arg := base64.StdEncoding.EncodeToString(b)
	fmt.Printf("Encoded message, base64 encoded: %s\n", arg)
	cmd := exec.Command("node", "build/src/index.js", arg)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
