package futaba

import (
	"fmt"
	"karakuripkgs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func setHostname(spec string) {
	// read specfile
	config_spec := karakuripkgs.ReadSpecFile(spec)

	if err := syscall.Sethostname([]byte(config_spec.Hostname)); err != nil {
		panic(err)
	}
}

func setNameeserver(config_spec karakuripkgs.ConfigSpec) {
	resolv_conf := config_spec.Root.Path + "/merged/etc/resolv.conf"
	if err := os.WriteFile(resolv_conf, []byte("nameserver "+config_spec.Network.Nameserver), 0644); err != nil {
		return
	}
}

func setEnv(config_spec karakuripkgs.ConfigSpec) {
	envs := config_spec.Process.Env
	for _, entry := range envs {
		if err := os.Setenv(entry.Key, entry.Value); err != nil {
			panic(err)
		}
	}
}

func mountOverlay(config_spec karakuripkgs.ConfigSpec) {
	// overlay
	if err := syscall.Mount(
		"overlay",
		config_spec.Root.Path+"/merged",
		"overlay",
		0,
		"lowerdir="+config_spec.Image.Path+",upperdir="+config_spec.Root.Path+"/diff,workdir="+config_spec.Root.Path+"/work",
	); err != nil {
		panic(err)
	}
}

func mountFs(config_spec karakuripkgs.ConfigSpec) {
	// mount file systems
	for _, mount_info := range config_spec.Mounts {
		var (
			mount_flag   int
			mount_option string = ""
		)
		// retrieve options
		if mount_info.Options != nil {
			var option_tmp = ""
			for _, option := range mount_info.Options {
				switch option {
				case "nosuid":
					mount_flag |= syscall.MS_NOSUID
				case "noexec":
					mount_flag |= syscall.MS_NOEXEC
				case "nodev":
					mount_flag |= syscall.MS_NODEV
				case "ro":
					mount_flag |= syscall.MS_RDONLY
				case "rw":
					mount_flag |= syscall.O_RDWR
				case "bind":
					mount_flag |= syscall.MS_BIND
				default:
					option_tmp += option + ","
				}
			}
			mount_option = strings.TrimRight(option_tmp, ",")
		} else {
			mount_flag = 0
		}
		// check directory
		if _, err := os.Stat(config_spec.Root.Path + "/merged" + mount_info.Destination); err != nil {
			if err := os.MkdirAll(config_spec.Root.Path+"/merged"+mount_info.Destination, os.ModePerm); err != nil {
				panic(err)
			}
		}
		// mount
		if err := syscall.Mount(mount_info.Source, config_spec.Root.Path+"/merged"+mount_info.Destination, mount_info.MountType, uintptr(mount_flag), mount_option); err != nil {
			fmt.Printf("[ERROR] Failed to mount %s to %s, type: %s\n", mount_info.Source, mount_info.Destination, mount_info.MountType)
			os.Exit(1)
		}
	}
}

func unmountFs() {
	// unmount proc
	syscall.Unmount("/proc", 0)
}

func pivotRoot(container_dir string) {
	// change direrctory to container directory
	if err := os.Chdir(container_dir); err != nil {
		panic(err)
	}
	// mount merged
	if err := syscall.Mount("merged", container_dir+"/merged", "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		panic(err)
	}
	// create put_old directory
	if err := os.MkdirAll(container_dir+"/merged/put_old", 0700); err != nil {
		panic(err)
	}
	// pivot_root
	if err := syscall.PivotRoot("merged", container_dir+"/merged/put_old"); err != nil {
		panic(err)
	}
	// change directory to root after pivot_root
	if err := os.Chdir("/"); err != nil {
		panic(err)
	}
	// unmount put_old
	if err := syscall.Unmount("/put_old", syscall.MNT_DETACH); err != nil {
		panic(err)
	}
	// delete put_old
	if err := syscall.Rmdir("/put_old"); err != nil {
		panic(err)
	}
}

