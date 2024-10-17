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

// Encode the BridgeMasterSpec into netlink attributes.
func (a bridgeMaster) Encode() ([]byte, error) {
	bridge := a.BridgeMasterSpec

	encoder := netlink.NewAttributeEncoder()

	stpEnabled := 0
	if bridge.STP.Enabled {
		stpEnabled = 1
	}

	vlanFiltering := 0
	if bridge.VLAN.FilteringEnabled {
		vlanFiltering = 1
	}

	vlanStats := 0
	if bridge.VLAN.StatsEnabled {
		vlanStats = 1
	}

	vlanStatsPerPort := 0
	if bridge.VLAN.StatsPerPort {
		vlanStatsPerPort = 1
	}

	encoder.Uint32(unix.IFLA_BR_STP_STATE, uint32(stpEnabled))
	encoder.Uint8(unix.IFLA_BR_VLAN_FILTERING, uint8(vlanFiltering))
	if vlanFiltering != 0 {
		encoder.Uint16(unix.IFLA_BR_VLAN_DEFAULT_PVID, bridge.VLAN.DefaultPVID)
	}
	encoder.Uint8(unix.IFLA_BR_VLAN_STATS_ENABLED, uint8(vlanStats))
	encoder.Uint8(unix.IFLA_BR_VLAN_STATS_PER_PORT, uint8(vlanStatsPerPort))

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
		case unix.IFLA_BR_VLAN_DEFAULT_PVID:
			bridge.VLAN.DefaultPVID = decoder.Uint16()
		case unix.IFLA_BR_VLAN_STATS_ENABLED:
			bridge.VLAN.StatsEnabled = decoder.Uint8() == 1
		case unix.IFLA_BR_VLAN_STATS_PER_PORT:
			bridge.VLAN.StatsPerPort = decoder.Uint8() == 1
		}
	}

	return decoder.Err()
}
