package robot

import (
	"errors"
	"log"
	"maploader/util"
	"os"
	"strings"
)

type Robot struct {
	model            string
	MapFolders       []string
	mapFiles         []string
	restartProcesses []Process
}

func (r *Robot) MapFilesAndFolders() []string {
	return append(r.MapFolders, r.mapFiles...)
}

var robots []Robot = []Robot{
	{model: "p2029", // Dreame L10 Pro
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap"},
		mapFiles:         []string{"/data/config/ava/mult_map.json"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
}

var CurrentRobot *Robot

func DetectRobot() {
	hostname, err := os.Hostname()

	util.CheckAndHandleError(err)

	robotName := strings.Split(hostname, "_")[0]
	var possibleRobots []Robot

	for _, cRobot := range robots {
		if cRobot.model == robotName {
			possibleRobots = append(possibleRobots, cRobot)
		}
	}
	if len(possibleRobots) == 0 {
		possibleRobots = robots
		log.Printf("Could not determine robot model using model: %s", robotName)
		log.Printf("Trying to identify compatible configuration. Please open a Issue/PR, if you think this configuration works well on your model")
	}

	CurrentRobot, err = SelectAndValidateRobot(possibleRobots...)
	util.CheckAndHandleError(err)
	log.Printf("Using configuration for model: %s", CurrentRobot.model)
}

func SelectAndValidateRobot(possibleRobots ...Robot) (*Robot, error) {

	for _, cRobot := range possibleRobots {
		if util.CheckDirectoriesExists(cRobot.MapFolders...) &&
			util.CheckFilesExists(cRobot.mapFiles...) {
			return &cRobot, nil
		}
	}
	return nil, errors.New("Robot model could not be determined or is not supported")
}
