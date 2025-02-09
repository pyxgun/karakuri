package karakuri

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hitoha"
	"io/fs"
	"karakuripkgs"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/google/uuid"
)

const BUILDFILE = "/Karakurifile"

// build book
type EnvParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BlobParam struct {
	Env []EnvParam `json:"env"`
	Cmd []string   `json:"cmd"`
}

type CopyParam struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type OutContainerParam struct {
	Copy []CopyParam `json:"copy"`
}

type InContainerParam struct {
	Run []string `json:"run"`
}

type BuildBook struct {
	Image        string            `json:"image"`
	InContainer  InContainerParam  `json:"in_container"`
	OutContainer OutContainerParam `json:"out_container"`
	Blob         BlobParam         `json:"blob"`
}

// blobs
type BlobConfig struct {
	Env []string `json:"Env"`
	Cmd []string `json:"Cmd"`
}

type BlobFile struct {
	Config BlobConfig `json:"Config"`
}

func readBuildFile(buildpath string) []string {
	file, err := os.Open(buildpath + BUILDFILE)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_str := scanner.Text()
		if line_str != "" {
			lines = append(lines, line_str)
		}
	}

	return lines
}

func getEnvList(image_tag string) []EnvParam {
	// parse image:tag
	image_tag_info := strings.Split(image_tag, ":")
	image := image_tag_info[0]
	tag := "latest"
	if len(image_tag_info) == 2 {
		tag = image_tag_info[1]
	}

	image_id := hitoha.GetImageId(image, tag)
	blob_path := karakuripkgs.IMAGE_ROOT + "/" + image_id + "/blob.json"

	// read blob
	var bytes []byte
	bytes, err := os.ReadFile(blob_path)
	if err != nil {
		panic(err)
	}

	var blob_file BlobFile
	if err := json.Unmarshal(bytes, &blob_file); err != nil {
		panic(err)
	}

	var env_list []EnvParam
	for _, entry := range blob_file.Config.Env {
		env := strings.Split(entry, "=")
		env_list = append(env_list, EnvParam{Key: env[0], Value: env[1]})
	}

	return env_list
}

func createBuildBook(build_commands []string) BuildBook {
	var build_book BuildBook
	for _, entry := range build_commands {
		command := strings.Split(entry, " ")
		operation := command[0]

		switch operation {
		case "FROM":
			build_book.Image = strings.Join(command[1:], " ")

		case "RUN":
			build_book.InContainer.Run = append(build_book.InContainer.Run, strings.Join(command[1:], " "))

		case "COPY":
			source := command[1]
			destination := command[2]
			build_book.OutContainer.Copy = append(build_book.OutContainer.Copy, CopyParam{Source: source, Destination: destination})

		case "ENV":
			env := strings.Split(command[1], "=")
			env_param := EnvParam{Key: env[0], Value: env[1]}
			build_book.Blob.Env = append(build_book.Blob.Env, env_param)

		case "CMD":
			for _, args := range command[1:] {
				new_arg := strings.Replace(args, "[", "", -1)
				new_arg = strings.Replace(new_arg, "]", "", -1)
				new_arg = strings.Replace(new_arg, "\"", "", -1)
				new_arg = strings.Replace(new_arg, ",", "", -1)
				build_book.Blob.Cmd = append(build_book.Blob.Cmd, new_arg)
			}
		}
	}

	return build_book
}

