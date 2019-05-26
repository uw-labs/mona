package command

// Test attempts to run the test command for all modules where changes
// are detected.
func Test() error {
	return rangeChangedModules(changeTypeTest, true, testModule)
}
