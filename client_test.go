package roger

import (
	"strings"
	"testing"
)

func TestConnection(t *testing.T) {
	if _, err := NewRClient("localhost", 6311); err != nil {
		t.Error("Failed to connect to RServe: " + err.Error())
	}
	if _, err := NewRClient("localhost", 6313); err != nil {
		t.Error("Failed to connect to RServe from CRAN: " + err.Error())
	}
}

func TestConnectionFailure(t *testing.T) {
	if _, err := NewRClient("localhost", 6310); err == nil {
		t.Error("Should fail when trying to connect to an incorrect port")
	}
	if _, err := NewRClient("%20", 6311); err == nil {
		t.Error("Should fail when trying to connect to an invalid host")
	}
}

func TestHandshakeFailure(t *testing.T) {
	_, err := NewRClient("localhost", 6315)
	if err == nil {
		t.Error("Should fail when trying to connect to a port not exposing RServe")
	}
	if strings.Contains(err.Error(), "Handshake") == false {
		t.Error("Should fail with a handshake error")
	}
}

func TestSecureConnection(t *testing.T) {
	if _, err := NewRClientWithAuth("localhost", 6312, "roger", "testpassword"); err != nil {
		t.Error("Failed to connect to secure RServe: " + err.Error())
	}
	if _, err := NewRClientWithAuth("localhost", 6314, "roger", "testpassword"); err != nil {
		t.Error("Failed to connect to secure RServe from CRAN: " + err.Error())
	}
}

func TestSecureConnectionFailure(t *testing.T) {
	if _, err := NewRClient("localhost", 6312); err == nil {
		t.Error("Should not have connected to a secure RServe with no username or password")
	}
	if _, err := NewRClientWithAuth("localhost", 6312, "roger", "incorrectpassword"); err == nil {
		t.Error("Should not have connected to a secure RServe with incorrect password")
	}
	if _, err := NewRClientWithAuth("localhost", 6312, "notauser", "testpassword"); err == nil {
		t.Error("Should not have connected to a secure RServe with incorrect username")
	}
}
