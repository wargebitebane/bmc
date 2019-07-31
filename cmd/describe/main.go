package main

// Describe shows various information about a BMC using the ASF Presence Pong,
// Get Channel Authentication Capabilities, Get System GUID and Get Device ID
// commands.

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/gebn/bmc"
	"github.com/gebn/bmc/internal/pkg/transport"
	"github.com/gebn/bmc/pkg/dcmi"
	"github.com/gebn/bmc/pkg/ipmi"

	"github.com/alecthomas/kingpin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	argBMCAddr = kingpin.Arg("addr", "IP[:port] of the BMC to describe.").
			Required().
			String()
	flgUsername = kingpin.Flag("username", "The username to connect as.").
			Required().
			String()
	flgPassword = kingpin.Flag("password", "The password of the user to connect as.").
			Required().
			String()
)

func main() {
	kingpin.Parse()

	machine, err := bmc.DialV2(*argBMCAddr) // TODO change to Dial (need to implement v1.5 sessionless communication...)
	if err != nil {
		log.Fatal(err)
	}
	defer machine.Close()

	log.Printf("connected to %v over IPMI v%v", machine.Address(), machine.Version())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if pong, err := presencePing(ctx, machine); err != nil {
		log.Printf("failed to get presence pong capabilities: %v", err)
	} else {
		printPong(pong)
	}

	if caps, err := machine.GetChannelAuthenticationCapabilities(ctx,
		&ipmi.GetChannelAuthenticationCapabilitiesReq{
			ExtendedData:      true, // only has effect if v2.0
			Channel:           ipmi.ChannelPresentInterface,
			MaxPrivilegeLevel: ipmi.PrivilegeLevelAdministrator,
		}); err != nil {
		log.Printf("failed to get channel auth capabilities: %v", err)
	} else {
		printChannelAuthCaps(caps)
	}

	if guid, err := machine.GetSystemGUID(ctx); err != nil {
		log.Printf("failed to get system GUID: %v", err)
	} else {
		printSystemGUID(guid)
	}

	c, m, p := getDCMICaps(ctx, machine)
	if c != nil {
		printDCMICaps(c)
	}
	if m != nil {
		printDCMIPlatformAttrs(m)
	}
	if p != nil {
		printDCMIPowerPeriods(p)
	}

	sess, err := machine.NewSession(ctx, *flgUsername, []byte(*flgPassword))
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close(ctx)

	if id, err := sess.GetDeviceID(ctx); err != nil {
		log.Printf("failed to get device id: %v", err)
	} else {
		printDeviceID(id)
	}

	if status, err := sess.GetChassisStatus(ctx); err != nil {
		log.Printf("failed to get chassis status: %v", err)
	} else {
		printChassisStatus(status)
	}
}

func presencePing(ctx context.Context, t transport.Transport) (*layers.ASFPresencePong, error) {
	asfRmcp := &layers.RMCP{
		Version:  layers.RMCPVersion1,
		Sequence: 0xFF, // do not send an ACK
		Class:    layers.RMCPClassASF,
	}
	asf := &layers.ASF{
		ASFDataIdentifier: layers.ASFDataIdentifierPresencePing,
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, asfRmcp, asf); err != nil {
		return nil, err
	}
	bytes, err := t.Send(ctx, buf.Bytes())
	if err != nil {
		return nil, err
	}
	packet := gopacket.NewPacket(bytes, layers.LayerTypeRMCP, gopacket.DecodeOptions{
		Lazy:   true,
		NoCopy: true,
	})
	pongLayer := packet.Layer(layers.LayerTypeASFPresencePong)
	if pongLayer == nil {
		return nil, fmt.Errorf("no presence pong layer in response")
	}
	return pongLayer.(*layers.ASFPresencePong), nil
}

func printPong(p *layers.ASFPresencePong) {
	fmt.Println("ASF Presence Pong capabilities:")
	fmt.Printf("\tIPMI:               %v\n", p.IPMI)
	fmt.Printf("\tASF v1.0:           %v\n", p.ASFv1)
	fmt.Printf("\tASF security exts:  %v\n", p.SecurityExtensions) // means the BMC uses the secure port in addition to the normal one
	fmt.Printf("\tDASH:               %v\n", p.DASH)
	fmt.Printf("\tDCMI:               %v\n", p.SupportsDCMI())
}

