package karakuri_mod

import (
	"errors"
	"karakuripkgs"
	"net"
	"os"
	"os/exec"
	"strings"
)

func createIngressModuleKarakurifile() {
	// create Karakurifile
	karakurifile := `FROM nginx:alpine

COPY /etc/karakuri/modules/ingress/script /script
COPY /etc/karakuri/modules/ingress/certconf /certconf
COPY /etc/karakuri/modules/ingress/conf /conf

RUN apk update && apk add openssl

RUN mkdir rootca interca server
RUN openssl genrsa -out rootca/rootca.key
RUN openssl req -new -key rootca/rootca.key -out rootca/rootca.csr -subj "/C=JP/ST=Tokyo/CN=Karakuri Ingress RootCA"
RUN openssl x509 -req -in rootca/rootca.csr -signkey rootca/rootca.key -days 365 -sha256 -extfile /certconf/ca_v3.ext -out rootca/rootca.crt

RUN openssl genrsa -out interca/interca.key
RUN openssl req -new -key interca/interca.key -out interca/interca.csr -subj "/C=JP/ST=Tokyo/CN=Karakuri Ingress InterCA"
RUN openssl x509 -req -in interca/interca.csr -CA rootca/rootca.crt -CAkey rootca/rootca.key -CAcreateserial -days 365 -sha256 -out interca/interca.crt -extfile /certconf/ca_v3.ext

RUN openssl genrsa -out server/server.key
RUN openssl req -new -key server/server.key -out server/server.csr -subj "/C=JP/ST=Tokyo/CN=Karakuri Ingress Host"
RUN openssl x509 -req -in server/server.csr -CA interca/interca.crt -CAkey interca/interca.key -CAcreateserial -days 365 -sha256 -out server/server.crt -extfile /certconf/server_v3.ext

RUN cat server/server.crt interca/interca.crt rootca/rootca.crt > server/chain.crt

RUN cp /conf/server.conf /etc/nginx/conf.d/

CMD ["sh", "/script/entrypoint.sh"]
`

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/Karakurifile"); stat != nil {
		fd, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/Karakurifile")
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

func isPrivateIP(ip net.IP) bool {
	var prvMasks []*net.IPNet

	for _, cidr := range []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	} {
		_, mask, _ := net.ParseCIDR(cidr)
		prvMasks = append(prvMasks, mask)
	}

	for _, mask := range prvMasks {
		if mask.Contains(ip) {
			return true
		}
	}
	return false
}

func getDeviceIpAddress() (net.IP, error) {
	ift, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifi := range ift {
		addrs, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

func createIngressModuleConfig() {
	// conf
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/conf"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS+"/conf", os.ModePerm); err != nil {
			panic(err)
		}
	}

	server_conf := `server{
    server_name    ingress.karakuri.container;

    listen 443 ssl;
    ssl_certificate /server/chain.crt;
    ssl_certificate_key /server/server.key;

    client_max_body_size 5G;
    keepalive_timeout  130;
    send_timeout 130;
    client_body_timeout 130;
    client_header_timeout 130;
    proxy_send_timeout 130;
    proxy_read_timeout 130;
}
`
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/conf/server.conf"); stat != nil {
		fd_1, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/conf/server.conf")
		if err != nil {
			panic(err)
		}
		defer fd_1.Close()
		bytes_1 := []byte(server_conf)
		if _, err := fd_1.Write(bytes_1); err != nil {
			panic(err)
		}
	}

	// cert
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/certconf"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS+"/certconf", os.ModePerm); err != nil {
			panic(err)
		}
	}

	device_address, err := getDeviceIpAddress()
	if err != nil {
		panic(err)
	}
	ca_v3 := `basicConstraints=critical, CA:true
authorityKeyIdentifier=keyid:always,issuer
subjectKeyIdentifier=hash
extendedKeyUsage=serverAuth, clientAuth
keyUsage=keyCertSign,cRLSign
`
	server_v3 := `authorityKeyIdentifier=keyid,issuer
basicConstraints=critical, CA:FALSE
keyUsage=digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth, clientAuth
subjectAltName=@alt_names

[alt_names]
DNS.1 = karakuri.container
DNS.2 = *.karakuri.container
IP.1 = ` + device_address.String() + "\n"

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/certconf/ca_v3.ext"); stat != nil {
		fd_2, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/certconf/ca_v3.ext")
		if err != nil {
			panic(err)
		}
		defer fd_2.Close()
		bytes_2 := []byte(ca_v3)
		if _, err := fd_2.Write(bytes_2); err != nil {
			panic(err)
		}
	}
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/certconf/server_v3.ext"); stat != nil {
		fd_3, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/certconf/server_v3.ext")
		if err != nil {
			panic(err)
		}
		defer fd_3.Close()
		bytes_3 := []byte(server_v3)
		if _, err := fd_3.Write(bytes_3); err != nil {
			panic(err)
		}
	}

	// script
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/script"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS+"/script", os.ModePerm); err != nil {
			panic(err)
		}
	}
	entrypoint := `#!/bin/sh

/docker-entrypoint.sh nginx -g "daemon off;"
`
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS + "/script/entrypoint.sh"); stat != nil {
		fd_4, err := os.Create(karakuripkgs.KARAKURI_MOD_INGRESS + "/script/entrypoint.sh")
		if err != nil {
			panic(err)
		}
		defer fd_4.Close()
		bytes_4 := []byte(entrypoint)
		if _, err := fd_4.Write(bytes_4); err != nil {
			panic(err)
		}
	}
}

func setupIngressModule() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_INGRESS); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_INGRESS, os.ModePerm); err != nil {
			panic(err)
		}
	}
	// Karakurifile
	createIngressModuleKarakurifile()
	// ingress config
	createIngressModuleConfig()
}

func buildIngressImage(mod_info ModInfo) {
	build_args := []string{"build", "--name", mod_info.ImageName, "--buildpath", mod_info.Path}
	build := exec.Command("karakuri", build_args...)
	if err := build.Run(); err != nil {
		panic(err)
	}
}

func createIngressContainer(mod_info ModInfo) {
	create_args := []string{
		"create",
		"--name", mod_info.Name,
		"--image", mod_info.ImageName,
		"--port", "4443:443:tcp",
		"--restart", "on-boot",
		"--ns", "system-mod",
	}
	create := exec.Command("karakuri", create_args...)
	if err := create.Run(); err != nil {
		panic(err)
	}
}

func startIngressContainer(mod_info ModInfo) {
	start_args := []string{"start", "--name", mod_info.Name}
	start := exec.Command("karakuri", start_args...)
	if err := start.Run(); err != nil {
		panic(err)
	}
}

func enableIngressModule(mod_info ModInfo) {
	image_info := strings.Split(mod_info.ImageName, ":")
	if !isImageExists(image_info[0], image_info[1]) {
		buildIngressImage(mod_info)
	}
	createIngressContainer(mod_info)
	startIngressContainer(mod_info)
}

func disableIngressModule() {
	// remove container
	removeModuleContainer("ingress")
}

func EditIngressEntry() {
	edit_args := []string{"exec", "--it", "--name", "ingress", "--cmd", "vi,/etc/nginx/conf.d/server.conf"}
	edit := exec.Command("karakuri", edit_args...)

	edit.Stdin = os.Stdin
	edit.Stdout = os.Stdout
	edit.Stderr = os.Stderr

	if err := edit.Start(); err != nil {
		panic(err)
	}
	edit.Wait()
	// restart
	restart_args := []string{"restart", "--name", "ingress"}
	restart := exec.Command("karakuri", restart_args...)
	if err := restart.Run(); err != nil {
		panic(err)
	}
}
