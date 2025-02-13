package hitoha

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"karakuripkgs"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func compressLayer(id string) {
	// create .tar
	tar_cmd := exec.Command("tar", "-cvf", karakuripkgs.IMAGE_ROOT+"/"+id+"/layer.tar", "-C", karakuripkgs.IMAGE_ROOT+"/"+id+"/rootfs", ".")
	if err := tar_cmd.Run(); err != nil {
		panic(err)
	}
	// create .tar.gz
	gzip_cmd := exec.Command("gzip", karakuripkgs.IMAGE_ROOT+"/"+id+"/layer.tar")
	if err := gzip_cmd.Run(); err != nil {
		panic(err)
	}
}

func getUploadUrl(registry, image string) string {
	url := "http://" + registry + "/v2/" + image + "/blobs/uploads/"

	req, _ := http.NewRequest("POST", url, nil)

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	upload_url := resp.Header.Get("Location")
	return upload_url
}

func calcSha256(file, id string) string {
	file_path := karakuripkgs.IMAGE_ROOT + "/" + id + "/" + file
	file_data, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer file_data.Close()

	// calculate sha256 hash
	hash := sha256.New()
	if _, err := io.Copy(hash, file_data); err != nil {
		panic(err)
	}
	hash_bytes := hash.Sum(nil)
	hash_string := hex.EncodeToString(hash_bytes)

	return hash_string
}

func calcFileSize(file, id string) int {
	file_path := karakuripkgs.IMAGE_ROOT + "/" + id + "/" + file
	file_info, err := os.Stat(file_path)
	if err != nil {
		panic(err)
	}
	file_size := file_info.Size()
	return int(file_size)
}

func uploadImage(file, digest, registry, image, id string) bool {
	// open layer.tar.gz
	layer_path := karakuripkgs.IMAGE_ROOT + "/" + id + "/" + file
	data, err := os.Open(layer_path)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	// get upload url
	upload_url := getUploadUrl(registry, image)
	url := upload_url + "&digest=sha256:" + digest

	req, _ := http.NewRequest("PUT", url, data)
	req.Header.Set("Content-Type", "application/octet-stream")

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}

func createManifest(id, config_digest, config_size, layer_digest, layer_size string) {
	manifest := `{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "size": ` + config_size + `,
    "digest": "sha256:` + config_digest + `"
  },
  "layers": [
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": ` + layer_size + `,
      "digest": "sha256:` + layer_digest + `"
    }
  ]
}
`
	fd, err := os.Create(karakuripkgs.IMAGE_ROOT + "/" + id + "/manifest.json")
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	bytes := []byte(manifest)
	if _, err := fd.Write(bytes); err != nil {
		panic(err)
	}
}

func uploadManifest(registry, image, tag, id string) bool {
	// open layer.tar.gz
	manifest_path := karakuripkgs.IMAGE_ROOT + "/" + id + "/manifest.json"
	data, err := os.Open(manifest_path)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	url := "http://" + registry + "/v2/" + image + "/manifests/" + tag

	req, _ := http.NewRequest("PUT", url, data)
	req.Header.Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")

	http_client := new(http.Client)
	resp, err := http_client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}

func removePushFiles(id string) {
	image_path := karakuripkgs.IMAGE_ROOT + "/" + id + "/"
	// layer.tar.gz
	if err := os.Remove(image_path + "layer.tar.gz"); err != nil {
		panic(err)
	}
	// manifest.json
	if err := os.Remove(image_path + "manifest.json"); err != nil {
		panic(err)
	}
}

func PushImage(registry, image_tag string) ResponsePushImage {
	// parse image:tag
	image_tag_info := strings.Split(image_tag, ":")
	image := image_tag_info[0]
	tag := "latest"
	if len(image_tag_info) == 2 {
		tag = image_tag_info[1]
	}

	// check if image exists
	if !isImageExists(image, tag) {
		return createResponsePushImage("error", image_tag, "no such image, "+image_tag)
	}

	// get image id
	image_id := GetImageId(image, tag)
	// create layer data
	compressLayer(image_id)

	// upload layer
	// calculate file sha256 sum and size
	layer_digest := calcSha256("layer.tar.gz", image_id)
	layer_size := calcFileSize("layer.tar.gz", image_id)
	if res := uploadImage("layer.tar.gz", layer_digest, registry, image, image_id); !res {
		return createResponsePushImage("error", image_tag, "layer: "+layer_digest+" upload failed")
	}
	// upload config
	config_digest := calcSha256("config.json", image_id)
	config_size := calcFileSize("config.json", image_id)
	if res := uploadImage("config.json", config_digest, registry, image, image_id); !res {
		return createResponsePushImage("error", image_tag, "config: "+config_digest+" upload failed")
	}
	// upload manifest
	createManifest(image_id, config_digest, strconv.Itoa(config_size), layer_digest, strconv.Itoa(layer_size))
	if res := uploadManifest(registry, image, tag, image_id); !res {
		return createResponsePushImage("error", image_tag, "manifest upload failed")
	}

	// remove files
	removePushFiles(image_id)

	return createResponsePushImage("success", image_tag, "push completed")
}
