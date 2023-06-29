package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

const programName = "crossjoin"
const programVersion = "1.1.0"
const programDescription = `Generate a cross join, also known as a Cartesian product, from the lines of the
specified files. If standard input (stdin) is provided, the program will use it
as the first input.`

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
	readerCount := len(filenames)

	if hasStdin() {
		readerCount++
	}

	if readerCount == 0 {
		return fmt.Errorf("no input specified")
	}

	readers := make([]io.ReadSeeker, readerCount)
	scanners := make([]*bufio.Scanner, readerCount)
	lines := make([][]byte, readerCount)

	// Initialize stdin
	inputIndex := 0
	if hasStdin() {
		readers[inputIndex] = os.Stdin
		inputIndex++
	}

	// Initialize files
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		readers[inputIndex] = file
		inputIndex++
	}

	// Initialize scanners and get first line from each
	for i := 0; i < readerCount; i++ {
		scanner := bufio.NewScanner(readers[i])
		scanner.Scan()

		scanners[i] = scanner
		lines[i] = scanner.Bytes()
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for {
		// Print current line combination
		for _, line := range lines {
			writer.Write(line)
		}
		writer.WriteByte('\n')

		// Update combination starting from the last file
		for i := readerCount - 1; i >= 0; i-- {
			if scanners[i].Scan() {
				lines[i] = scanners[i].Bytes()
				break
			} else {
				// Check if we've cycled through all combinations
				if i == 0 {
					return nil
				}

				// Reached the end of this file, rewind and continue with the next one
				if _, err := readers[i].Seek(0, 0); err != nil {
					return err
				}

				scanners[i] = bufio.NewScanner(readers[i])
				scanners[i].Scan()
				lines[i] = scanners[i].Bytes()
			}
		}
	}
}

func hasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}

	return true
}
