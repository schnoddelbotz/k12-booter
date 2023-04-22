package ossrepo

import "fmt"

func (p *Package) getInstallCommand() string {
	return fmt.Sprintf("chocolatey install %s", p.ChocolateyPackageName)
}
