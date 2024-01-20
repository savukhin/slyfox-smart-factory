package server

import (
	"errors"
	"fmt"
)

func isConnect(controlHeader byte) bool {
	return controlHeader == 0x10
}

func isPublish(controlHeader byte) bool {
	return (controlHeader & 0x30) == 0x30
}

func getPacketLength(packet []byte) (result int, shiftEnd int) {
	return calculateVarLenFromShift(packet, 1)
}

func calculateVarLenFromShift(packet []byte, shift int) (result int, shiftEnd int) {
	result = 0
	shiftEnd = shift

	for {
		b := packet[shiftEnd] & 0x7F
		result = (result << 8) | int(b)

		if (packet[shiftEnd] >> 7) == 0 {
			return
		}
		shiftEnd += 1
	}
}

func calculateNBytesFromShift(packet []byte, shift, n int) (result int, shiftEnd int) {
	result = 0
	shiftEnd = shift - 1
	_ = packet[shift+n-1] // compilier-only optimization

	for i := 0; i < n; i++ {
		b := int(packet[shift+i])
		result = (result << 8) | b
		shiftEnd++
	}
	return
}

func extractCredentials(packet []byte) (username, hashedPassword string, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("err extracting credentials %v", r)
		}
	}()

	_, packetLenEnd := getPacketLength(packet)
	protocolLen, protocolLenEnd := calculateNBytesFromShift(packet, packetLenEnd+1, 2)

	versionInd := protocolLenEnd + 1 + protocolLen
	version := packet[versionInd]
	connectFlags := packet[versionInd+1]

	hasUsername := (connectFlags & 0x80) != 0
	if !hasUsername {
		err = errors.New("no username in connect packet")
		return
	}
	hasPassword := (connectFlags & 0x40) != 0
	if !hasPassword {
		err = errors.New("no password in connect packet")
		return
	}

	var clientIDStart int
	if version == 5 {
		propertiesStart := versionInd + 1 + 3 // version + flags
		propertiesLen, propertiesLenEnd := calculateNBytesFromShift(packet, propertiesStart, 1)
		clientIDStart = propertiesLenEnd + 1 + propertiesLen
	} else {
		clientIDStart = versionInd + 1 + 3 // version + flags
	}

	cliendIDLength, cliendIDLenEnd := calculateNBytesFromShift(packet, clientIDStart, 2)

	usernameStart := cliendIDLenEnd + 1 + cliendIDLength
	usernameLength, usernameLenEnd := calculateNBytesFromShift(packet, usernameStart, 2)

	passwordStart := usernameLenEnd + 1 + usernameLength
	_, passwordLenEnd := calculateNBytesFromShift(packet, passwordStart, 2)

	username = string(packet[usernameLenEnd+1 : passwordStart])
	hashedPassword = string(packet[passwordLenEnd+1:])

	return
}
