package spa

// RunnerType runner type.
type RunnerType string

const (
	// RunnerTypeNpm runner type NPM.
	RunnerTypeNpm RunnerType = "npm"

	// RunnerTypeNpx runner type NPX.
	RunnerTypeNpx RunnerType = "npx"

	// RunnerTypeYarn runner type Yarn.
	RunnerTypeYarn RunnerType = "yarn"

	// RunnerTypeCustom custom runner type.
	RunnerTypeCustom RunnerType = "custom"
)

func prepareRunner(runnerType RunnerType, scriptName string, args ...string) (string, []string) {
	if runnerType == RunnerTypeCustom {
		return scriptName, args
	}

	var path string
	switch runnerType {
	case RunnerTypeNpm:
		path = "npm"
	case RunnerTypeNpx:
		path = "npx"
	case RunnerTypeYarn:
		path = "yarn"
	}

	a := make([]string, 0, len(args)+1)

	if runnerType == RunnerTypeNpm {
		a = append(a, "run")
	}

	a = append(a, scriptName)

	if runnerType == RunnerTypeNpm {
		a = append(a, "--")
	}

	a = append(a, args...)

	return path, a
}
