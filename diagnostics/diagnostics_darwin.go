package diagnostics

import "os/exec"

func getPackageManager() packageManager {
	return PackageManagerHomebrew
}

func getOSRelease() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, _ := cmd.CombinedOutput()
	return string(output)
}
