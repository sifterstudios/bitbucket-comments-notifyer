package main

import (
	"testing"

	"github.com/sifterstudios/bitbucket-notifier/data"
)

// Mock data.FileOrFolderExists to control its behavior
func mockFileOrFolderExists(path string) bool {
	switch path {
	case data.SecurityFile:
		return true
	case data.ConfigFile:
		return true
	default:
		return false
	}
}

func TestInitialize(t *testing.T) {
	// Backup original function and restore it after the test
	originalFileOrFolderExists := data.FileOrFolderExists
	defer func() {
		data.FileOrFolderExists = originalFileOrFolderExists
	}()

	data.FileOrFolderExists = mockFileOrFolderExists

	// Test case 1: Security file and Config file exist
	initialize()

	// Test case 2: Security file exists, but Config file doesn't
	data.FileOrFolderExists = func(path string) bool {
		return path == data.SecurityFile
	}
	initialize()

	// Test case 3: Neither Security file nor Config file exist
	data.FileOrFolderExists = func(path string) bool {
		return false
	}
	initialize()
}
