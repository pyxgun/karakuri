package main

import (
	"futaba"
	"karakuripkgs"
	"os"

	"github.com/urfave/cli"
)

func main() {
	entrypoint := cli.NewApp()

	// flags
	var (
		spec_flag      karakuripkgs.SpecFlag
		flag_it        bool
		flag_init_mode string
		flag_command   string
		flag_start_id  string
		flag_spec      string
		flag_delete_id string
		flag_kill_id   string
	)

	entrypoint.Name = "Futaba"
	entrypoint.Usage = "low level container runtime"

	entrypoint.Commands = []cli.Command{
		// create
		{
			Name:  "create",
			Usage: "create container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "spec",
					Value:       ".",
					Destination: &flag_spec,
				},
			},
			Action: func(c *cli.Context) {
				futaba.CreateContainer(flag_spec)
			},
		},

		// start
		{
			Name:  "start",
			Usage: "start container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_start_id,
				},
				cli.BoolFlag{
					Name:        "it",
					Destination: &flag_it,
				},
			},
			Action: func(c *cli.Context) {
				futaba.StartContainer(flag_start_id, flag_it)
			},
		},

		// run
		{
			Name:  "run",
			Usage: "run container",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "it",
					Destination: &flag_it,
				},
				cli.StringFlag{
					Name:        "spec",
					Value:       ".",
					Destination: &flag_spec,
				},
			},
			Action: func(c *cli.Context) {
				container_id := futaba.CreateContainer(flag_spec)
				futaba.StartContainer(container_id, flag_it)
			},
		},

		// exec
		{
			Name: "exec",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_start_id,
					Required:    true,
				},
				cli.BoolFlag{
					Name:        "it",
					Destination: &flag_it,
				},
				cli.StringFlag{
					Name:        "cmd",
					Value:       "",
					Destination: &flag_command,
				},
			},
			Action: func(c *cli.Context) {
				futaba.ExecContainer(flag_start_id, flag_it, flag_command)
			},
		},

		// kill
		{
			Name: "kill",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_kill_id,
					Required:    true,
				},
			},
			Action: func(c *cli.Context) {
				futaba.KillContainer(flag_kill_id)
			},
		},

		// spec
		{
			Name:  "spec",
			Usage: "create spec file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "output",
					Value:       ".",
					Destination: &spec_flag.OutputPath,
				},
				cli.StringFlag{
					Name:        "root",
					Destination: &spec_flag.RootPath,
				},
				cli.StringFlag{
					Name:        "cmd",
					Value:       "sh",
					Destination: &spec_flag.Command,
				},
				cli.StringFlag{
					Name:        "image",
					Destination: &spec_flag.ImagePath,
				},
				cli.StringFlag{
					Name:        "hostdevice",
					Value:       "karakuri0",
					Destination: &spec_flag.HostDevice,
				},
				cli.StringFlag{
					Name:        "address",
					Value:       "10.157.0.10/24",
					Destination: &spec_flag.Address,
				},
				cli.StringFlag{
					Name:        "gateway",
					Value:       "10.157.0.1",
					Destination: &spec_flag.Gateway,
				},
				cli.StringFlag{
					Name:        "nameserver",
					Value:       "8.8.8.8",
					Destination: &spec_flag.Nameserver,
				},
				cli.StringFlag{
					Name:        "mount",
					Destination: &spec_flag.Mount,
				},
				cli.StringFlag{
					Name:        "port",
					Destination: &spec_flag.PortForward,
				},
				cli.StringFlag{
					Name:        "env",
					Destination: &spec_flag.EnvVars,
				},
			},
			Action: func(c *cli.Context) {
				karakuripkgs.CreateSpecFile(spec_flag)
			},
		},

		// delete
		{
			Name: "delete",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_delete_id,
					Required:    true,
				},
			},
			Action: func(c *cli.Context) {
				futaba.DeleteContainer(flag_delete_id)
			},
		},

		// list
		{
			Name: "list",
			Action: func(c *cli.Context) {
				futaba.ShowContainerList()
			},
		},

		// hidden: init
		{
			Name: "init",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "mode",
					Destination: &flag_init_mode,
				},
				cli.StringFlag{
					Name:        "spec",
					Value:       ".",
					Destination: &flag_spec,
				},
			},
			Action: func(c *cli.Context) {
				futaba.InitContainer(flag_spec)
			},
			Hidden: true,
		},
	}

	// start
	entrypoint.Run(os.Args)
}
