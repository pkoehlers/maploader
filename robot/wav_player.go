package robot

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// PlayMapLoadedSound plays the sound for map loaded event
func PlayMapLoadedSound(mapName string) error {
	wavFilePath := os.Getenv("WAV_FILE_MAP_LOADED")
	if wavFilePath == "" {
		return nil
	}

	// Replace {map_name} placeholder
	if strings.Contains(wavFilePath, "{map_name}") {
		wavFilePath = strings.Replace(wavFilePath, "{map_name}", mapName, -1)
	}

	// Check if the WAV file exists, use default if not
	if _, err := os.Stat(wavFilePath); os.IsNotExist(err) {
		wavFilePath = strings.Replace(wavFilePath, mapName, "default", -1)
	}
	// Check for the existence of the WAV file
	if _, err := os.Stat(wavFilePath); os.IsNotExist(err) {
		log.Printf("WAV file does not exist: %s", wavFilePath)
		return errors.New("WAV file does not exist")
	}

	aplayCmdArgs := os.Getenv("WAV_APLAY_ARGS")
	if aplayCmdArgs == "" {
		aplayCmdArgs = "-Dhw:0,0"
	}
	log.Printf("Calling aplay (wav) with args: %s and file %s", aplayCmdArgs, wavFilePath)

	cmd := exec.Command("aplay", aplayCmdArgs, wavFilePath)
	return cmd.Run()
}
