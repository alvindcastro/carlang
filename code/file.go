package code

import (
	"bytes"
	"fmt"
	"os"
)

var Magic = []byte("CARLBC1\n")

func WriteFile(path string, instructions Instructions) error {
	data := append([]byte{}, Magic...)
	data = append(data, instructions...)
	return os.WriteFile(path, data, 0o644)
}

func ReadFile(path string) (Instructions, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(data, Magic) {
		return nil, fmt.Errorf("%s is not a Carlang bytecode file", path)
	}
	return Instructions(data[len(Magic):]), nil
}
