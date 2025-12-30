package antivirus

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type ClamAV struct {
	socketPath string
}

func NewClamAV(socketPath string) *ClamAV {
	return &ClamAV{socketPath: socketPath}
}

func (c *ClamAV) Scan(data []byte) (*ScanResult, error) {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clamd: %w", err)
	}
	defer conn.Close()

	// Send INSTREAM command
	_, err = fmt.Fprintf(conn, "zINSTREAM\x00")
	if err != nil {
		return nil, fmt.Errorf("failed to send INSTREAM command: %w", err)
	}

	// Send data in chunks
	chunkSize := 2048
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		// Write chunk size (4 bytes, network byte order)
		if err := binary.Write(conn, binary.BigEndian, uint32(len(chunk))); err != nil {
			return nil, fmt.Errorf("failed to write chunk size: %w", err)
		}
		if _, err := conn.Write(chunk); err != nil {
			return nil, fmt.Errorf("failed to write chunk data: %w", err)
		}
	}

	// End stream
	if err := binary.Write(conn, binary.BigEndian, uint32(0)); err != nil {
		return nil, fmt.Errorf("failed to write end of stream: %w", err)
	}

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\x00')
	if err != nil {
		return nil, fmt.Errorf("failed to read response from clamd: %w", err)
	}

	return parseResponse(response), nil
}

type ScanResult struct {
	Clean bool
	Virus string
	Error string
}

func parseResponse(response string) *ScanResult {
	response = strings.TrimSuffix(response, "\x00")

	if strings.HasSuffix(response, "OK") {
		return &ScanResult{Clean: true}
	}

	if strings.Contains(response, "FOUND") {
		parts := strings.Split(response, ":")
		virus := strings.TrimSpace(strings.TrimSuffix(parts[len(parts)-1], "FOUND"))
		return &ScanResult{Clean: false, Virus: virus}
	}

	return &ScanResult{Clean: false, Error: response}
}
