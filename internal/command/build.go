package command

// Build will execute the build commands for all modules where changes
// are detected.
func Build() error {
	return rangeChangedModules(changeTypeBuild, true, buildModule)
}
