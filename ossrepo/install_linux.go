package ossrepo

import "fmt"

func (p *Package) getInstallCommand() string {
	return fmt.Sprintf("apt install %s", p.DebianPackageName)
}
