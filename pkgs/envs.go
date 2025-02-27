package karakuripkgs

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

// mod
const (
	KARAKURI_MOD_ROOT = KARAKURI_ROOT + "/modules"
	KARAKURI_MOD_LIST = KARAKURI_MOD_ROOT + "/module_list.json"

	// modules
	KARAKURI_MOD_DNS              = KARAKURI_MOD_ROOT + "/dns"
	KARAKURI_MOD_REGISTRY         = KARAKURI_MOD_ROOT + "/registry"
	KARAKURI_MOD_REGISTRY_BROWSER = KARAKURI_MOD_ROOT + "/registry-browser"
	KARAKURI_MOD_INGRESS          = KARAKURI_MOD_ROOT + "/ingress"
)

// registry controller
const (
	KARAKURI_REGCTL_ROOT    = KARAKURI_ROOT + "/regctl"
	KARAKURI_REGCTL_REGINFO = KARAKURI_REGCTL_ROOT + "/registry.json"
)

// cluster controller
const (
	KARAKURI_NODECTL_ROOT            = KARAKURI_ROOT + "/node"
	KARAKURI_NODECTL_NODEINFO        = KARAKURI_NODECTL_ROOT + "/node.json"
	KARAKURI_NODECTL_AUTHCODE        = KARAKURI_NODECTL_ROOT + "/authcode.json"
	KARAKURI_NODECTL_REMOTE_AUTHCODE = KARAKURI_NODECTL_ROOT + "/remote_authcode.json"
)

// runtime
const (
	RUNTIME = "/bin/futaba"
	SERVER  = "localhost:9806"
)
