package ossrepo

import "fmt"

func (p *Package) getInstallCommand() string {
	return fmt.Sprintf("brew install %s", p.HomebrewPackageName)
}
