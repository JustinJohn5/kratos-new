package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CmdRun run project command.
var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run project",
	Long:  "Run project. Example: kratos run",
	RunE:  Run,
}
var targetDir string

func init() {
	CmdRun.Flags().StringVarP(&targetDir, "work", "w", "", "target working directory")
}

// Run run project.
func Run(cmd *cobra.Command, args []string) error {
	var dir string
	cmdArgs, programArgs := splitArgs(cmd, args)
	base, err := os.Getwd()
	if err != nil {
		return err
	}
	if len(cmdArgs) > 0 {
		dir = cmdArgs[0]
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(base, dir)
		}
		if info, statErr := os.Stat(dir); statErr != nil || !info.IsDir() {
			return fmt.Errorf("%s is not a valid directory", cmdArgs[0])
		}
	}
	if dir == "" {
		// find the directory containing the cmd/*
		cmdPath, err := findCMD(base)
		if err != nil {
			return err
		}
		switch len(cmdPath) {
		case 0:
			return fmt.Errorf("the cmd directory cannot be found in the current directory")
		case 1:
			for _, v := range cmdPath {
				dir = v
			}
		default:
			var cmdPaths []string
			for k := range cmdPath {
				cmdPaths = append(cmdPaths, k)
			}
			prompt := &survey.Select{
				Message:  "Which directory do you want to run?",
				Options:  cmdPaths,
				PageSize: 10,
			}
			e := survey.AskOne(prompt, &dir)
			if e != nil || dir == "" {
				return e
			}
			dir = cmdPath[dir]
		}
	}
	fd := exec.Command("go", append([]string{"run", "."}, programArgs...)...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = dir
	changeWorkingDirectory(fd, targetDir)
	if err := fd.Run(); err != nil {
		return fmt.Errorf("run project: %w", err)
	}
	return nil
}

func splitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	return args, []string{}
}

func findCMD(base string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	wd = filepath.Clean(wd)
	if !strings.HasSuffix(wd, string(filepath.Separator)) {
		wd += string(filepath.Separator)
	}
	var root bool
	next := func(dir string) (map[string]string, error) {
		cmdPath := make(map[string]string)
		err := filepath.Walk(dir, func(walkPath string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			// multi level directory is not allowed under the cmdPath directory, so it is judged that the path ends with cmdPath.
			if info.IsDir() && filepath.Base(walkPath) == "cmd" {
				paths, err := os.ReadDir(walkPath)
				if err != nil {
					return err
				}
				for _, fileInfo := range paths {
					if fileInfo.IsDir() {
						abs := filepath.Join(walkPath, fileInfo.Name())
						cmdPath[strings.TrimPrefix(abs, wd)] = abs
					}
				}
				return nil
			}
			if info.Name() == "go.mod" {
				root = true
			}
			return nil
		})
		return cmdPath, err
	}
	for i := 0; i < 5; i++ {
		tmp := base
		cmd, err := next(tmp)
		if err != nil {
			return nil, err
		}
		if len(cmd) > 0 {
			return cmd, nil
		}
		if root {
			break
		}
		base = filepath.Join(base, "..")
	}
	return map[string]string{"": base}, nil
}

func changeWorkingDirectory(cmd *exec.Cmd, targetDir string) {
	targetDir = strings.TrimSpace(targetDir)
	if targetDir != "" {
		cmd.Dir = targetDir
	}
}
