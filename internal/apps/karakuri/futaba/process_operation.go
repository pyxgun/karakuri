package futaba

import (
	"errors"
	"fmt"
	"karakuripkgs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func setHostname(spec string) error {
	// read specfile
	config_spec := karakuripkgs.ReadSpecFile(spec)

	if err := syscall.Sethostname([]byte(config_spec.Hostname)); err != nil {
		return errors.New("failed to set hostname")
	}
	return nil
}

func setNameeserver(config_spec karakuripkgs.ConfigSpec) error {
	resolv_conf := config_spec.Root.Path + "/merged/etc/resolv.conf"
	if err := os.WriteFile(resolv_conf, []byte("nameserver "+config_spec.Network.Nameserver+"\nsearch karakuri.container\n"), 0644); err != nil {
		return errors.New("failed to create /etc/resolv.conf")
	}
	return nil
}

func setDevNull(config_spec karakuripkgs.ConfigSpec) error {
	//cmd := exec.Command("mknod", config_spec.Root.Path+"/merged/dev/null", "c", "1", "3")
	if err := unix.Mknod(config_spec.Root.Path+"/merged/dev/null", 0666, 0); err != nil {
		return errors.New("failed to create /dev/null")
	}
	return nil
}

func setEnv(config_spec karakuripkgs.ConfigSpec) error {
	envs := config_spec.Process.Env
	for _, entry := range envs {
		if err := os.Setenv(entry.Key, entry.Value); err != nil {
			return errors.New("failed to set environmental variables: " + entry.Key + "=" + entry.Value)
		}
	}
	return nil
}

func mountOverlay(config_spec karakuripkgs.ConfigSpec) error {
	// overlay
	if err := syscall.Mount(
		"overlay",
		config_spec.Root.Path+"/merged",
		"overlay",
		0,
		"lowerdir="+config_spec.Image.Path+",upperdir="+config_spec.Root.Path+"/diff,workdir="+config_spec.Root.Path+"/work",
	); err != nil {
		return errors.New("failed to mount overlay filesystem")
	}
	return nil
}

func mountFs(config_spec karakuripkgs.ConfigSpec) error {
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
				return errors.New("failed to create mount destination directory")
			}
		}
		// mount
		if err := syscall.Mount(mount_info.Source, config_spec.Root.Path+"/merged"+mount_info.Destination, mount_info.MountType, uintptr(mount_flag), mount_option); err != nil {
			return errors.New("failed to mount " + mount_info.Source + " to " + mount_info.Destination + ", type: " + mount_info.MountType)
		}
	}
	return nil
}

func unmountFs() error {
	// unmount proc
	if err := syscall.Unmount("/proc", 0); err != nil {
		return errors.New("failed to unmount /proc")
	}
	return nil
}

func pivotRoot(container_dir string) error {
	// change direrctory to container directory
	if err := os.Chdir(container_dir); err != nil {
		return errors.New("failed to enter container layer directory")
	}
	// mount merged
	if err := syscall.Mount("merged", container_dir+"/merged", "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return errors.New("failed to mount merged filesystem")
	}
	// create put_old directory
	if err := os.MkdirAll(container_dir+"/merged/put_old", 0700); err != nil {
		return errors.New("failed to create /put_old")
	}
	// pivot_root
	if err := syscall.PivotRoot("merged", container_dir+"/merged/put_old"); err != nil {
		return errors.New("failed to execute pivot_root")
	}
	// change directory to root after pivot_root
	if err := os.Chdir("/"); err != nil {
		return errors.New("failed to change directory to /")
	}
	// unmount put_old
	if err := syscall.Unmount("/put_old", syscall.MNT_DETACH); err != nil {
		return errors.New("failed to unmount /put_old")
	}
	// delete put_old
	if err := syscall.Rmdir("/put_old"); err != nil {
		return errors.New("failed to remove /put_old")
	}
	return nil
}

func createFifo(fifo_path string) error {
	if err := unix.Mkfifo(fifo_path+"/exec.fifo", 0o622); err != nil {
		return errors.New("failed to create named pipe: " + fifo_path + "/exec.fifo")
	}
	return nil
}

func waitParant(fifo_path string) error {
	_, err := unix.Open(fifo_path+"/exec.fifo", unix.O_WRONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return errors.New("failed to open named pipe: " + fifo_path + "/exec.fifo for waiting parant process")
	}
	return nil
}

func startChild(fifo_path string) error {
	_, err := os.OpenFile(fifo_path+"/exec.fifo", os.O_RDONLY, 0)
	if err != nil {
		return errors.New("failed to open named pipe: " + fifo_path + "/exec.fifo for starting child process")
	}
	return nil
}

func StartContainer(id string, terminal bool) {
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		fmt.Println(err)
		return
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
			fmt.Println(err)
			return
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
				Size:        1,
			},
			{
				ContainerID: 1,
				HostID:      100000,
				Size:        65535,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
			{
				ContainerID: 1,
				HostID:      100000,
				Size:        65535,
			},
		},
	}

	// execute command
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
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
	if err := createCgroup(id); err != nil {
		fmt.Println(err)
		return
	}
	// set cpu limit
	if err := setCpuLimit(id, cpu_max); err != nil {
		fmt.Println(err)
		return
	}
	// set memory limit
	if err := setMemoryLimit(id, mem_max); err != nil {
		fmt.Println(err)
		return
	}
	// set pid to cgroup
	if err := setCgourpPid(id, pid); err != nil {
		fmt.Println(err)
		return
	}

	// start child
	if err := startChild(config_spec.Fifo); err != nil {
		fmt.Println(err)
		return
	}

	if terminal {
		cmd.Wait()
	}
}

func execEntrypoint(config_spec karakuripkgs.ConfigSpec) error {
	args := config_spec.Process.Args
	cmd := exec.Command(args[0], args[1:]...)

	// set standard i/o/e
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.New("failed to execute container entrypoint: " + strings.Join(args, " "))
	}
	cmd.Wait()

	return nil
}

func InitContainer(spec string) {
	// read specfile
	config_spec := karakuripkgs.ReadSpecFile(spec)

	// wait parent process
	if err := waitParant(config_spec.Fifo); err != nil {
		fmt.Println(err)
		return
	}

	// set hostname
	if err := setHostname(spec); err != nil {
		fmt.Println(err)
		return
	}
	// mount overlay
	if err := mountOverlay(config_spec); err != nil {
		fmt.Println(err)
		return
	}
	// mount system
	if err := mountFs(config_spec); err != nil {
		fmt.Println(err)
		return
	}
	// env
	if err := setEnv(config_spec); err != nil {
		fmt.Println(err)
		return
	}
	// set nameserver
	if err := setNameeserver(config_spec); err != nil {
		fmt.Println(err)
		return
	}
	// set /dev/null
	if err := setDevNull(config_spec); err != nil {
		fmt.Println(err)
		return
	}
	// pivot root
	if err := pivotRoot(config_spec.Root.Path); err != nil {
		fmt.Println(err)
		return
	}

	// execute entry point
	if err := execEntrypoint(config_spec); err != nil {
		fmt.Println(err)
		return
	}

	// unmount when container exit
	if err := unmountFs(); err != nil {
		fmt.Println(err)
		return
	}
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
		fmt.Println(err)
		return
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
		fmt.Println(err)
		return
	}
	if err := proc.Kill(); err != nil {
		fmt.Println(err)
		return
	}
}
