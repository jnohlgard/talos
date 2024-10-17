// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package network

import (
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"

	"github.com/siderolabs/talos/pkg/machinery/resources/network"
)

// BridgeMasterSpec adapter provides encoding/decoding to netlink structures.
//
//nolint:revive
func BridgeMasterSpec(r *network.BridgeMasterSpec) bridgeMaster {
	return bridgeMaster{
		BridgeMasterSpec: r,
	}
}

// bridgeMaster contains the bridge master spec and provides methods for encoding/decoding it to netlink structures.
type bridgeMaster struct {
	*network.BridgeMasterSpec
}

func bool01(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// Encode the BridgeMasterSpec into netlink attributes.
func (a bridgeMaster) Encode() ([]byte, error) {
	bridge := a.BridgeMasterSpec

	encoder := netlink.NewAttributeEncoder()

	encoder.Uint32(unix.IFLA_BR_STP_STATE, uint32(bool01(bridge.STP.Enabled)))
	encoder.Uint8(unix.IFLA_BR_VLAN_FILTERING, bool01(bridge.VLAN.FilteringEnabled))
	encoder.Uint8(unix.IFLA_BR_VLAN_STATS_ENABLED, bool01(bridge.VLAN.StatsEnabled))
	encoder.Uint8(unix.IFLA_BR_VLAN_STATS_PER_PORT, bool01(bridge.VLAN.StatsPerPort))
	if bridge.VLAN.FilteringEnabled {
		encoder.Uint16(unix.IFLA_BR_VLAN_DEFAULT_PVID, bridge.VLAN.DefaultPVID)
	}

	return encoder.Encode()
}

// Decode the BridgeMasterSpec from netlink attributes.
func (a bridgeMaster) Decode(data []byte) error {
	bridge := a.BridgeMasterSpec

	decoder, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}

	for decoder.Next() {
		switch decoder.Type() {
		case unix.IFLA_BR_STP_STATE:
			bridge.STP.Enabled = decoder.Uint32() == 1
		case unix.IFLA_BR_VLAN_FILTERING:
			bridge.VLAN.FilteringEnabled = decoder.Uint8() == 1
		case unix.IFLA_BR_VLAN_STATS_ENABLED:
			bridge.VLAN.StatsEnabled = decoder.Uint8() == 1
		case unix.IFLA_BR_VLAN_STATS_PER_PORT:
			bridge.VLAN.StatsPerPort = decoder.Uint8() == 1
		case unix.IFLA_BR_VLAN_DEFAULT_PVID:
			bridge.VLAN.DefaultPVID = decoder.Uint16()
		}
	}

	return decoder.Err()
}
