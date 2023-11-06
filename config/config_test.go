package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	defaultString = "field1"
	defaultInt    = 1
	defaultDbUrl  = "postgres://localhost:5432"
	defaultAPIKey = "abc123"
)

// mockConfig is a mock struct implementing Config interface for testing purposes.
type mockConfig struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

// Default returns a default mock configuration.
func (c mockConfig) Default() Config {
	return &mockConfig{
		Field1: defaultString,
		Field2: defaultInt,
	}
}

type mockAppConfig struct {
	DatabaseUrl string
	ApiKey      string
}

func (c mockAppConfig) Default() Config {
	return &mockAppConfig{
		DatabaseUrl: "postgres://localhost:5432",
		ApiKey:      "abc123",
	}
}

func TestNewConfig_AppConfig(t *testing.T) {
	defaultAppConfig := NewConfig(&mockAppConfig{})

	if defaultAppConfig.DatabaseUrl != defaultDbUrl || defaultAppConfig.ApiKey != defaultAPIKey {
		t.Errorf("NewConfig didn't return the expected default values, got: %+v", defaultAppConfig)
	}
}

func TestNewConfig(t *testing.T) {
	config := &mockConfig{}

	newConfig := NewConfig(config)

	if newConfig.Field1 != defaultString || newConfig.Field2 != defaultInt {
		t.Errorf("NewConfig didn't return the expected default values, got: %+v", newConfig)

	}
}
func TestLoadConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "go-config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	var conf = &mockConfig{}
	defaultConf := conf.Default()

	err = SaveConfig(tmpFile.Name(), defaultConf)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	err = LoadConfig(tmpFile.Name(), conf)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if conf.Field1 != defaultString || conf.Field2 != defaultInt {
		t.Errorf("LoadConfig didn't load the expected values")
	}
}

func TestSaveConfig(t *testing.T) {
	conf := &mockConfig{
		Field1: "test1",
		Field2: 1,
	}

	tmpDir, err := os.MkdirTemp("", "go-config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "go-config.json")

	conf.Field1 = "changedString"
	conf.Field2 = 2

	err = SaveConfig(configPath, conf)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	err = LoadConfig(configPath, conf)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if conf.Field1 != "changedString" || conf.Field2 != 2 {
		t.Errorf("SaveConfig didn't save the expected values")
	}
}

func TestInitConfig_FileDoesNotExist(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")

	if _, err := os.Stat(configPath); err == nil {
		t.Fatalf("Config file already exists: %s", configPath)
	}

	defaultConf := &mockConfig{}
	err = InitConfig(configPath, defaultConf)
	if err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	loadedConf := &mockConfig{}
	err = LoadConfig(configPath, loadedConf)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConf.Field1 != defaultString || loadedConf.Field2 != defaultInt {
		t.Errorf("InitConfig did not create a config file with the default configuration. Got: %+v", loadedConf)
	}
}

func TestInitConfig_FileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")

	initialConf := &mockConfig{}
	err = InitConfig(configPath, initialConf)
	if err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	err = InitConfig(configPath, initialConf)
	if err == nil {
		t.Fatalf("Second call to InitConfig passed when it should have failed")
	}

	loadedConf := &mockConfig{}
	err = LoadConfig(configPath, loadedConf)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	fmt.Println(loadedConf)
	if loadedConf.Field1 != defaultString || loadedConf.Field2 != defaultInt {
		t.Errorf("InitConfig altered an existing config file. Got: %+v", loadedConf)
	}
}
