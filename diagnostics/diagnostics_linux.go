package diagnostics

func getPackageManager() packageManager {
	return PackageManagerAPT // || YUM || APK ... TODO
}

func getOSRelease() string {
	return "?TODO?"
}
