package futaba

import (
	"errors"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func createCgroup(container_id string) error {
	cmd := exec.Command("cgcreate", "-g", "memory,cpu:"+container_id)
	if err := cmd.Run(); err != nil {
		return errors.New("failed to create cgroup:memory,cpu for " + container_id)
	}
	return nil
}

func setCpuLimit(container_id string, limit string) error {
	// calculate limit
	limit_int, _ := strconv.Atoi(strings.TrimRight(limit, "%"))
	limit_int = 10000 * limit_int
	limit_value := strconv.Itoa(limit_int)

	cmd1 := exec.Command("echo", limit_value, "1000000")
	cmd2 := exec.Command("sudo", "tee", "/sys/fs/cgroup/"+container_id+"/cpu.max")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r

	if err := cmd1.Start(); err != nil {
		return errors.New("failed to echo limit value")
	}
	if err := cmd2.Start(); err != nil {
		return errors.New("failed to write cpu limit: /sys/fs/cgroup/" + container_id + "/cpu.max")
	}

	cmd1.Wait()
	w.Close()
	cmd2.Wait()

	return nil
}

func setMemoryLimit(container_id string, limit string) error {
	cmd1 := exec.Command("echo", limit)
	cmd2 := exec.Command("sudo", "tee", "/sys/fs/cgroup/"+container_id+"/memory.max")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r

	if err := cmd1.Start(); err != nil {
		return errors.New("failed to echo limit value")
	}
	if err := cmd2.Start(); err != nil {
		return errors.New("failed to write memory limit: /sys/fs/cgroup/" + container_id + "/memory.max")
	}

	cmd1.Wait()
	w.Close()
	cmd2.Wait()

	return nil
}

func setCgourpPid(container_id string, pid int) error {
	cmd1 := exec.Command("echo", strconv.Itoa(pid))
	cmd2 := exec.Command("sudo", "tee", "/sys/fs/cgroup/"+container_id+"/cgroup.procs")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r

	if err := cmd1.Start(); err != nil {
		return errors.New("failed to echo pid")
	}
	if err := cmd2.Start(); err != nil {
		return errors.New("failed to write pid: /sys/fs/cgroup/" + container_id + "/cgroup.procs")
	}

	cmd1.Wait()
	w.Close()
	cmd2.Wait()

	return nil
}

func deleteCgroup(container_id string) error {
	cmd := exec.Command("cgdelete", "memory,cpu:"+container_id)
	if err := cmd.Run(); err != nil {
		return errors.New("failed to delete cgroup:memory,cpu for " + container_id)
	}
	return nil
}
