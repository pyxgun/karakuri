package main

import (
	"karakuri"
	"karakuri_mod"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	// flags
	var (
		flag_image     string
		flag_name      string
		flag_restart   string
		flag_namespace string
		flag_os        string
		flag_repositry string
		flag_port      string
		flag_command   string
		flag_mount     string
		flag_id        string
		flag_it        bool
		flag_rm        bool
		flag_buildpath string
		// mod
		flag_mod_name         string
		flag_mod_ingress_edit bool
	)

	app.Name = "Karakuri"
	app.Usage = "karakuri is cli for handling container"

	app.Commands = []cli.Command{
		// container control
		// create
		{
			Name:  "create",
			Usage: "create container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "image",
					Required:    true,
					Destination: &flag_image,
				},
				cli.StringFlag{
					Name:        "port",
					Value:       "none",
					Destination: &flag_port,
				},
				cli.StringFlag{
					Name:        "mount",
					Value:       "none",
					Destination: &flag_mount,
				},
				cli.StringFlag{
					Name:        "cmd",
					Value:       "none",
					Destination: &flag_command,
				},
				cli.StringFlag{
					Name:        "repo",
					Value:       "public",
					Destination: &flag_repositry,
				},
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
				cli.StringFlag{
					Name:        "restart",
					Value:       "no",
					Destination: &flag_restart,
				},
				cli.StringFlag{
					Name:        "ns",
					Value:       "none",
					Destination: &flag_namespace,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.CreateContainer(karakuri.RequestCreateContainer{
					Image:     flag_image,
					Name:      flag_name,
					Namespace: flag_namespace,
					Port:      flag_port,
					Mount:     flag_mount,
					Cmd:       flag_command,
					Repositry: flag_repositry,
					Restart:   flag_restart,
				})
			},
		},
		// start
		{
			Name:  "start",
			Usage: "start container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_id,
				},
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
				cli.BoolFlag{
					Name:        "it",
					Destination: &flag_it,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.StartContainer(karakuri.RequestStartContainer{
					Id:       flag_id,
					Name:     flag_name,
					Terminal: flag_it,
				})
			},
		},
		// run
		{
			Name:  "run",
			Usage: "run container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
				cli.StringFlag{
					Name:        "image",
					Required:    true,
					Destination: &flag_image,
				},
				cli.StringFlag{
					Name:        "port",
					Value:       "none",
					Destination: &flag_port,
				},
				cli.StringFlag{
					Name:        "mount",
					Value:       "none",
					Destination: &flag_mount,
				},
				cli.StringFlag{
					Name:        "cmd",
					Value:       "none",
					Destination: &flag_command,
				},
				cli.BoolFlag{
					Name:        "it",
					Destination: &flag_it,
				},
				cli.StringFlag{
					Name:        "repo",
					Value:       "public",
					Destination: &flag_repositry,
				},
				cli.StringFlag{
					Name:        "ns",
					Value:       "none",
					Destination: &flag_namespace,
				},
				cli.BoolFlag{
					Name:        "rm",
					Destination: &flag_rm,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.RunContainer(karakuri.RequestRunContainer{
					Name:      flag_name,
					Namespace: flag_namespace,
					Image:     flag_image,
					Port:      flag_port,
					Mount:     flag_mount,
					Terminal:  flag_it,
					Cmd:       flag_command,
					Repositry: flag_repositry,
					Remove:    flag_rm,
				})
			},
		},
		// exec
		{
			Name:  "exec",
			Usage: "exec container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_id,
				},
				cli.BoolFlag{
					Name:        "it",
					Destination: &flag_it,
				},
				cli.StringFlag{
					Name:        "cmd",
					Value:       "",
					Required:    true,
					Destination: &flag_command,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.ExecContainer(karakuri.RequestExecContainer{
					Id:       flag_id,
					Name:     flag_name,
					Terminal: flag_it,
					Cmd:      flag_command,
				})
			},
		},
		// stop
		{
			Name:  "stop",
			Usage: "stop container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_id,
				},
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.StopContainer(karakuri.RequestStopContainer{
					Id:   flag_id,
					Name: flag_name,
				})
			},
		},
		// all stop
		{
			Name:  "stopall",
			Usage: "stop all container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "ns",
					Value:       "all",
					Destination: &flag_namespace,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.StopAllContaier(flag_namespace)
			},
		},
		// restart
		{
			Name:  "restart",
			Usage: "restart container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_id,
				},
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.RestartContainer(karakuri.RequsetRestartContainer{
					Id:       flag_id,
					Name:     flag_name,
					Terminal: false,
				})
			},
		},
		// ls
		{
			Name:  "ls",
			Usage: "show container list",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "ns",
					Value:       "none",
					Destination: &flag_namespace,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.ShowContainerList(flag_namespace)
			},
		},
		// rm
		{
			Name:  "rm",
			Usage: "delete container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_id,
				},
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.DeleteContainer(karakuri.RequestDeleteContainer{
					Id:   flag_id,
					Name: flag_name,
				})
			},
		},
		// spec
		{
			Name:  "spec",
			Usage: "show container spec",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_id,
				},
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.ShowContainerSpec(karakuri.RequestShowContainerSpec{
					Id:   flag_id,
					Name: flag_name,
				})
			},
		},
		// logs
		{
			Name:  "logs",
			Usage: "show container logs",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Destination: &flag_id,
				},
				cli.StringFlag{
					Name:        "name",
					Value:       "none",
					Destination: &flag_name,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.ShowContainerLog(flag_id, flag_name)
			},
		},

		// image control
		// images
		{
			Name:  "images",
			Usage: "show images",
			Action: func(c *cli.Context) {
				karakuri.ShowImage()
			},
		},
		// pull
		{
			Name:  "pull",
			Usage: "pull image",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "image",
					Required:    true,
					Destination: &flag_image,
				},
				cli.StringFlag{
					Name:        "os",
					Value:       "linux:amd64",
					Destination: &flag_os,
				},
				cli.StringFlag{
					Name:        "repo",
					Value:       "public",
					Destination: &flag_repositry,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.PullImage(flag_image, flag_os, flag_repositry)
			},
		},
		// rmi
		{
			Name:  "rmi",
			Usage: "delete image",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "id",
					Required:    true,
					Destination: &flag_id,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.DeleteImage(flag_id)
			},
		},
		// build
		{
			Name:  "build",
			Usage: "build image",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Required:    true,
					Destination: &flag_name,
				},
				cli.StringFlag{
					Name:        "buildpath",
					Value:       ".",
					Destination: &flag_buildpath,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.BuildImage(flag_name, flag_buildpath)
			},
		},

		// namespace control
		// show namespace
		{
			Name:  "ns",
			Usage: "show namespaces",
			Action: func(c *cli.Context) {
				karakuri.ShowNamespace()
			},
		},
		// create namespace
		{
			Name:  "createns",
			Usage: "create namespace",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Required:    true,
					Destination: &flag_namespace,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.CreateNamespace(flag_namespace)
			},
		},
		// delete namespace
		{
			Name:  "rmns",
			Usage: "delete namespace",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Required:    true,
					Destination: &flag_namespace,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.DeleteNamespace(flag_namespace)
			},
		},

		// module
		// enable
		{
			Name:  "mod",
			Usage: "show module",
			Action: func(c *cli.Context) {
				karakuri.ShowModuleList()
			},
			Subcommands: []cli.Command{
				{
					Name:  "enable",
					Usage: "enable karakuri module",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "name",
							Required:    true,
							Destination: &flag_mod_name,
						},
					},
					Action: func(c *cli.Context) {
						karakuri.EnableModule(flag_mod_name)
					},
				},
				{
					Name:  "disable",
					Usage: "disable karakuri module",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "name",
							Required:    true,
							Destination: &flag_mod_name,
						},
					},
					Action: func(c *cli.Context) {
						karakuri.DisableModule(flag_mod_name)
					},
				},
				{
					Name:  "ingress",
					Usage: "edit ingress condition",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:        "edit",
							Destination: &flag_mod_ingress_edit,
						},
					},
					Action: func(c *cli.Context) {
						if flag_mod_ingress_edit {
							karakuri_mod.EditIngressEntry()
						}
					},
				},
			},
		},

		// other
		{
			Name:  "version",
			Usage: "show karakuri version",
			Action: func(c *cli.Context) {
				karakuri.PrintKarakuriVersion()
			},
		},
	}

	app.Run(os.Args)
}
