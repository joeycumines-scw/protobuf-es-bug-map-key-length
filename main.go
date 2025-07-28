package main

import (
  "encoding/base64"
  "fmt"
  "os"
  "os/exec"
  "strings"

  "github.com/joeycumines-scw/protobuf-es-bug-map-key-length/go/example"
  "google.golang.org/protobuf/proto"
)

func main() {
  fmt.Println("=== Message Test ===")
  keyLengths := []int{120, 121, 122, 123, 124, 125, 126, 127, 128}
  fmt.Println("\nTesting Message")
  testEncode(keyLengths)
  fmt.Println("\nFinding exact boundary")
  findExactBoundary()
}

func testEncode(keyLengths []int) {
  for _, keyLen := range keyLengths {
    fmt.Printf("Key length %d: ", keyLen)
    msg := &example.Message{
      MyMap: map[string]uint32{
        strings.Repeat("a", keyLen): 1,
      },
    }
    if testMessage(msg) {
      fmt.Println("✅ PASS")
    } else {
      fmt.Println("❌ FAIL")
    }
  }
}

func findExactBoundary() {
  fmt.Println("Finding exact failure boundary...")

  // Binary search for exact boundary
  low, high := 0, 10000
  lastWorking := -1
  firstFailing := -1

  for low <= high {
    mid := (low + high) / 2
    fmt.Printf("Testing key length %d: ", mid)
    msg := &example.Message{
      MyMap: map[string]uint32{
        strings.Repeat("a", mid): 1,
      },
    }
    if testMessage(msg) {
      fmt.Println("✅ PASS")
      lastWorking = mid
      low = mid + 1
    } else {
      fmt.Println("❌ FAIL")
      firstFailing = mid
      high = mid - 1
    }
  }
  fmt.Printf("\n🎯 EXACT BOUNDARY:\n")
  fmt.Printf("   Last working key length: %d\n", lastWorking)
  fmt.Printf("   First failing key length: %d\n", firstFailing)
  // Test the exact boundary with wire format analysis
  if lastWorking != -1 && firstFailing != -1 {
    fmt.Printf("\n📊 WIRE FORMAT ANALYSIS:\n")
    analyzeWireFormat(lastWorking, "WORKING")
    analyzeWireFormat(firstFailing, "FAILING")
  }
}

func analyzeWireFormat(keyLen int, label string) {
  msg := &example.Message{
    MyMap: map[string]uint32{
      strings.Repeat("a", keyLen): 1,
    },
  }

  b, err := proto.Marshal(msg)
  if err != nil {
    fmt.Printf("Error marshaling %s case: %v\n", label, err)
    return
  }

  fmt.Printf("%s (key len %d): %d bytes total\n", label, keyLen, len(b))
  fmt.Printf("Hex: %x\n", b)

  // Calculate expected map entry size
  mapEntrySize := 1 + getVarintLen(keyLen) + keyLen + 1 + 1 // tag + key_len + key + tag + value
  fmt.Printf("Expected map entry size: %d bytes\n", mapEntrySize)
  fmt.Printf("Map entry varint encoding: %s\n", getVarintInfo(mapEntrySize))
  fmt.Printf("Base64: %s\n", base64.StdEncoding.EncodeToString(b))
  fmt.Println()
}

func getVarintLen(val int) int {
  if val < 128 {
    return 1
  }
  return 2
}

func getVarintInfo(val int) string {
  if val < 128 {
    return fmt.Sprintf("1 byte (0x%02x)", val)
  } else {
    return fmt.Sprintf("2+ bytes (0x%02x 0x%02x)", val&0x7f|0x80, val>>7)
  }
}

func testMessage(msg proto.Message) bool {
  b, err := proto.Marshal(msg)
  if err != nil {
    return false
  }
  arg := base64.StdEncoding.EncodeToString(b)
  cmd := exec.Command("node", "build/src/index.js", arg)
  cmd.Stderr = os.Stderr
  cmd.Stdout = os.Stdout
  return cmd.Run() == nil
}
