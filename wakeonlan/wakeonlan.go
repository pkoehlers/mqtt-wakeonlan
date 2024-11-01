package wakeonlan

import (
	"errors"
	"net"
)

// createMagicPacket creates the magic packet for the specified MAC address.
func createMagicPacket(mac string) ([]byte, error) {
	macBytes, err := net.ParseMAC(mac)
	if err != nil {
		return nil, errors.New("invalid MAC address")
	}

	if len(macBytes) != 6 {
		return nil, errors.New("invalid MAC address length")
	}

	packet := make([]byte, 102)
	// Fill the first 6 bytes with 0xFF
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	// Repeat the MAC address 16 times
	for i := 1; i <= 16; i++ {
		copy(packet[i*6:], macBytes)
	}
	return packet, nil
}

func SendDefaultWOL(mac string) error {
	return SendWOL("255.255.255.255:9", mac)
}

// sendWOL sends the magic packet to the broadcast address over UDP.
func SendWOL(broadcastAddr, mac string) error {
	packet, err := createMagicPacket(mac)
	if err != nil {
		return err
	}

	conn, err := net.Dial("udp", broadcastAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return err
	}
	return nil
}
