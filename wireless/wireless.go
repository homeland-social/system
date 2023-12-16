package wireless

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"system/logging"
	"system/network"
	"system/service"
)

const DEFAULT_HOSTAPD_CONTROL_PATH = "root/run/hostapd.sock"
const DEFAULT_HOSTAPD_CONF_PATH = "root/etc/hostapd/hostapd.conf"
const DEFAULT_WPA_CONF_PATH = "root/etc/wpa_supplicant/wpa_supplicant.conf"

const (
	AccessPoint string = "AP"
	Client      string = "CLIENT"
)

type WirelessConfig struct {
	Mode       string             `json:"mode"`
	SSID       string             `json:"ssid"`
	Passphrase string             `json:"passphrase"`
	Channel    int                `json:"channel"`
	Interface  *network.Interface `json:"interface"`
}

var wirelessConfig *WirelessConfig = nil

func IsValidMode(mode string) bool {
	if mode == "AP" {
		return true
	} else if mode == "CLIENT" {
		return true
	} else {
		return false
	}
}

func (s *WirelessConfig) UnmarshalJSON(data []byte) error {
	// Perform validation for WirelessConfig attributes
	type Temp WirelessConfig
	var a *Temp = (*Temp)(s)
	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	if !IsValidMode(s.Mode) {
		return errors.New("Invalid value for Mode")
	}

	return nil
}

func getHostAPDControlPath() string {
	path, exists := os.LookupEnv("HOSTAPD_CONTROL_PATH")

	if !exists {
		path = DEFAULT_HOSTAPD_CONTROL_PATH
	}

	return path
}

func getHostAPDConfPath() string {
	path, exists := os.LookupEnv("HOSTAPD_CONF_PATH")

	if !exists {
		path = DEFAULT_HOSTAPD_CONF_PATH
	}

	return path
}

func getWpaSupplicantConfPath() string {
	path, exists := os.LookupEnv("WPA_CONF_PATH")

	if !exists {
		path = DEFAULT_WPA_CONF_PATH
	}

	return path
}

func ReadHostAPDConf(wifi *WirelessConfig) error {
	body, err := os.ReadFile(getHostAPDConfPath())
	if err != nil {
		return err
	}
	ifaces, err := network.Load()
	if err != nil {
		return err
	}

	lines := strings.Split(string(body[:]), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "=")
		if parts[0] == "ssid" {
			wifi.SSID = parts[1]
		} else if parts[0] == "interface" {
			wifi.Interface = ifaces.Get(parts[1])
		} else if parts[0] == "channel" {
			wifi.Channel, err = strconv.Atoi(parts[1])
			if err != nil {
				logging.Error("Error converting channel %s to in: %s", parts[1], err)
			}
		} else if parts[0] == "wpa_passphrase" {
			wifi.Passphrase = strings.Repeat("*", len(parts[1]))
		}
	}

	return nil
}

func Load() (*WirelessConfig, error) {
	wifi := &WirelessConfig{}

	if service.Check("hostapd") {
		wifi.Mode = "AP"
		err := ReadHostAPDConf(wifi)
		if err != nil {
			return nil, err
		}
	} else if service.Check("wpa_supplicant") {
		wifi.Mode = "CLIENT"
	}

	return wifi, nil
}
