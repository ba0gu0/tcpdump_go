//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

//go:embed libs/wpcap.dll
var wpcapDLL []byte

//go:embed libs/Packet.dll
var packetDLL []byte

var content embed.FS

func loadEmbeddedDLL(dllName string, dllBytes []byte) error {
	dllPath := filepath.Join(os.TempDir(), dllName)
	err := ioutil.WriteFile(dllPath, dllBytes, 0644)
	if err != nil {
		return err
	}

	if os.IsNotExist(err) {
		return fmt.Errorf("DLL file does not exist after writing: %s", dllPath)
	}

	_, err = syscall.LoadDLL(dllPath)
	if err != nil {
		return fmt.Errorf("Failed to load %s: %v", dllPath, err)
	}
	return nil
}

func init() {

	err := loadEmbeddedDLL("Packet.dll", packetDLL)
	if err != nil {
		log.Fatalf("Error loading Packet.dll: %v\n", err)
	}

	err = loadEmbeddedDLL("wpcap.dll", wpcapDLL)
	if err != nil {
		log.Fatalf("Error loading wpcap.dll: %v\n", err)
	}

	fmt.Println("DLLs loaded successfully")
}
