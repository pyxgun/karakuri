package karakuri_mod

import (
	"karakuripkgs"
	"os"
	"os/exec"
	"strings"
)

func createRegistryBrowserModuleKarakurifile() {
	// create Karakurifile
	device_address, _ := getDeviceIpAddress()
	karakurifile := `FROM klausmeyer/docker-registry-browser

COPY /etc/karakuri/modules/registry-browser/script /script

ENV DOCKER_REGISTRY_URL=http://` + device_address.String() + `:5000
ENV ENABLE_DELETE_IMAGES=true
ENV SECRET_KEY_BASE=karakuri

CMD ["sh", "/script/entrypoint.sh"]
`

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER + "/Karakurifile"); stat != nil {
		fd, err := os.Create(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER + "/Karakurifile")
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		bytes := []byte(karakurifile)
		if _, err := fd.Write(bytes); err != nil {
			panic(err)
		}
	}
}

func createRegistryBrowserModuleConfig() {
	// script
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER + "/script"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER+"/script", os.ModePerm); err != nil {
			panic(err)
		}
	}

	server_conf := `#!/bin/bash

cd /app
/docker-entrypoint.sh web
`
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER + "/script/entrypoint.sh"); stat != nil {
		fd_1, err := os.Create(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER + "/script/entrypoint.sh")
		if err != nil {
			panic(err)
		}
		defer fd_1.Close()
		bytes_1 := []byte(server_conf)
		if _, err := fd_1.Write(bytes_1); err != nil {
			panic(err)
		}
	}
}

func setupRegistryBrowserModule() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// Karakurifile
	createRegistryBrowserModuleKarakurifile()
	// script
	createRegistryBrowserModuleConfig()
}

func buildRegistryBrowserImage(mod_info ModInfo) {
	build_args := []string{
		"build",
		"--name",
		mod_info.ImageName,
		"--buildpath",
		mod_info.Path,
	}
	build := exec.Command("karakuri", build_args...)
	if err := build.Run(); err != nil {
		panic(err)
	}
}

func createRegistryBrowserContainer(mod_info ModInfo) {
	create_args := []string{
		"create",
		"--name", mod_info.Name,
		"--image", mod_info.ImageName,
		"--port", "8081:8080:tcp",
		"--restart", "on-boot",
		"--ns", "system-mod",
	}
	create := exec.Command("karakuri", create_args...)
	if err := create.Run(); err != nil {
		panic(err)
	}
}

func startRegistryBrowserContainer(mod_info ModInfo) {
	start_args := []string{"start", "--name", mod_info.Name}
	start := exec.Command("karakuri", start_args...)
	if err := start.Run(); err != nil {
		panic(err)
	}
}

func enableRegistryBrowserModule(mod_info ModInfo) {
	image_info := strings.Split(mod_info.ImageName, ":")
	if !isImageExists(image_info[0], image_info[1]) {
		buildRegistryBrowserImage(mod_info)
	}
	createRegistryBrowserContainer(mod_info)
	startRegistryBrowserContainer(mod_info)
}

func disableRegistryBrowserModule() {
	removeModuleContainer("registry-browser")
}
