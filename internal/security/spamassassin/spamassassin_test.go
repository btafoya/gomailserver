package spamassassin

import (
	"strconv"
	"strings"
	"testing"

	"net"
)

// mockSpamdServer simulates a SpamAssassin spamd server for testing
type mockSpamdServer struct {
	response string
}

func (m *mockSpamdServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the request (simplified)
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		return
	}

	// Send mock response
	conn.Write([]byte(m.response))
}

func TestSpamAssassin_Scan(t *testing.T) {
	// Mock spamd response with proper format
	mockResponse := "SPAMD/1.2 0 EX_OK\r\nSpam: True ; 5.5 / 5.0\r\n\r\n5.5 TEST_RULE Test rule description\r\n"

	server := &mockSpamdServer{response: mockResponse}

	// Start mock server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		server.handleConnection(conn)
	}()

	// Get the port
	addr := listener.Addr().String()
	parts := strings.Split(addr, ":")
	portStr := parts[len(parts)-1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatal(err)
	}

	// Test the scanner
	scanner := NewSpamAssassin("127.0.0.1", port)
	result, err := scanner.Scan([]byte("test message"))

	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if !result.IsSpam {
		t.Error("Expected message to be marked as spam")
	}

	if result.Score != 5.5 {
		t.Errorf("Expected score 5.5, got %f", result.Score)
	}

	if result.Threshold != 5.0 {
		t.Errorf("Expected threshold 5.0, got %f", result.Threshold)
	}

	if len(result.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(result.Rules))
	}

	if result.Rules[0].Rule != "TEST_RULE" {
		t.Errorf("Expected rule 'TEST_RULE', got '%s'", result.Rules[0].Rule)
	}
}
