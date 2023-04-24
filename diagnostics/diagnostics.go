package diagnostics

import (
	"fmt"
	"log"
	"net"
	"os/user"
	"runtime"
)

type SysInfoData struct {
	OS             string
	OSRelease      string
	LinuxDistro    string
	Architecture   string
	Username       string
	PackageManager packageManager
	LANAddresses   []string
}

type packageManager string

const (
	PackageManagerHomebrew   packageManager = "homebrew"
	PackageManagerAPT        packageManager = "apt"
	PackageManagerYUM        packageManager = "yum"
	PackageManagerAPK        packageManager = "apk"
	PackageManagerChocolatey packageManager = "chocolatey"
)

func GetSysInfoData() *SysInfoData {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &SysInfoData{
		Username:       user.Username,
		OS:             runtime.GOOS,
		OSRelease:      getOSRelease(),
		Architecture:   runtime.GOARCH,
		PackageManager: getPackageManager(),
		LANAddresses:   getLocalAddresses(),
	}
}

func SysInfo() string {
	data := GetSysInfoData()
	return fmt.Sprintf("%+v\n", data)
}

func getLocalAddresses() []string {
	// https://stackoverflow.com/questions/23529663/how-to-get-all-addresses-and-masks-from-local-interfaces-in-go
	var res []string
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if ipv4Addr := a.(*net.IPNet).IP.To4(); ipv4Addr != nil {
					res = append(res, fmt.Sprintf("%v(%s)", i.Name, v))
				}
			}
		}
	}
	return res
}
