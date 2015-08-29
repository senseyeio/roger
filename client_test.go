package roger

import "testing"

func TestConnection(t *testing.T) {
	if _, err := NewRClient("localhost", 6311); err != nil {
		t.Error("Failed to connect to RServe: " + err.Error())
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

func TestSecureConnection(t *testing.T) {
	if _, err := NewRClientWithAuth("localhost", 6312, "roger", "testpassword"); err != nil {
		t.Error("Failed to connect to secure RServe: " + err.Error())
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
