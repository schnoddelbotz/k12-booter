package ossrepo

type Platform int
type Category int
type Subject int

const (
	Platform_MacOS Platform = iota
	Platform_Linux
	Platform_Windows
	Category_EDU Category = iota
	Category_Media
	Category_IT
	Subject_Maths Subject = iota
	Subject_Physics
	Subject_Chemistry
	Subject_History
	Subject_Ethics
	Subject_Languages
	Subject_Programming
	Subject_Networking
	Subject_Hardware

	AllPlatforms = Platform_MacOS & Platform_Linux & Platform_Windows
)

type Package struct {
	Name                  string
	Platforms             Platform
	Category              Category
	DescriptionShort      string
	DescriptionLong       string
	Homepage              string
	HomebrewPackageName   string
	DebianPackageName     string
	ChocolateyPackageName string
}

var DB []Package = []Package{
	{
		Name:                  "Blender",
		Platforms:             AllPlatforms,
		Category:              Category_Media,
		DescriptionShort:      "3D-modelling & animation",
		Homepage:              "https://blender.org/",
		HomebrewPackageName:   "blender",
		ChocolateyPackageName: "blender",
		DebianPackageName:     "blender",
	},
}
