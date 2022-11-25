package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// version is the published version of the utility
var version = 0

const (
	// VerboseFlag is the Verbose Flag
	VerboseFlag string = "verbose"
)

// initFlags is where command line flags are instantiated
func initFlags(flag *pflag.FlagSet) {

	// Verbose
	flag.BoolP(VerboseFlag, "v", false, "log messages at the debug level.")

	flag.SortFlags = false
}

// checkConfig is how the input to command line flags are checked
func checkConfig(v *viper.Viper) error {

	return nil
}

func main() {
	root := cobra.Command{
		Use:   "linefix [flags]",
		Short: "LineFix",
		Long:  "LineFix",
	}

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\nlinefix completion > /usr/local/etc/bash_completion.d/linefix",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	scanCommand := &cobra.Command{
		Use:                   "scan [flags]",
		DisableFlagsInUseLine: true,
		Short:                 "Fixes Newlines in given directory",
		Long:                  "Fixes Newlines in given directory and all subdirectories",
		RunE:                  scanFunction,
	}
	initFlags(scanCommand.Flags())
	root.AddCommand(scanCommand)

	fixCommand := &cobra.Command{
		Use:                   "fix [flags]",
		DisableFlagsInUseLine: true,
		Short:                 "Fixes Newlines in given directory",
		Long:                  "Fixes Newlines in given directory and all subdirectories",
		RunE:                  fixFunction,
	}
	initFlags(fixCommand.Flags())
	root.AddCommand(fixCommand)

	versionCommand := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 "Print the version",
		Long:                  "Print the version",
		RunE:                  versionFunction,
	}
	root.AddCommand(versionCommand)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func versionFunction(cmd *cobra.Command, args []string) error {
	if version == 0 {
		fmt.Println("Version: development")
		return nil
	}

	fmt.Println("Version: ", version)
	return nil
}

func scanFunction(cmd *cobra.Command, args []string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	err := cmd.ParseFlags(args)
	if err != nil {
		return err
	}

	flag := cmd.Flags()

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		return bindErr
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	verbose := v.GetBool(VerboseFlag)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(io.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Check the config and exit with usage details if there is a problem
	checkConfigErr := checkConfig(v)
	if checkConfigErr != nil {
		return checkConfigErr
	}

	if verbose {

		fmt.Println("Verbose: ", verbose)

	} else {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		// get all files in the directory and subdirectories
		files, err := getFiles(dir)
		if err != nil {
			return err
		}

		bar := progressbar.Default(
			int64(len(files)),
			"Scanning",
		)

		affectedFiles := []string{}

		// loop through the files and check for newlines
		for _, file := range files {
			bar.Add(1)
			// check for newlines
			if !checkForNewlines(file) {
				affectedFiles = append(affectedFiles, file)
			}
			time.Sleep(40 * time.Millisecond)
		}

		fmt.Println("")
		// print in color orange
		fmt.Printf("\033[33m%d\033[0m files affected by newline issues in \033[33m%s\033[0m and subdirectories: ", len(affectedFiles), dir)
		fmt.Println("")

		for _, file := range affectedFiles {
			// print in color orange
			fmt.Printf("\033[33m%s\033[0m ", file)
			fmt.Println("")
		}

	}

	return nil
}

func fixFunction(cmd *cobra.Command, args []string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	err := cmd.ParseFlags(args)
	if err != nil {
		return err
	}

	flag := cmd.Flags()

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		return bindErr
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	verbose := v.GetBool(VerboseFlag)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(io.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Check the config and exit with usage details if there is a problem
	checkConfigErr := checkConfig(v)
	if checkConfigErr != nil {
		return checkConfigErr
	}

	if verbose {

		fmt.Println("Verbose: ", verbose)

	} else {
		// scan recursively through the directory and search for files with no newlines
		fmt.Println("Scanning for files with no newlines")

		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		// get all files in the directory and subdirectories
		files, err := getFiles(dir)
		if err != nil {
			return err
		}

		bar := progressbar.Default(
			int64(len(files)),
			"Scanning",
		)

		affectedFiles := []string{}

		// loop through the files and check for newlines
		for _, file := range files {
			bar.Add(1)
			// check for newlines
			if !checkForNewlines(file) {
				affectedFiles = append(affectedFiles, file)
				addNewline(file)
			}
			time.Sleep(40 * time.Millisecond)
		}

		fmt.Println("")
		// print in color green
		fmt.Printf("\033[32m%d\033[0m files fixed newline issues in \033[32m%s\033[0m and subdirectories: ", len(affectedFiles), dir)
		fmt.Println("")

		for _, file := range affectedFiles {
			// print in color green
			fmt.Printf("\033[32m%s\033[0m ", file)
			fmt.Println("")
		}

	}

	return nil
}

func getFiles(dir string) ([]string, error) {
	var files []string

	// exclude the .git directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			// exclude the .git directory
			if info.Name() == ".git" {
				return nil
			}

			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return []string{}, nil
	}

	return files, nil
}

func checkForNewlines(file string) bool {
	// check file if it has newline at the end
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	// get the file size
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	// read the last byte
	bs := make([]byte, 1)
	_, err = f.ReadAt(bs, stat.Size()-1)
	if err != nil {
		return false
	}

	// check if the last byte is a newline
	if bs[0] == '\r' || bs[0] == '\n' {
		return true
	}

	return false
}

func addNewline(file string) {
	// open the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// write a newline to the file
	if _, err = f.WriteString("\r"); err != nil {
		return
	}
}