func printChannelAuthCaps(c *ipmi.GetChannelAuthenticationCapabilitiesRsp) {
	fmt.Println("Channel Authentication Capabilities:")
	fmt.Printf("\tChannel:            %v\n", c.Channel)
	fmt.Printf("\tExtended:           %v\n", c.ExtendedCapabilities)
	fmt.Printf("\tSupportsV2:         %v\n", c.SupportsV2)
	fmt.Printf("\tK_G configured:     %v\n", c.TwoKeyLogin)
	fmt.Printf("\tPer-message auth:   %v\n", c.PerMessageAuthentication)
	fmt.Printf("\tUser-level auth:    %v\n", c.UserLevelAuthentication)
	fmt.Printf("\tNon-null usernames: %v\n", c.NonNullUsernamesEnabled)
	fmt.Printf("\tNull usernames:     %v\n", c.NullUsernamesEnabled)
	fmt.Printf("\tAnon login:         %v\n", c.AnonymousLoginEnabled)
	fmt.Printf("\tOEM:                %v\n", c.OEM)
}

func printSystemGUID(guid [16]byte) {
	buf := [36]byte{}
	encodeHex(buf[:], guid)
	fmt.Println("System:")
	fmt.Printf("\tGUID:               %v\n", string(buf[:]))
}

func encodeHex(dst []byte, guid [16]byte) {
	hex.Encode(dst, guid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], guid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], guid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], guid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], guid[10:])
}

func getDCMICaps(ctx context.Context, s bmc.Sessionless) (
	*dcmi.GetDCMICapabilitiesInfoSupportedCapabilitiesRsp,
	*dcmi.GetDCMICapabilitiesInfoMandatoryPlatformAttrsRsp,
	*dcmi.GetDCMICapabilitiesInfoEnhancedSystemPowerStatisticsAttrsRsp,
) {
	commander := dcmi.SessionlessCommander(s)
	c, err := commander.GetDCMICapabilitiesInfoSupportedCapabilities(ctx)
	if err != nil {
		log.Printf("failed to fetch DCMI supported capabilities: %v", err)
	}
	m, err := commander.GetDCMICapabilitiesInfoMandatoryPlatformAttrs(ctx)
	if err != nil {
		log.Printf("failed to fetch DCMI mandatory platform attrs: %v", err)
	}
	p, err := commander.GetDCMICapabilitiesInfoEnhancedSystemPowerStatisticsAttrs(ctx)
	if err != nil {
		log.Printf("failed to fetch DCMI enhanced power stats attrs: %v", err)
	}
	return c, m, p
}

func printDCMICaps(c *dcmi.GetDCMICapabilitiesInfoSupportedCapabilitiesRsp) {
	fmt.Println("DCMI Capabilities:")
	fmt.Printf("\tMajor version:      %v\n", c.MajorVersion)
	fmt.Printf("\tMinor version:      %v\n", c.MinorVersion)
	fmt.Printf("\tSupports pwr mgmt:  %v\n", c.PowerManagement)
}

func printDCMIPlatformAttrs(m *dcmi.GetDCMICapabilitiesInfoMandatoryPlatformAttrsRsp) {
	fmt.Println("DCMI Mandatory Platform Attributes:")
	fmt.Printf("\tMax SEL entries:    %v\n", m.SELMaxEntries)
	fmt.Printf("\tTemp sampling freq: %v\n", m.TemperatureSamplingFrequency)
}

func printDCMIPowerPeriods(p *dcmi.GetDCMICapabilitiesInfoEnhancedSystemPowerStatisticsAttrsRsp) {
	fmt.Println("DCMI Power Average Time Periods:")
	for _, duration := range p.PowerRollingAvgTimePeriods {
		fmt.Printf("\t%v\n", duration)
	}
}

func printDeviceID(id *ipmi.GetDeviceIDRsp) {
	fmt.Println("Device:")
	fmt.Printf("\tID:                 %v\n", id.ID)
	fmt.Printf("\tRevision:           %v\n", id.Revision)
	fmt.Printf("\tManufacturer:       %v\n", id.Manufacturer)
	fmt.Printf("\tProduct:            %v\n", id.Product)
	fmt.Printf("\tFirmware (major):   %v\n", id.MajorFirmwareRevision)
	fmt.Printf("\tFirmware (minor):   %v\n", id.MinorFirmwareRevision)
	fmt.Printf("\tFirmware (aux):     %v\n", hex.EncodeToString(id.AuxiliaryFirmwareRevision[:]))
	fmt.Printf("\tFirmware:           %v\n", bmc.FirmwareVersion(id))
}

func printChassisStatus(status *ipmi.GetChassisStatusRsp) {
	fmt.Println("Chassis:")
	fmt.Printf("\tPowered on:         %v\n", status.PoweredOn)
	fmt.Printf("\tOn power restore:   %v\n", status.PowerRestorePolicy)
	fmt.Printf("\tIdentification:     %v\n", status.ChassisIdentifyState)
	fmt.Printf("\tIntrusion:          %v\n", status.Intrusion)
	fmt.Printf("\tPower fault:        %v\n", status.PowerFault)
	fmt.Printf("\tCooling fault:      %v\n", status.CoolingFault)
	fmt.Printf("\tDrive fault:        %v\n", status.DriveFault)
}
