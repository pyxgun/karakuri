package futaba

import (
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func createCgroup(container_id string) {
	cmd := exec.Command("cgcreate", "-g", "memory,cpu:"+container_id)
	if err := cmd.Run(); err != nil {
		return
	}
}

func setCpuLimit(container_id string, limit string) {
	// calculate limit
	limit_int, _ := strconv.Atoi(strings.TrimRight(limit, "%"))
	limit_int = 10000 * limit_int
	limit_value := strconv.Itoa(limit_int)

	cmd1 := exec.Command("echo", limit_value, "1000000")
	cmd2 := exec.Command("sudo", "tee", "/sys/fs/cgroup/"+container_id+"/cpu.max")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r

	cmd1.Start()
	cmd2.Start()

	cmd1.Wait()
	w.Close()
	cmd2.Wait()
}

func setMemoryLimit(container_id string, limit string) {
	cmd1 := exec.Command("echo", limit)
	cmd2 := exec.Command("sudo", "tee", "/sys/fs/cgroup/"+container_id+"/memory.max")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r

	cmd1.Start()
	cmd2.Start()

	cmd1.Wait()
	w.Close()
	cmd2.Wait()
}

func setCgourpPid(container_id string, pid int) {
	cmd1 := exec.Command("echo", strconv.Itoa(pid))
	cmd2 := exec.Command("sudo", "tee", "/sys/fs/cgroup/"+container_id+"/cgroup.procs")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r

	cmd1.Start()
	cmd2.Start()

	cmd1.Wait()
	w.Close()
	cmd2.Wait()
}

func deleteCgroup(container_id string) {
	cmd := exec.Command("cgdelete", "memory,cpu:"+container_id)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
