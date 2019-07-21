package ipmi

import (
	"github.com/gebn/bmc/pkg/layerexts"

	"github.com/google/gopacket"
)

var (
	LayerTypeSessionSelector = gopacket.RegisterLayerType(
		1000,
		gopacket.LayerTypeMetadata{
			Name:    "IPMI Session Selector",
			Decoder: layerexts.BuildDecoder(&SessionSelector{}),
		},
	)
	LayerTypeV1Session = gopacket.RegisterLayerType(
		1001,
		gopacket.LayerTypeMetadata{
			Name:    "Session v1.5",
			Decoder: layerexts.BuildDecoder(&V1Session{}),
		},
	)
	LayerTypeGetChannelAuthenticationCapabilitiesReq = gopacket.RegisterLayerType(
		1002,
		gopacket.LayerTypeMetadata{
			Name: "Get Channel Authentication Capabilities Request",
		},
	)
	LayerTypeGetChannelAuthenticationCapabilitiesRsp = gopacket.RegisterLayerType(
		1003,
		gopacket.LayerTypeMetadata{
			Name:    "Get Channel Authentication Capabilities Response",
			Decoder: layerexts.BuildDecoder(&GetChannelAuthenticationCapabilitiesRsp{}),
		},
	)
	LayerTypeV2Session = gopacket.RegisterLayerType(
		1004,
		gopacket.LayerTypeMetadata{
			Name: "Session v2.0",
			// by default this layer can only encode and decode unauthenticated
			// packets; to deal with authenticated packets, the
			// IntegrityAlgorithm attribute must be set
			Decoder: layerexts.BuildDecoder(&V2Session{}),
		},
	)
	LayerTypeOpenSessionReq = gopacket.RegisterLayerType(
		1005,
		gopacket.LayerTypeMetadata{
			Name: "RMCP+ Open Session Request",
		},
	)
	LayerTypeOpenSessionRsp = gopacket.RegisterLayerType(
		1006,
		gopacket.LayerTypeMetadata{
			Name:    "RMCP+ Open Session Response",
			Decoder: layerexts.BuildDecoder(&OpenSessionRsp{}),
		},
	)
	LayerTypeRAKPMessage1 = gopacket.RegisterLayerType(
		1007,
		gopacket.LayerTypeMetadata{
			Name: "RAKP Message 1",
		},
	)
	LayerTypeRAKPMessage2 = gopacket.RegisterLayerType(
		1008,
		gopacket.LayerTypeMetadata{
			Name:    "RAKP Message 2",
			Decoder: layerexts.BuildDecoder(&RAKPMessage2{}),
		},
	)
	LayerTypeRAKPMessage3 = gopacket.RegisterLayerType(
		1009,
		gopacket.LayerTypeMetadata{
			Name: "RAKP Message 3",
		},
	)
	LayerTypeRAKPMessage4 = gopacket.RegisterLayerType(
		1010,
		gopacket.LayerTypeMetadata{
			Name:    "RAKP Message 4",
			Decoder: layerexts.BuildDecoder(&RAKPMessage4{}),
		},
	)
	layerTypeAES128CBC = gopacket.RegisterLayerType(
		1011,
		gopacket.LayerTypeMetadata{
			Name: "AES-128-CBC Encrypted IPMI Message",
			// decoder not specified here as default struct not usable
		},
	)
	LayerTypeMessage = gopacket.RegisterLayerType(
		1012,
		gopacket.LayerTypeMetadata{
			Name:    "IPMI Message",
			Decoder: layerexts.BuildDecoder(&Message{}),
		},
	)
	LayerTypeCloseSessionReq = gopacket.RegisterLayerType(
		1013,
		gopacket.LayerTypeMetadata{
			Name: "Close Session Request",
		},
	)
	LayerTypeGetSystemGUIDRsp = gopacket.RegisterLayerType(
		1014,
		gopacket.LayerTypeMetadata{
			Name:    "Get System GUID Response",
			Decoder: layerexts.BuildDecoder(&GetSystemGUIDRsp{}),
		},
	)
	LayerTypeGetDeviceIDRsp = gopacket.RegisterLayerType(
		1015,
		gopacket.LayerTypeMetadata{
			Name:    "Get Device ID Response",
			Decoder: layerexts.BuildDecoder(&GetDeviceIDRsp{}),
		},
	)
)