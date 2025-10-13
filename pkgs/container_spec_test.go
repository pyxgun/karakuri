package karakuripkgs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSpecFile(t *testing.T) {
	spec_flag := SpecFlag{
		RootPath:    "/etc/karakuri/futaba/container",
		ImagePath:   "/etc/karakuri/hitoha/image",
		Command:     "echo,hello",
		OutputPath:  ".",
		HostDevice:  "host_dev",
		Address:     "10.157.0.2/24",
		Gateway:     "10.157.0.1",
		Nameserver:  "8.8.8.8",
		Mount:       "/host/conf:/conf,/host/data:/data",
		PortForward: "8080:80:tcp,5353:53:udp",
		EnvVars:     "PATH=/bin:/usr/local/bin!ME=admin*user",
	}
	CreateSpecFile(spec_flag)
	actual := ReadSpecFile(".")
	os.Remove("./config.json")
	hostname := actual.Hostname

	expected := ConfigSpec{
		Version: KARAKURI_VERSION,
		Process: SpecProcess{
			Args: []string{
				"echo",
				"hello",
			},
			Env: []SpecEnv{
				{
					Key:   "PATH",
					Value: "/bin:/usr/local/bin",
				},
				{
					Key:   "ME",
					Value: "admin user",
				},
			},
			Pid: 0,
		},
		Cgroup: SpecCgroup{
			Path: "/sys/fs/cgroup/" + hostname,
			Cpu: SpecCpu{
				Max: "80%",
			},
			Memory: SpecMemory{
				Max: "1024M",
			},
		},
		Root: SpecRoot{
			Path: "/etc/karakuri/futaba/container",
		},
		Image: SpecImage{
			Path: "/etc/karakuri/hitoha/image",
		},
		Hostname: hostname,
		Fifo:     "/etc/karakuri/futaba/container/fifo",
		Mounts: []SpecMount{
			{
				Destination: "/proc",
				MountType:   "proc",
				Source:      "proc",
				Options: []string{
					"nosuid",
					"noexec",
					"nodev",
				},
			},
			{
				Destination: "/sys",
				MountType:   "sysfs",
				Source:      "sysfs",
				Options: []string{
					"nosuid",
					"noexec",
					"nodev",
					"ro",
				},
			},
			{
				Destination: "/sys/fs/cgroup",
				MountType:   "cgroup2",
				Source:      "cgroup",
			},
			{
				Destination: "/dev",
				MountType:   "tmpfs",
				Source:      "tmpfs",
				Options: []string{
					"nosuid",
					"mode=755",
					"size=65536k",
				},
			},
			{
				Destination: "/dev/pts",
				MountType:   "devpts",
				Source:      "devpts",
				Options: []string{
					"rw",
					"nosuid",
					"noexec",
					"newinstance",
					"mode=620",
					"gid=5",
					"ptmxmode=0666",
				},
			},
			{
				Destination: "/dev/mqueue",
				MountType:   "mqueue",
				Source:      "mqueue",
				Options: []string{
					"rw",
					"nosuid",
					"nodev",
					"noexec",
				},
			},
			{
				Destination: "/dev/shm",
				MountType:   "tmpfs",
				Source:      "shm",
				Options: []string{
					"rw",
					"nosuid",
					"nodev",
					"noexec",
					"mode=1777",
					"size=65536k",
				},
			},
			{
				Destination: "/conf",
				MountType:   "",
				Source:      "/host/conf",
				Options: []string{
					"bind",
				},
			},
			{
				Destination: "/data",
				MountType:   "",
				Source:      "/host/data",
				Options: []string{
					"bind",
				},
			},
		},
		Network: SpecNetwork{
			HostDevice: "host_dev",
			Address:    "10.157.0.2/24",
			Gateway:    "10.157.0.1",
			Nameserver: "8.8.8.8",
			Port: []SpecPortForward{
				{
					HostPort:   8080,
					TargetPort: 80,
					Protocol:   "tcp",
				},
				{
					HostPort:   5353,
					TargetPort: 53,
					Protocol:   "udp",
				},
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestUpdateSpecPid(t *testing.T) {
	spec_flag := SpecFlag{
		RootPath:    "/etc/karakuri/futaba/container",
		ImagePath:   "/etc/karakuri/hitoha/image",
		Command:     "echo,hello",
		OutputPath:  ".",
		HostDevice:  "host_dev",
		Address:     "10.157.0.2/24",
		Gateway:     "10.157.0.1",
		Nameserver:  "8.8.8.8",
		Mount:       "/host/conf:/conf,/host/data:/data",
		PortForward: "8080:80:tcp,5353:53:udp",
		EnvVars:     "PATH=/bin:/usr/local/bin!ME=admin*user",
	}
	CreateSpecFile(spec_flag)
	created_spec_file := ReadSpecFile(".")
	hostname := created_spec_file.Hostname

	UpdateSpecPid(".", 12345)
	actual := ReadSpecFile(".")
	os.Remove("./config.json")

	expected := ConfigSpec{
		Version: KARAKURI_VERSION,
		Process: SpecProcess{
			Args: []string{
				"echo",
				"hello",
			},
			Env: []SpecEnv{
				{
					Key:   "PATH",
					Value: "/bin:/usr/local/bin",
				},
				{
					Key:   "ME",
					Value: "admin user",
				},
			},
			Pid: 12345,
		},
		Cgroup: SpecCgroup{
			Path: "/sys/fs/cgroup/" + hostname,
			Cpu: SpecCpu{
				Max: "80%",
			},
			Memory: SpecMemory{
				Max: "1024M",
			},
		},
		Root: SpecRoot{
			Path: "/etc/karakuri/futaba/container",
		},
		Image: SpecImage{
			Path: "/etc/karakuri/hitoha/image",
		},
		Hostname: hostname,
		Fifo:     "/etc/karakuri/futaba/container/fifo",
		Mounts: []SpecMount{
			{
				Destination: "/proc",
				MountType:   "proc",
				Source:      "proc",
				Options: []string{
					"nosuid",
					"noexec",
					"nodev",
				},
			},
			{
				Destination: "/sys",
				MountType:   "sysfs",
				Source:      "sysfs",
				Options: []string{
					"nosuid",
					"noexec",
					"nodev",
					"ro",
				},
			},
			{
				Destination: "/sys/fs/cgroup",
				MountType:   "cgroup2",
				Source:      "cgroup",
			},
			{
				Destination: "/dev",
				MountType:   "tmpfs",
				Source:      "tmpfs",
				Options: []string{
					"nosuid",
					"mode=755",
					"size=65536k",
				},
			},
			{
				Destination: "/dev/pts",
				MountType:   "devpts",
				Source:      "devpts",
				Options: []string{
					"rw",
					"nosuid",
					"noexec",
					"newinstance",
					"mode=620",
					"gid=5",
					"ptmxmode=0666",
				},
			},
			{
				Destination: "/dev/mqueue",
				MountType:   "mqueue",
				Source:      "mqueue",
				Options: []string{
					"rw",
					"nosuid",
					"nodev",
					"noexec",
				},
			},
			{
				Destination: "/dev/shm",
				MountType:   "tmpfs",
				Source:      "shm",
				Options: []string{
					"rw",
					"nosuid",
					"nodev",
					"noexec",
					"mode=1777",
					"size=65536k",
				},
			},
			{
				Destination: "/conf",
				MountType:   "",
				Source:      "/host/conf",
				Options: []string{
					"bind",
				},
			},
			{
				Destination: "/data",
				MountType:   "",
				Source:      "/host/data",
				Options: []string{
					"bind",
				},
			},
		},
		Network: SpecNetwork{
			HostDevice: "host_dev",
			Address:    "10.157.0.2/24",
			Gateway:    "10.157.0.1",
			Nameserver: "8.8.8.8",
			Port: []SpecPortForward{
				{
					HostPort:   8080,
					TargetPort: 80,
					Protocol:   "tcp",
				},
				{
					HostPort:   5353,
					TargetPort: 53,
					Protocol:   "udp",
				},
			},
		},
	}

	assert.Equal(t, expected, actual)
}
