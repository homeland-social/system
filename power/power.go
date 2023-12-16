package power

import (
	"os/exec"
)

func Shutdown() error {
	cmd := exec.Command("/sbin/halt")
	return cmd.Run()
}

func Restart() error {
	cmd := exec.Command("/sbin/reboot")
	return cmd.Run()
}
