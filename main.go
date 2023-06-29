package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

const programName = "crossjoin"
const programVersion = "1.0.0"
const programDescription = `Generate the cross join (or Cartesian product) of lines from the files specified.

For instance, it can combine http:// or https:// from file1, various domains
from file2, and assorted endpoints from file3, effectively creating a
comprehensive list for tasks such as fuzzing or penetration testing.`

type arguments struct {
	help    bool
	version bool

	files []string
}

func parseArguments() arguments {
	args := arguments{}

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s file1 file2 file3 [...fileN]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		pflag.PrintDefaults()
	}

	flags := pflag.NewFlagSet("", pflag.ExitOnError)
	flags.BoolVarP(&args.help, "help", "h", false, "Show help")
	flags.BoolVarP(&args.version, "version", "v", false, "Show version")

	pflag.CommandLine.AddFlagSet(flags)
	pflag.Parse()

	if args.help {
		usage()
		os.Exit(0)
	}

	if args.version {
		fmt.Println(programName, programVersion)
		os.Exit(0)
	}

	args.files = pflag.Args()

	return args
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s %s\n", programName, programVersion)
	fmt.Fprintf(os.Stderr, "%s\n\n", programDescription)
	pflag.Usage()
}

func main() {
	args := parseArguments()
	if err := process(args.files); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func process(filenames []string) error {
	fileCount := len(filenames)
	if fileCount == 0 {
		return fmt.Errorf("no input file specified")
	}

	files := make([]*os.File, fileCount)
	scanners := make([]*bufio.Scanner, fileCount)
	lines := make([][]byte, fileCount)

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Initialize scanners and get first line from each file
	for i, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		files[i] = file
		scanners[i] = scanner
		lines[i] = scanner.Bytes()
	}

	for {
		// Print current line combination
		for _, line := range lines {
			writer.Write(line)
		}
		writer.WriteByte('\n')

		// Update combination starting from the last file
		for i := fileCount - 1; i >= 0; i-- {
			if scanners[i].Scan() {
				lines[i] = scanners[i].Bytes()
				break
			} else {
				// Check if we've cycled through all combinations
				if i == 0 {
					return nil
				}

				// Reached the end of this file, rewind and continue with the next one
				if _, err := files[i].Seek(0, 0); err != nil {
					return err
				}

				scanners[i] = bufio.NewScanner(files[i])
				scanners[i].Scan()
				lines[i] = scanners[i].Bytes()
			}
		}
	}
}
