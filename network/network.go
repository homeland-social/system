package network

import (
	"fmt"
	"os"
	"strings"
)

/*
Structure of /etc/network/interfaces...

auto <foo> - interface <foo> will be started on boot, not if omitted.
iface <foo> inet dhcp
iface <foo> inet static
        address <ip>
        netmask <mask>
        gateway <gw>
*/

const DEFAULT_NETWORK_INTERFACES = "/root/etc/network/interfaces"

type Interface struct {
	Name    string `json:"name"`
	Auto    bool   `json:"auto"`
	DHCP    bool   `json:"dhcp"`
	Address string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

func getNetworkPath() string {
	path, exists := os.LookupEnv("NETWORK_INTERFACES_PATH")

	if !exists {
		path = DEFAULT_NETWORK_INTERFACES
	}

	return path
}

func ParseInterfaces() (map[string]*Interface, error) {
	path := getNetworkPath()
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ifaces := make(map[string]*Interface)
	lines := strings.Split(string(body[:]), "\n")

	var iface *Interface = nil
	var name string
	var ok bool
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")

		if parts[0] == "auto" {
			name = parts[1]
			if iface != nil {
				ifaces[iface.Name] = iface
			}
			iface, ok = ifaces[name]
			if !ok {
				iface = &Interface{Name: name}
			}
			iface.Auto = true
		} else if parts[0] == "iface" {
			name = parts[1]
			if iface != nil {
				ifaces[iface.Name] = iface
			}
			iface, ok = ifaces[name]
			if !ok {
				iface = &Interface{Name: name}
			}
			iface.DHCP = (parts[3] == "dhcp")
		} else if parts[0] == "address" {
			iface.Address = parts[1]
		} else if parts[0] == "netmask" {
			iface.Netmask = parts[1]
		} else if parts[0] == "gateway" {
			iface.Gateway = parts[1]
		}
	}
	if iface != nil {
		ifaces[iface.Name] = iface
	}

	return ifaces, nil
}

func getMode(iface *Interface) string {
	if iface.DHCP {
		return "dhcp"
	} else {
		return "static"
	}
}

func SaveInterfaces(ifaces map[string]*Interface) error {
	path := getNetworkPath()
	var lines []string

	for _, iface := range ifaces {
		if iface.Auto {
			lines = append(lines, fmt.Sprintf("auto %s", iface.Name))
		}
		mode := getMode(iface)
		lines = append(lines, fmt.Sprintf("iface %s inet %s", iface.Name, mode))
		if !iface.DHCP {
			lines = append(lines, fmt.Sprintf("\taddress %s", iface.Address))
			lines = append(lines, fmt.Sprintf("\tnetmask %s", iface.Netmask))
			lines = append(lines, fmt.Sprintf("\tgateway %s", iface.Gateway))
		}
	}

	err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0666)
	if err != nil {
		return err
	}

	return nil
}