func copyFile(source string, destination string, is_directory bool) {
	var args []string
	if is_directory {
		args = []string{"cp", "-r", source, destination}
	} else {
		args = []string{"cp", source, destination}
	}
	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func renameFile(source string, destination string) {
	cmd := exec.Command("mv", source, destination)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func mountOverlay(container_layer string, image_path string) {
	if err := syscall.Mount(
		"overlay",
		container_layer+"/merged",
		"overlay",
		0,
		"lowerdir="+image_path+",upperdir="+container_layer+"/diff,workdir="+container_layer+"/work",
	); err != nil {
		panic(err)
	}
}

func unmountOverlay(mount_path string) {
	if err := syscall.Unmount(mount_path, 0); err != nil {
		panic(err)
	}
}

func buildProcCreateBuildBook(buildpath string) BuildBook {
	// read build file
	build_commands := readBuildFile(buildpath)
	// create build book
	build_book := createBuildBook(build_commands)

	return build_book
}

func buildProcCreateScript(run_command []string) {
	command_str := strings.Join(run_command, "\n")
	script := fmt.Sprintf("#!/bin/bash\n\n%s\n", command_str)

	if err := os.WriteFile("./karakuribuild.sh", []byte(script), 0644); err != nil {
		panic(err)
	}
}

func buildProcCreateContainer(build_book BuildBook, image_id string) string {
	container_name := "for_build_image_" + image_id
	container_command := "sh,/karakuribuild.sh"
	CreateContainer(RequestCreateContainer{
		Image:     build_book.Image,
		Name:      container_name,
		Namespace: "system",
		Cmd:       container_command,
		Port:      "none",
		Mount:     "none",
		Repositry: "public",
		Restart:   "no",
	})

	var container_id string
	if res, resp_id, message := karakuripkgs.RequestContainerId(container_name); !res {
		fmt.Println(message)
		return ""
	} else {
		container_id = resp_id
	}

	return container_id
}

func buildProcCopyHostFile(container_id string, build_book BuildBook) {
	container_layer := karakuripkgs.FUTABA_ROOT + "/" + container_id
	copyFile("./karakuribuild.sh", container_layer+"/diff", false)
	// remove karakuribuild.sh
	if err := os.Remove("./karakuribuild.sh"); err != nil {
		panic(err)
	}

	for i, entry := range build_book.OutContainer.Copy {
		source := entry.Source
		destination := entry.Destination
		fmt.Printf("    [%d] Copy \"%s\" to \"%s\"\n", i+1, source, destination)
		copyFile(source, container_layer+"/diff"+destination, true)
	}
}

func buildProcCreateImageLayer(container_id string, image_dir string) {
	container_layer := karakuripkgs.FUTABA_ROOT + "/" + container_id

	// remove karakuribuild.sh from container image
	if err := os.Remove(container_layer + "/diff/karakuribuild.sh"); err != nil {
		panic(err)
	}

	// image path
	spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + container_id)
	image_path := spec.Image.Path
	// re-mount overlay on parant process
	mountOverlay(container_layer, image_path)
	// copy image to new image dir
	if err := os.MkdirAll(image_dir, fs.ModePerm); err != nil {
		panic(err)
	}
	copyFile(container_layer+"/merged", image_dir, true)
	renameFile(image_dir+"/merged", image_dir+"/rootfs")

	// unmount
	unmountOverlay(container_layer + "/merged")
}

func buildProcCreateBlobFile(build_book BuildBook, image_layer string) {
	var blob_file BlobFile
	// set base images's env
	base_env := getEnvList(build_book.Image)
	build_book.Blob.Env = append(build_book.Blob.Env, base_env...)

	for _, entry := range build_book.Blob.Env {
		blob_file.Config.Env = append(blob_file.Config.Env, entry.Key+"="+entry.Value)
	}
	// set cmd
	blob_file.Config.Cmd = build_book.Blob.Cmd

	// write file
	data, _ := json.MarshalIndent(blob_file, "", "  ")
	if err := os.WriteFile(image_layer+"/blob.json", data, fs.ModePerm); err != nil {
		panic(err)
	}
}

func BuildImage(image string, buildpath string) {
	// parse image tag
	image_info := strings.Split(image, ":")
	image_name := image_info[0]
	tag := "latest"
	if len(image_info) == 2 {
		tag = image_info[1]
	}

	// create build book
	fmt.Printf("[1] Create build book... ")
	build_book := buildProcCreateBuildBook(buildpath)
	fmt.Printf("Done.\n")

	// create shell script for command in container
	fmt.Printf("[2] Create shell script for container... ")
	buildProcCreateScript(build_book.InContainer.Run)
	fmt.Printf("Done.\n")

	// create new id
	new_image_id := (uuid.NewString())[24:]
	new_image_dir := karakuripkgs.IMAGE_ROOT + "/" + new_image_id

	// create container for build
	fmt.Printf("[3] Create container for build... ")
	// new image id
	container_id := buildProcCreateContainer(build_book, new_image_id)
	fmt.Printf("Done.\n")

	// COPY: copy local file to container
	fmt.Printf("[4] Copy local files to container...\n")
	buildProcCopyHostFile(container_id, build_book)
	fmt.Printf("Done.\n")

	// start container and execute shell script
	fmt.Printf("[5] Execute RUN command on container...\n")
	StartContainer(RequestStartContainer{
		Id:       container_id,
		Name:     "none",
		Terminal: true,
	})
	fmt.Printf("Done.\n")

	// create image layer
	fmt.Printf("[6] Create image layer... ")
	buildProcCreateImageLayer(container_id, new_image_dir)
	fmt.Printf("Done.\n")

	// create blobs
	fmt.Printf("[7] Create blob file... ")
	buildProcCreateBlobFile(build_book, new_image_dir)
	fmt.Printf("Done.\n")

	// add image list
	fmt.Printf("[8] Add local image list... ")
	hitoha.AddImageList(image_name, tag, new_image_id, new_image_dir+"/rootfs")
	fmt.Printf("Done.\n")

	// delete container
	fmt.Printf("[9] delete containere... ")
	DeleteContainer(RequestDeleteContainer{
		Id:   container_id,
		Name: "none",
	})
	fmt.Println("Build image \"" + image_name + ":" + tag + "\", id: " + new_image_id + " has been completed.")
}
