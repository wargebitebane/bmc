package ipmi

import (
	"encoding/binary"
	"fmt"

	"github.com/google/gopacket"
)

// RAKPMessage1 represents a RAKP Message 1, defined in 13.20 of the spec. It
// begins the session authentication process.
type RAKPMessage1 struct {

	// Tag is an arbitrary 8-bit quantity used by the remote console to match
	// this message with RAKP Message 2.
	Tag uint8

	// ManagedSystemSessionID is the session ID returned by the BMC in the RMCP+
	// Open Session Response message.
	ManagedSystemSessionID uint32

	// RemoteConsoleRandom is a 16-byte random value selected by the remote
	// console. Although it is referred to as a number, its byte order is not
	// reversed on the wire.
	RemoteConsoleRandom [16]byte

	// PrivilegeLevelLookup indicates whether to use both the MaxPrivilegeLevel
	// and Username to search for the relevant user entry. If this is true and
	// the username is empty, we effectively use role-based authentication. If
	// this is false, the supplied MaxPrivilegeLevel will be ignored when
	// searching for the Username.
	PrivilegeLevelLookup bool

	// MaxPrivilegeLevel is the upper privilege limit for the session. If
	// PrivilegeLevelLookup is true, it is also used in the user entry lookup.
	// Regardless of this value, if PrivilegeLevelLookup is false, the channel
	// or user privilege level limit may further constrain allowed commands.
	MaxPrivilegeLevel PrivilegeLevel

	// Username is the name of the user to search for in user entries. Only
	// ASCII characters allowed (excluding \0), maximum length 16.
	Username string
}

func (*RAKPMessage1) LayerType() gopacket.LayerType {
	return LayerTypeRAKPMessage1
}

func (r *RAKPMessage1) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	if len(r.Username) > 16 {
		return fmt.Errorf("Username cannot be more than 16 characters long, got %v", len(r.Username))
	}
	d, err := b.PrependBytes(28 + len(r.Username))
	if err != nil {
		return err
	}
	d[0] = r.Tag
	d[1] = 0x00
	d[2] = 0x00
	d[3] = 0x00
	binary.LittleEndian.PutUint32(d[4:8], r.ManagedSystemSessionID)
	copy(d[8:24], r.RemoteConsoleRandom[:])
	d[24] = uint8(r.MaxPrivilegeLevel) & 0xF
	if !r.PrivilegeLevelLookup {
		d[24] |= 1 << 4
	}
	d[25] = 0x00
	d[26] = 0x00
	d[27] = uint8(len(r.Username))
	if len(r.Username) > 0 {
		copy(d[28:], []byte(r.Username))
	}
	return nil
}