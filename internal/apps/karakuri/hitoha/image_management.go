package hitoha

import (
	"encoding/json"
	"karakuripkgs"
	"os"
)

// docker registory api
type RespToken struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	IssuedAt  string `json:"issued_at"`
}

type PlatformInfo struct {
	Architecture string `json:"architecture"`
	Os           string `json:"os"`
}

type Manifest struct {
	Digest   string       `json:"digest"`
	Size     int          `json:"size"`
	Platform PlatformInfo `json:"platform"`
}

type ManifestList struct {
	ManifetList   []Manifest `json:"manifests"`
	MediaType     string     `json:"mediaType"`
	SchemeVersion int        `json:"schemaVersion"`
}

type ManifestConfig struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type ManifestLayer struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type ManifestBlob struct {
	SchemaVersion int             `json:"schemaVersion"`
	MediaType     string          `json:"mediaType"`
	Config        ManifestConfig  `json:"config"`
	Layers        []ManifestLayer `json:"layers"`
}

// private repo
type PrivManifest struct {
	SchemaVersion int             `json:"schemaVersion"`
	MediaType     string          `json:"mediaType"`
	Config        ManifestConfig  `json:"config"`
	Layers        []ManifestLayer `json:"layers"`
}

// image list
type ImageInfo struct {
	Image   string `json:"image"`
	Tag     string `json:"tag"`
	ImageId string `json:"image_id"`
	Rootfs  string `json:"rootfs"`
}

type ImageList struct {
	List []ImageInfo `json:"list"`
}

// blob
type BlobConfig struct {
	Cmd        []string `json:"Cmd"`
	Entrypoint []string `json:"Entrypoint"`
	Env        []string `json:"Env"`
}

type BlobFile struct {
	Config BlobConfig `json:"config"`
}

func createImageDirectory(image_id string) {
	if err := os.MkdirAll(karakuripkgs.IMAGE_ROOT+"/"+image_id+"/rootfs", os.ModePerm); err != nil {
		panic(err)
	}
}

func newImageList() {
	var image_list_data ImageList
	data, _ := json.MarshalIndent(image_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_IMAGE_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func AddImageList(image string, tag string, image_id string, rootfs string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_IMAGE_LIST)
	if err != nil {
		panic(err)
	}

	var image_list_data ImageList
	if err := json.Unmarshal(bytes, &image_list_data); err != nil {
		panic(err)
	}

	var image_info = ImageInfo{
		Image:   image,
		Tag:     tag,
		ImageId: image_id,
		Rootfs:  rootfs,
	}
	image_list_data.List = append(image_list_data.List, image_info)

	data, _ := json.MarshalIndent(image_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_IMAGE_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func isImageExists(image string, tag string) bool {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_IMAGE_LIST)
	if err != nil {
		newImageList()
		return false
	}

	var image_list_data ImageList
	if err := json.Unmarshal(bytes, &image_list_data); err != nil {
		panic(err)
	}

	for _, entry := range image_list_data.List {
		if entry.Image == image && entry.Tag == tag {
			return true
		}
	}
	return false
}

// delete
func deleteImageDirectory(image_id string) {
	image_path := karakuripkgs.IMAGE_ROOT + "/" + image_id
	if err := os.RemoveAll(image_path); err != nil {
		panic(err)
	}
}

func DeleteImage(image_id string) ResponseDeleteImage {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_IMAGE_LIST)
	if err != nil {
		return createResponseDeleteImage("error", image_id)
	}

	var image_list_data ImageList
	if err := json.Unmarshal(bytes, &image_list_data); err != nil {
		panic(err)
	}

	// delete image directory
	var new_image_list ImageList
	for _, entry := range image_list_data.List {
		if entry.ImageId == image_id {
			deleteImageDirectory(image_id)
		} else {
			new_image_list.List = append(new_image_list.List, entry)
		}
	}

	// update image list
	new_data, _ := json.MarshalIndent(new_image_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_IMAGE_LIST, new_data, os.ModePerm); err != nil {
		panic(err)
	}

	return createResponseDeleteImage("success", image_id)
}

func ShowImageList() ResponseImageList {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_IMAGE_LIST)
	if err != nil {
		panic(err)
	}

	var image_list_data ImageList
	if err := json.Unmarshal(bytes, &image_list_data); err != nil {
		panic(err)
	}

	return createResponseImageList("success", image_list_data)
}

func getImageRootfs(image string, tag string) string {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_IMAGE_LIST)
	if err != nil {
		panic(err)
	}

	var image_list_data ImageList
	if err := json.Unmarshal(bytes, &image_list_data); err != nil {
		panic(err)
	}

	for _, entry := range image_list_data.List {
		if entry.Image == image && entry.Tag == tag {
			return entry.Rootfs
		}
	}

	return ""
}

func GetImageId(image string, tag string) string {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_IMAGE_LIST)
	if err != nil {
		panic(err)
	}

	var image_list_data ImageList
	if err := json.Unmarshal(bytes, &image_list_data); err != nil {
		panic(err)
	}

	for _, entry := range image_list_data.List {
		if entry.Image == image && entry.Tag == tag {
			return entry.ImageId
		}
	}

	return ""
}

// blob
func readBlobFile(image_id string) BlobFile {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.IMAGE_ROOT + "/" + image_id + "/config.json")
	if err != nil {
		panic(err)
	}

	var blob_file BlobFile
	if err := json.Unmarshal(bytes, &blob_file); err != nil {
		panic(err)
	}

	return blob_file
}
