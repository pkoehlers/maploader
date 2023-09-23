package robot

import (
	"maploader/config"
	"maploader/util"
	"os"
	"os/exec"
)

type Process struct {
	StartCommand string
	StartArgs    []string
	StopCommand  string
	StopArgs     []string
}

var MiioClientProcess = Process{StartCommand: "sh", StartArgs: []string{"/etc/rc.d/miio.sh"}, StopCommand: "sh", StopArgs: []string{"/etc/rc.d/miio.sh", "stop"}}
var AvaProcess = Process{StartCommand: "sh", StartArgs: []string{"/etc/rc.d/ava.sh"}, StopCommand: "killall", StopArgs: []string{"-9", "ava"}}

func StopProcesses() {
	if os.Getenv("MAPLOADER_RESTART_VALETUDO") != "" {
		stopValetudo()
	}

	for _, restartProcess := range CurrentRobot.restartProcesses {
		ExcuteCmd(restartProcess.StopCommand, restartProcess.StopArgs...)
	}
}

func StartProcesses() {
	for _, restartProcess := range CurrentRobot.restartProcesses {
		ExcuteCmd(restartProcess.StartCommand, restartProcess.StartArgs...)
	}

	if os.Getenv("MAPLOADER_RESTART_VALETUDO") != "" {
		startValetudo()
	}
}

func ExcuteCmd(cmdStr string, cmdArgs ...string) {

	cmd := exec.Command(cmdStr, cmdArgs...)
	err := cmd.Run()

	if err != nil {
		util.CheckAndHandleError(err)
	}
}

func stopValetudo() {
	ExcuteCmd("killall", "-9", "valetudo")
}

func startValetudo() {
	devnull, dnerr := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if dnerr != nil {
		util.CheckAndHandleError(dnerr)
	}

	cmd := exec.Command("/data/valetudo")
	cmd.Stdout = devnull
	cmd.Env = os.Environ()
	valetudoConfigEnv := "VALETUDO_CONFIG_PATH=" + config.Getenv("VALETUDO_CONFIG_PATH", "/data/valetudo_config.json")
	cmd.Env = append(cmd.Env, valetudoConfigEnv)
	err := cmd.Start()

	if err != nil {
		util.CheckAndHandleError(err)
	} else {
		cmd.Process.Release()
	}
}
