package service

import (
	"os/exec"

	"system/logging"
)

const RC_SERVICE = "/sbin/rc-service"

func Start(name string) error {
	cmd := exec.Command(RC_SERVICE, name, "start")
	return cmd.Run()
}

func Stop(name string) error {
	cmd := exec.Command(RC_SERVICE, name, "stop")
	return cmd.Run()
}

func Restart(name string) error {
	err := Stop(name)
	if err != nil {
		return err
	}
	return Start(name)
}

func Check(name string) bool {
	cmd := exec.Command(RC_SERVICE, name, "status")
	err := cmd.Run()
	if err != nil {
		logging.Error("Command error: %s", err)
	}
	return err == nil
}
