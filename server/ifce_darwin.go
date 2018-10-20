// +build darwin
package main

import (
	"fmt"
	"github.com/Doridian/wsvpn/shared"
	"github.com/songgao/water"
	"net"
)

func configIface(dev *water.Interface, ipConfig bool, mtu int, ipClient net.IP, ipServer net.IP) error {
	err := shared.ExecCmd("ifconfig", dev.Name(), "mtu", fmt.Sprintf("%d", mtu))
	if err != nil {
		return err
	}

	if !ipConfig {
		return shared.ExecCmd("ifconfig", dev.Name(), "up")
	}

	if dev.IsTAP() {
		return shared.ExecCmd("ifconfig", dev.Name(), ipServer.String(), "up")
	}
	return shared.ExecCmd("ifconfig", dev.Name(), ipServer.String(), ipClient.String(), "up")
}

func extendTAPConfig(config *water.Config) {

}
