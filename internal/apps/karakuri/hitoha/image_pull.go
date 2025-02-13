package hitoha

import (
	"encoding/json"
	"io"
	"karakuripkgs"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// -------------------------------------
// public repositry
func getToken(image string) string {
	image_info := strings.Split(image, "/")
	var (
		repository string
		image_name string
	)
	if len(image_info) == 2 {
		repository = image_info[0]
		image_name = image_info[1]
	} else {
		repository = "library"
		image_name = image
	}
	// url := "https://auth.docker.io/token?scope=repository:library/" + image + ":pull&service=registry.docker.io"
	url := "https://auth.docker.io/token?scope=repository:" + repository + "/" + image_name + ":pull&service=registry.docker.io"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	var response RespToken
	if err := json.Unmarshal(byte_array, &response); err != nil {
		panic(err)
	}

	return response.Token
}

func getManifestList(image string, tag string, token string) ManifestList {
	image_info := strings.Split(image, "/")
	var (
		repository string
		image_name string
	)
	if len(image_info) == 2 {
		repository = image_info[0]
		image_name = image_info[1]
	} else {
		repository = "library"
		image_name = image
	}
	url := "https://registry-1.docker.io/v2/" + repository + "/" + image_name + "/manifests/" + tag

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Add("Authorization", "Bearer "+token)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	var manifest_list ManifestList
	if err := json.Unmarshal(byte_array, &manifest_list); err != nil {
		panic(err)
	}

	return manifest_list
}

func getBlob(manifest_digest string, image string, token string) ManifestBlob {
	image_info := strings.Split(image, "/")
	var (
		repository string
		image_name string
	)
	if len(image_info) == 2 {
		repository = image_info[0]
		image_name = image_info[1]
	} else {
		repository = "library"
		image_name = image
	}
	url := "https://registry-1.docker.io/v2/" + repository + "/" + image_name + "/blobs/" + manifest_digest

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	var manifest_blob ManifestBlob
	if err := json.Unmarshal(byte_array, &manifest_blob); err != nil {
		panic(err)
	}

	return manifest_blob
}

func getConfig(digest string, image string, token string, image_id string) {
	image_info := strings.Split(image, "/")
	var (
		repository string
		image_name string
	)
	if len(image_info) == 2 {
		repository = image_info[0]
		image_name = image_info[1]
	} else {
		repository = "library"
		image_name = image
	}
	url := "https://registry-1.docker.io/v2/" + repository + "/" + image_name + "/blobs/" + digest

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	// store file
	file, err := os.OpenFile(karakuripkgs.IMAGE_ROOT+"/"+image_id+"/config.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(byte_array)
}

func getLayer(digest string, image string, token string, image_id string, index int) {
	image_info := strings.Split(image, "/")
	var (
		repository string
		image_name string
	)
	if len(image_info) == 2 {
		repository = image_info[0]
		image_name = image_info[1]
	} else {
		repository = "library"
		image_name = image
	}
	url := "https://registry-1.docker.io/v2/" + repository + "/" + image_name + "/blobs/" + digest

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	// store file
	filepath := karakuripkgs.IMAGE_ROOT + "/" + image_id + "/"
	filename := "layer" + strconv.Itoa(index) + ".tar.gz"
	destination := filepath + filename
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(byte_array)
	// extract layer
	cmd := exec.Command("tar", "zxvf", destination, "-C", filepath+"rootfs")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	// remove tar file
	if err := os.Remove(destination); err != nil {
		panic(err)
	}
}

// -------------------------------------

// -------------------------------------
// private repository
func privGetManifest(repositry string, image string, tag string) PrivManifest {
	url := "http://" + repositry + "/v2/" + image + "/manifests/" + tag

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	var manifest PrivManifest
	if err := json.Unmarshal(byte_array, &manifest); err != nil {
		panic(err)
	}

	return manifest
}

func privGetConfig(repositry string, digest string, image string, image_id string) {
	url := "http://" + repositry + "/v2/" + image + "/blobs/" + digest

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	// store file
	file, err := os.OpenFile(karakuripkgs.IMAGE_ROOT+"/"+image_id+"/config.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(byte_array)
}

func privGetLayer(repositry string, digest string, image string, image_id string, index int) {
	url := "http://" + repositry + "/v2/" + image + "/blobs/" + digest

	req, _ := http.NewRequest("GET", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byte_array, _ := io.ReadAll(resp.Body)

	// store file
	filepath := karakuripkgs.IMAGE_ROOT + "/" + image_id + "/"
	filename := "layer" + strconv.Itoa(index) + ".tar.gz"
	destination := filepath + filename
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(byte_array)
	// extract layer
	cmd := exec.Command("tar", "zxvf", destination, "-C", filepath+"rootfs")
	//if err := cmd.Run(); err != nil {
	//	panic(err)
	//}
	cmd.Run()
	// remove tar file
	if err := os.Remove(destination); err != nil {
		panic(err)
	}
}

func pullProcess(image string, tag string, os string, arch string) {
	// get token
	token := getToken(image)
	// get manifest
	manifets_list := getManifestList(image, tag, token)
	// retrieve manifest digest
	var target_manifest_digest string
	for _, entry := range manifets_list.ManifetList {
		if entry.Platform.Os == os && entry.Platform.Architecture == arch {
			target_manifest_digest = entry.Digest
		}
	}
	// get manifest blog
	blob := getBlob(target_manifest_digest, image, token)

	// config digest
	config_digest := blob.Config.Digest
	// layer
	layers := blob.Layers
	// image id
	image_id := config_digest[7:19]
	// rootfs
	rootfs := karakuripkgs.IMAGE_ROOT + "/" + image_id + "/rootfs"

	// create image directory
	createImageDirectory(image_id)

	// get config file
	getConfig(config_digest, image, token, image_id)

	// get layers
	for i, entry := range layers {
		digest := entry.Digest
		getLayer(digest, image, token, image_id, i)
	}

	// add image list
	AddImageList(image, tag, image_id, rootfs)
}

func pullPrivProcess(image string, tag string, repositry string) {
	// get manifest
	manifest := privGetManifest(repositry, image, tag)
	// config digest
	config_digest := manifest.Config.Digest
	// layer
	layers := manifest.Layers
	// image id
	image_id := config_digest[7:19]
	// rootfs
	rootfs := karakuripkgs.IMAGE_ROOT + "/" + image_id + "/rootfs"

	// create image directory
	createImageDirectory(image_id)

	// get config file
	privGetConfig(repositry, config_digest, image, image_id)

	// get layers
	for i, entry := range layers {
		digest := entry.Digest
		privGetLayer(repositry, digest, image, image_id, i)
	}

	// add image list
	AddImageList(image, tag, image_id, rootfs)
}

func PullImage(image_tag string, os_arch string, repositry string) ResponsePullImage {
	// parse image:tag
	image_tag_info := strings.Split(image_tag, ":")
	image := image_tag_info[0]
	tag := "latest"
	if len(image_tag_info) == 2 {
		tag = image_tag_info[1]
	}

	// parse os/archtecture
	os_arch_info := strings.Split(os_arch, ":")
	os := os_arch_info[0]
	arch := os_arch_info[1]

	// check if image already exists
	if isImageExists(image, tag) {
		return createResponsePullImage("success", true, image, tag, os, arch)
	} else {
		if repositry == "public" {
			pullProcess(image, tag, os, arch)
		} else {
			pullPrivProcess(image, tag, repositry)
		}
	}

	return createResponsePullImage("success", false, image, tag, os, arch)
}
