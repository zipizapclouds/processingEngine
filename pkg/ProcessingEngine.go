package processingEngine

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// ProcessingEngine is a struct that can run a binary with arguments and environment variables
// and capture the stdout, stderr, and exit code
// The environment variables are read from an envFilePath, which should contain lines of the form KEY=VALUE (i.e. matching regexp ^[^#]*=.*)
// that will be appended to the command's environment
type ProcessingEngine struct {
	// binPath should be path to executable file
	binPath string

	// "" or a readable file
	envFilePath string
	args        []string

	stdout   string
	stderr   string
	exitCode int
}

func NewProcessingEngine(binPath string, envFilePath string, args []string) *ProcessingEngine {
	return &ProcessingEngine{
		binPath:     binPath,
		envFilePath: envFilePath,
		args:        args,
	}
}

// Run the binpath with args and envFilePath, and set stdout, stderr, and exitCode
//
// Verifications
// - Verify that pe.binPath is a file with executable permissions
// - Verify that pe.envFilePath is a readable file
//
// # Execute pe.binPath with pe.args and pe.envFilePath, and set pe.stdout, pe.stderr, and pe.exitCode
//
// The envFilePath lines of the form KEY=VALUE (i.e. matching regexp ^[^#]*=.*) will be appended to the command's environment
func (pe *ProcessingEngine) Run() (exitCode int, er error) {
	// Verifications
	// - Verify that pe.binPath is a file with executable permissions
	// - Verify that pe.envFilePath is a readable file
	{
		// Verify that pe.binPath is a file with executable permissions
		if stat, err := os.Stat(pe.binPath); err != nil {
			return 0, fmt.Errorf("error verifying binary at %s: %s", pe.binPath, err)
		} else if !stat.Mode().IsRegular() {
			return 0, fmt.Errorf("binary at %s is not a regular file", pe.binPath)
		} else if stat.Mode()&0111 == 0 {
			return 0, fmt.Errorf("binary at %s is not executable", pe.binPath)
		}

		// Verify that pe.envFilePath is "" or a readable file
		if pe.envFilePath != "" {
			if stat, err := os.Stat(pe.envFilePath); err != nil {
				return 0, fmt.Errorf("error verifying environment file at %s: %s", pe.envFilePath, err)
			} else if !stat.Mode().IsRegular() {
				return 0, fmt.Errorf("environment file at %s is not a regular file", pe.envFilePath)
			} else if stat.Mode()&0444 == 0 {
				return 0, fmt.Errorf("environment file at %s is not readable", pe.envFilePath)
			}
		}
	}

	// Execute pe.binPath with pe.args and pe.envFilePath, and set pe.stdout, pe.stderr, and pe.exitCode.
	// The envFilePath lines of the form KEY=VALUE (i.e. matching regexp ^[^#]*=.*) will be appended to the command's environment
	// - set command.Env, containing os.Environ() appended with pe.envFilePath lines that match regexp ^[^#]*=.*
	{
		command := exec.Command(pe.binPath, pe.args...)
		// The envFilePath lines of the form KEY=VALUE (i.e. matching regexp ^[^#]*=.*) will be appended to the command's environment
		{
			command.Env = os.Environ()

			if pe.envFilePath != "" {
				dat, err := os.ReadFile(pe.envFilePath)
				if err != nil {
					return 0, fmt.Errorf("error opening environment file at %s: %s", pe.envFilePath, err)
				}

				for _, line := range strings.Split(string(dat), "\n") {
					// if line of envFilePath matches regexp ^[^#]*=.*
					if !regexp.MustCompile(`^[^#]*=.*`).MatchString(line) {
						continue
					}
					// - append matching line to command.Env
					command.Env = append(command.Env, line)
				}
			}
		}
		var stdoutb, stderrb bytes.Buffer
		command.Stdout = &stdoutb
		command.Stderr = &stderrb
		err := command.Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				pe.exitCode = exitError.ExitCode()
			} else {
				return 0, fmt.Errorf("error when running the command: %s", err)
			}
		}
		pe.stdout = stdoutb.String()
		pe.stderr = stderrb.String()
	}

	exitCode = pe.exitCode
	er = nil
	return
}

func (pe *ProcessingEngine) GetStdout() string {
	return pe.stdout
}
func (pe *ProcessingEngine) GetStderr() string {
	return pe.stderr
}
func (pe *ProcessingEngine) GetExitCode() int {
	return pe.exitCode
}
func (pe *ProcessingEngine) GetBinPath() string {
	return pe.binPath
}
func (pe *ProcessingEngine) GetEnvFilePath() string {
	return pe.envFilePath
}
func (pe *ProcessingEngine) GetArgs() []string {
	return pe.args
}
