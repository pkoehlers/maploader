package tar

import (
	"log"
	"os/exec"
)

// Use shell command tar to add supplied filePaths to the tar
func Tar(tarballFilePath string, filePaths ...string) error {
	return tarCommand("zcvf", tarballFilePath, filePaths...)
}

// Use shell command tar to untar into the root file system
func Untar(tarballFilePath string, target string) error {
	return tarCommand("xvf", tarballFilePath, "-C", "/")
}

// Private methods

func tarCommand(mode string, tarballFilePath string, filePaths ...string) error {
	app := "tar"
	args := append([]string{mode, tarballFilePath})
	argsFinal := append(args, filePaths...)

	cmd := exec.Command(app, argsFinal...)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	// Print the output
	log.Println(string(stdout))
	return nil
}
