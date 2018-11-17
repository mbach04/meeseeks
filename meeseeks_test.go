package main

import (
	"testing"
)

func TestConfig(t *testing.T) {
	//test default items
	testConfig, err := readConfig("fake_config", map[string]interface{}{
		"HOSTNAME": "localhost",
		"API_PORT": 8080,
		"DEBUG":    true,
		"TLS_CERT": "cert.pem",
		"TLS_KEY":  "key.pem",
	})

	if testConfig.GetBool("DEBUG") != true {
		t.Error("Error setting default value for DEBUG")
	}

	if testConfig.GetString("API_PORT") != "8080" {
		t.Error("Error setting default value for API_PORT")
	}

	//test items from a config file
	testConfig, err = readConfig("config", map[string]interface{}{
		"HOSTNAME": "localhost",
		"API_PORT": 8080,
		"DEBUG":    true,
		"TLS_CERT": "cert.pem",
		"TLS_KEY":  "key.pem",
	})

	if err != nil {
		t.Error("Error reading config file")
	}

	//ensure the config file has a different port than the default
	if testConfig.GetString("API_PORT") != "9191" {
		t.Error("Error reading config file for API_PORT")
	}

}

func TestGenerateHTTPSCert(t *testing.T) {

}

func TestStartEcho(t *testing.T) {

}
