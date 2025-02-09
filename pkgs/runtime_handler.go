package karakuripkgs

import (
	"os"
	"os/exec"
)

type ParamsRuntimeSpec struct {
	ImagePath  string
	Port       string
	Mount      string
	HostDevice string
	Address    string
	Gateway    string
	Nameserver string
	Command    string
	EnvVars    string
	Restart    string
}

// runtime: spec command
func RuntimeSpec(params ParamsRuntimeSpec) {
	args := []string{
		"spec",
		"--image=" + params.ImagePath,
		"--hostdevice=" + params.HostDevice,
		"--address=" + params.Address,
		"--gateway=" + params.Gateway,
		"--restart=" + params.Restart,
	}
	// set nameserver
	if params.Nameserver != "none" {
		args = append(args, "--nameserver="+params.Nameserver)
	}
	// set port
	if params.Port != "none" {
		args = append(args, "--port="+params.Port)
	}
	// set mount
	if params.Mount != "none" {
		args = append(args, "--mount="+params.Mount)
	}
	// set command
	if params.Command != "none" {
		args = append(args, "--cmd="+params.Command)
	}
	// set env
	if params.EnvVars != "none" {
		args = append(args, "--env="+params.EnvVars)
	}

	cmd := exec.Command(RUNTIME, args[0:]...)

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}

// runtime: create command
func RuntimeCreate() error {
	args := []string{"create"}
	cmd := exec.Command(RUNTIME, args[0:]...)

	if err := cmd.Start(); err != nil {
		return (err)
	}
	cmd.Wait()

	return nil
}

// runtime: start command
func RuntimeStart(id string, terminal bool) {
	args := []string{"start", "--id=" + id}
	if terminal {
		args = append(args, "--it")
	}
	cmd := exec.Command(RUNTIME, args[0:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}

// runtime: run command
func RuntimeRun(terminal bool) {
	args := []string{"run"}
	if terminal {
		args = append(args, "--it")
	}
	cmd := exec.Command(RUNTIME, args[0:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}

// runtime: exec command
func RuntimeExec(id string, terminal bool, cmd_args string) {
	args := []string{"exec", "--id=" + id}
	if terminal {
		args = append(args, "--it")
	}
	if cmd_args != "none" {
		args = append(args, "--cmd="+cmd_args)
	}
	cmd := exec.Command(RUNTIME, args[0:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}

// runtime: kill command
func RuntimeKill(id string) {
	args := []string{"kill", "--id=" + id}
	cmd := exec.Command(RUNTIME, args[0:]...)

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}

// runtime: delete commang
func RuntimeDelete(id string) {
	args := []string{"delete", "--id=" + id}
	cmd := exec.Command(RUNTIME, args[0:]...)

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}
