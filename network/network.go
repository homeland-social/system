package network

import (
	"fmt"
	"os"
	"os/exec"
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
const IFUP = "/usr/sbin/ifup"
const IFDOWN = "/usr/sbin/ifdown"

func getNetworkPath() string {
	path, exists := os.LookupEnv("NETWORK_INTERFACES_PATH")

	if !exists {
		path = DEFAULT_NETWORK_INTERFACES
	}

	return path
}

type Interface struct {
	Name    string `json:"name"`
	Auto    bool   `json:"auto"`
	DHCP    bool   `json:"dhcp"`
	Address string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}
type Interfaces map[string]*Interface

func (self *Interface) GetModeString() string {
	if self.DHCP {
		return "dhcp"
	} else {
		return "static"
	}
}

func (self *Interface) Up() error {
	s := fmt.Sprintf("%s %s", IFUP, self.Name)
	cmd := exec.Command(s)
	return cmd.Run()
}

func (self *Interface) Down() error {
	s := fmt.Sprintf("%s %s", IFDOWN, self.Name)
	cmd := exec.Command(s)
	return cmd.Run()
}

func (self *Interface) Reload() error {
	err := self.Down()
	if err != nil {
		return err
	}
	return self.Up()
}

func (self Interfaces) Save() error {
	var lines []string

	for _, iface := range self {
		if iface.Auto {
			lines = append(lines, fmt.Sprintf("auto %s", iface.Name))
		}
		lines = append(lines, fmt.Sprintf("iface %s inet %s", iface.Name, iface.GetModeString()))
		if !iface.DHCP {
			lines = append(lines, fmt.Sprintf("\taddress %s", iface.Address))
			lines = append(lines, fmt.Sprintf("\tnetmask %s", iface.Netmask))
			lines = append(lines, fmt.Sprintf("\tgateway %s", iface.Gateway))
		}
	}

	err := os.WriteFile(getNetworkPath(), []byte(strings.Join(lines, "\n")), 0666)
	if err != nil {
		return err
	}

	return nil
}

func (self Interfaces) Create(name string) *Interface {
	iface := &Interface{Name: name}
	self[name] = iface
	return iface
}

func (self Interfaces) Add(iface *Interface) {
	self[iface.Name] = iface
}

func (self Interfaces) Get(name string) *Interface {
	iface, ok := self[name]
	if !ok {
		return nil
	}
	return iface
}

func (self Interfaces) GetOrAdd(name string) *Interface {
	iface := self.Get(name)
	if iface == nil {
		iface = &Interface{Name: name}
		self.Add(iface)
	}
	return iface
}

func (self Interfaces) Has(name string) bool {
	_, ok := self[name]
	return ok
}

func (self Interfaces) Remove(name string) {
	delete(self, name)
}

func (self Interfaces) Values() []Interface {
	iface_list := make([]Interface, 0, len(self))
	for _, v := range self {
		iface_list = append(iface_list, *v)
	}
	return iface_list
}

func FromList(iface_list []*Interface) *Interfaces {
	ifaces := &Interfaces{}
	for _, v := range iface_list {
		ifaces.Add(v)
	}
	return ifaces
}

func Load() (*Interfaces, error) {
	ifaces := &Interfaces{}
	f, err := os.Open("/sys/class/net")
	if err != nil {
		return nil, err
	}
	files, err := f.ReadDir(0)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		name := v.Name()
		if name == "docker0" {
			continue
		}
		ifaces.Create(v.Name())
	}

	body, err := os.ReadFile(getNetworkPath())
	if err != nil {
		return nil, err
	}

	var iface *Interface = nil
	var name string

	lines := strings.Split(string(body[:]), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")

		if parts[0] == "auto" {
			name = parts[1]
			iface = ifaces.GetOrAdd(name)
			iface.Auto = true
			iface = nil
		} else if parts[0] == "iface" {
			name = parts[1]
			iface = ifaces.GetOrAdd(name)
			iface.DHCP = (parts[3] == "dhcp")
		} else if iface != nil {
			if parts[0] == "address" {
				iface.Address = parts[1]
			} else if iface != nil && parts[0] == "netmask" {
				iface.Netmask = parts[1]
			} else if parts[0] == "gateway" {
				iface.Gateway = parts[1]
			}
		}
	}

	return ifaces, nil
}