func createFifo(fifo_path string) {
	if err := unix.Mkfifo(fifo_path+"/exec.fifo", 0o622); err != nil {
		panic(err)
	}
}

func openFifo(fifo_path string) {
	_, err := unix.Open(fifo_path+"/exec.fifo", unix.O_WRONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		panic(err)
	}
}

func waitParant(fifo_path string) {
	openFifo(fifo_path)
}

func startChild(fifo_path string) {
	_, err := os.OpenFile(fifo_path+"/exec.fifo", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
}

func StartContainer(id string, terminal bool) {
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		panic(err)
	}

	// config spec path
	container_spec_path := karakuripkgs.FUTABA_ROOT + "/" + id

	args := []string{"init", "--spec=" + container_spec_path}

	cmd := exec.Command(initCmd, args[0:]...)
	// set standard i/o/e
	if terminal {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		//devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
		logfile, err := os.OpenFile(container_spec_path+"/container.log", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = logfile
		cmd.Stderr = logfile
	}

	// set clone flags
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWCGROUP |
			syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        65535,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        65535,
			},
		},
	}

	// execute command
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	// retrieve pid
	pid := cmd.Process.Pid
	// update pid
	karakuripkgs.UpdateSpecPid(container_spec_path, pid)

	// retrieve config spec
	config_spec := karakuripkgs.ReadSpecFile(container_spec_path)

	// setup network interface
	setupContainerNetwork(pid, config_spec.Network)

	// retrieve resource limit
	cpu_max := config_spec.Cgroup.Cpu.Max
	mem_max := config_spec.Cgroup.Memory.Max
	// create cgroup
	createCgroup(id)
	// set cpu limit
	setCpuLimit(id, cpu_max)
	// set memory limit
	setMemoryLimit(id, mem_max)
	// set pid to cgroup
	setCgourpPid(id, pid)

	// start child
	startChild(config_spec.Fifo)

	if terminal {
		cmd.Wait()
	}
}

func execEntrypoint(config_spec karakuripkgs.ConfigSpec) {
	args := config_spec.Process.Args
	cmd := exec.Command(args[0], args[1:]...)

	// set standard i/o/e
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}

func InitContainer(spec string) {
	// read specfile
	config_spec := karakuripkgs.ReadSpecFile(spec)

	// wait parent process
	waitParant(config_spec.Fifo)

	// set hostname
	setHostname(spec)
	// mount overlay
	mountOverlay(config_spec)
	// mount system
	mountFs(config_spec)
	// env
	setEnv(config_spec)
	// set nameserver
	setNameeserver(config_spec)
	// pivot root
	pivotRoot(config_spec.Root.Path)

	// execute entry point
	execEntrypoint(config_spec)

	// unmount when container exit
	unmountFs()
}

func ExecContainer(id string, terminal bool, command string) {
	// config spec path
	container_spec_path := karakuripkgs.FUTABA_ROOT + "/" + id
	// read config spec
	config_spec := karakuripkgs.ReadSpecFile(container_spec_path)
	// command
	cmd := strings.Split(command, ",")

	pid := config_spec.Process.Pid
	args := []string{"-t", strconv.Itoa(pid), "--all"}
	args = append(args, cmd...)
	nsenter := exec.Command("nsenter", args...)

	nsenter.Stdin = os.Stdin
	nsenter.Stdout = os.Stdout
	nsenter.Stderr = os.Stderr

	if err := nsenter.Start(); err != nil {
		panic(err)
	}
	if terminal {
		nsenter.Wait()
	}
}

func KillContainer(id string) {
	// config spec path
	container_spec_path := karakuripkgs.FUTABA_ROOT + "/" + id
	// read config spec
	config_spec := karakuripkgs.ReadSpecFile(container_spec_path)

	// kill process
	proc, err := os.FindProcess(config_spec.Process.Pid)
	if err != nil {
		panic(err)
	}
	if err := proc.Kill(); err != nil {
		panic(err)
	}
}
