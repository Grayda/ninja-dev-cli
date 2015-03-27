package main

import (
	"log"
	"os"
	"os/exec"
)

const developerCrossPackageName = "ninjasphere/ninjasphere-developer-cross-base"
const developerPackageName = "ninjasphere/ninjasphere-developer"

func verifyDockerRunning() {
	if !isPackageInstalled("docker") {
		log.Fatalln("Could not find a local installation of 'docker', is it installed and in your PATH?")
	}

	// verify that docker works (that docker daemon is running, etc)
	err := exec.Command("docker", "version").Run()
	if err != nil {
		// doesn't work, so give some guidance
		if isPackageInstalled("boot2docker") {
			log.Fatalln("Docker is installed but not working, did you bring up boot2docker and set it up using $(boot2docker shellinit)?")
		} else {
			log.Fatalln("Docker is installed but not working, is the docker daemon installed and running?")
		}
	}
}

func verifyDockerImages() {
	verifyDockerRunning()

	ensureDockerImage(developerCrossPackageName)
	ensureDockerImage(developerPackageName)
}

func ensureDockerImage(image string) {
	if exec.Command("docker", "inspect", image).Run() != nil {
		log.Printf("Docker package '%s' does not exist, pulling...\n", image)

		cmd := exec.Command("docker", "pull", image)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if cmd.Run() != nil {
			log.Fatalln("Could not retrieve docker package, see above log output for details.")
		}
	}
}

func isPackageInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return (err == nil)
}

func runNativeDockerCrossCommand(dockerArgs []string, command ...string) error {
	args := []string{"run", "--rm", "--privileged", "-it"}
	args = append(args, dockerArgs...)
	args = append(args, developerCrossPackageName)
	args = append(args, command...)
	cmd := exec.Command("docker", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runNativeDockerCommand(dockerArgs []string, command ...string) error {
	args := []string{"run", "--rm", "-it"}
	args = append(args, dockerArgs...)
	args = append(args, developerPackageName)
	args = append(args, command...)
	cmd := exec.Command("docker", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
