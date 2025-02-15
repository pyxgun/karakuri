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
		flag_registry  string
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
		// registry controller
		flag_regctl_address    string
		flag_regctl_repository string
		flag_regctl_image_tag  string
		// cluster controller
		flag_clsctl_target string
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
					Name:        "registry",
					Value:       "public",
					Destination: &flag_registry,
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
					Registry:  flag_registry,
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
					Name:        "registry",
					Value:       "public",
					Destination: &flag_registry,
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
					Registry:  flag_registry,
					Remove:    flag_rm,
					Restart:   flag_restart,
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
					Name:        "registry",
					Value:       "public",
					Destination: &flag_registry,
				},
			},
			Action: func(c *cli.Context) {
				karakuri.PullImage(flag_image, flag_os, flag_registry)
			},
		},
		// push
		{
			Name:  "push",
			Usage: "push image",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "registry",
					Value:       "connected_registry",
					Destination: &flag_registry,
				},
				cli.StringFlag{
					Name:        "image",
					Required:    true,
					Destination: &flag_image,
				},
			}, Action: func(c *cli.Context) {
				karakuri.PushImage(flag_image, flag_registry)
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

		// registry controller
		{
			Name:  "regctl",
			Usage: "registry controller",
			Subcommands: []cli.Command{
				{
					Name:  "target",
					Usage: "show target registry",
					Action: func(c *cli.Context) {
						karakuri.ShowTargetRegistry()
					},
				},
				{
					Name:  "connect",
					Usage: "connect registry",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "registry",
							Required:    true,
							Destination: &flag_regctl_address,
						},
					},
					Action: func(c *cli.Context) {
						karakuri.ConnectRegistry(flag_regctl_address)
					},
				},
				{
					Name:  "disconnect",
					Usage: "disconnect registry",
					Action: func(c *cli.Context) {
						karakuri.DisconnectRegistry()
					},
				},
				{
					Name: "get",
					Subcommands: []cli.Command{
						{
							Name:  "repository",
							Usage: "get repository",
							Action: func(c *cli.Context) {
								karakuri.ShowRepository()
							},
						},
						{
							Name:  "tag",
							Usage: "get tags",
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:        "repository",
									Required:    true,
									Destination: &flag_regctl_repository,
								},
							},
							Action: func(c *cli.Context) {
								karakuri.ShowTag(flag_regctl_repository)
							},
						},
					},
				},
				{
					Name:  "delete",
					Usage: "delete image",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "image",
							Required:    true,
							Destination: &flag_regctl_image_tag,
						},
					},
					Action: func(c *cli.Context) {
						karakuri.DeleteImageManifest(flag_regctl_image_tag)
					},
				},
			},
		},

		// cluster controller
		{
			Name:  "clsctl",
			Usage: "cluster controller",
			Subcommands: []cli.Command{
				{
					Name:  "info",
					Usage: "show target cluster",
					Action: func(c *cli.Context) {
						karakuri.ShowTargetCluster()
					},
				},
				{
					Name:  "mode",
					Usage: "change cluster mode",
					Subcommands: []cli.Command{
						{
							Name:  "cluster",
							Usage: "change mode to cluster-mode",
							Action: func(c *cli.Context) {
								karakuri.EnableClusterMode()
							},
						},
						{
							Name:  "stand-alone",
							Usage: "change mode to stand-alone-mode",
							Action: func(c *cli.Context) {
								karakuri.DisableClusterMode()
							},
						},
					},
				},
				{
					Name:  "connect",
					Usage: "connect cluster",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "cluster",
							Required:    true,
							Destination: &flag_clsctl_target,
						},
					},
					Action: func(c *cli.Context) {
						karakuri.ConnectCluster(flag_clsctl_target)
					},
				},
				{
					Name:  "disconnect",
					Usage: "disconnect cluster",
					Action: func(c *cli.Context) {
						karakuri.DisconnectCluster()
					},
				},
			},
		},

		// version
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
