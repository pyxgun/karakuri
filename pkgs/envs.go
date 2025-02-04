package karakuripkgs

// ====================================
//
//	Edit to match your environment
//
// ====================================
const (
	HOST_NIC = "eth0"
)

// ====================================
//
//	Don't edit from here on
//
// ====================================
// directories
const (
	// root directory
	KARAKURI_ROOT = "/etc/karakuri"
	HITOHA_ROOT   = KARAKURI_ROOT + "/hitoha"
	IMAGE_ROOT    = HITOHA_ROOT + "/images"
	FUTABA_ROOT   = KARAKURI_ROOT + "/futaba"
)

// management files
const (
	// hitoha
	HITOHA_CONTAINER_LIST = HITOHA_ROOT + "/container_list.json"
	HITOHA_IMAGE_LIST     = IMAGE_ROOT + "/image_list.json"
	HITOHA_NETWORK_LIST   = HITOHA_ROOT + "/network_list.json"
	HITOHA_NAMESPACE_LIST = HITOHA_ROOT + "/namespace_list.json"

	// futaba
	FUTABA_CONTAINER_LIST = FUTABA_ROOT + "/container_list.json"
)

// runtime
const (
	RUNTIME = "/bin/futaba"
	SERVER  = "http://localhost:9806"
)
