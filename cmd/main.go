package main

import (
	"flag"
	"fmt"
	"github.com/BGrewell/tgams/internal/state"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	version = "debug"
	tag     = "debug"
	commit  = "debug"
	branch  = "debug"

	yellow    = color.New(color.FgHiYellow).FprintfFunc()
	red       = color.New(color.FgHiRed).FprintfFunc()
	green     = color.New(color.FgHiGreen).FprintfFunc()
	cyan      = color.New(color.FgHiCyan).FprintfFunc()
	white     = color.New(color.FgHiWhite).FprintfFunc()
	darkWhite = color.New(color.FgWhite).FprintfFunc()
	blue      = color.New(color.FgHiBlue).FprintfFunc()
	magenta   = color.New(color.FgHiMagenta).FprintfFunc()

	y  = color.New(color.FgHiYellow).SprintFunc()
	r  = color.New(color.FgHiRed).SprintFunc()
	g  = color.New(color.FgHiGreen).SprintFunc()
	c  = color.New(color.FgHiCyan).SprintFunc()
	w  = color.New(color.FgHiWhite).SprintFunc()
	dw = color.New(color.FgWhite).SprintFunc()
	b  = color.New(color.FgHiBlue).SprintFunc()
	m  = color.New(color.FgHiMagenta).SprintFunc()

	level1 = "    "
	level2 = "      "
	level3 = "        "

	parameterFormat = "%s%-22s %-20v %s\n"

	flagKeys = []string{} // used to control the order of the sections when help is printed
	flagMap = map[string][]*CommandLineFlag{}
)

func init() {
	flagKeys = make([]string, 0)
	flagMap = make(map[string][]*CommandLineFlag)
}

// PrintAppHeader prints the application header
func PrintAppHeader(output io.Writer) {
	white(output, "[+] Traffic Generation Analysis and Measuremenet System\n")
	white(output, "%sVersion: %s [ Branch: %s | Commit: %s | Tag: %s ]\n", level1, y(version), dw(branch), dw(commit), dw(tag))
}

// PrintUsageSectionHeader prints a section header
func PrintUsageSectionHeader(output io.Writer, section string) {
	blue(output, "\n%s%s:\n", level1, section)
	white(output, parameterFormat, level1, "Parameter", "Default", "Description")
}

// PrintUsageLine prints the usage information for the given parameter
func PrintUsageLine(output io.Writer, parameter string, defaultValue interface{}, description string) {
	switch defaultValue.(type) {
	case string:
		if defaultValue.(string) == "" {
			defaultValue = "\"\""
		}
	}
	darkWhite(output, parameterFormat, level1, parameter, defaultValue, description)
}

// Usage automatically generates the usage text in a more readable format than the default flag.Usage() the
// flags package provides.
func Usage(output io.Writer) (usage func()) {
	return func() {
		PrintAppHeader(output)
		for _, key := range flagKeys {
			PrintUsageSectionHeader(output, fmt.Sprintf("%s Options", key))
			for _, f := range flagMap[key] {
				p := []string{f.Name}
				if f.AltNames != nil {
					p = append(p, *f.AltNames...)
				}
				parameter := strings.Join(p, ", ")
				PrintUsageLine(output, parameter, f.Value, f.Usage)
			}
		}
		fmt.Println("")
	}
}

func main() {

	// Setup command line parsing
	defaultOutput := os.Stdout
	sectionGeneral := "General"
	sectionServer := "Server"
	sectionClient := "Client"


	// Add flags
	AddFlagBool("help", sectionGeneral, false, "Print this help message", &[]string{"h"})
	printVersion := AddFlagBool("version", sectionGeneral, false, "Print the version information", &[]string{"v"})
	runServer := AddFlagBool("server", sectionGeneral, false, "Run in server mode", &[]string{"s"})
	runClient := AddFlagBool("client", sectionGeneral, false, "Run in client mode", &[]string{"c"})
	runDaemon := AddFlagBool("daemon", sectionGeneral, false, "Run as a service", &[]string{"d"})
	runDebug := AddFlagBool("debug", sectionGeneral, false, "Run in debug mode", nil)
	port := AddFlagInt("port", sectionGeneral, 9550, "Control port", &[]string{"p"})
	defaults := AddFlagString("defaults", sectionGeneral, "~/.tgams.yaml", "Load defaults from file", nil)

	// Client specific flags
	hostAddr := AddFlagString("host", sectionClient, "", "Server ip address", &[]string{"H"})

	// Server specific flags
	listenAddr := AddFlagString("listen", sectionServer, "0.0.0.0", "Listen address", &[]string{"l"})


	// Parse command line
	flag.Usage = Usage(defaultOutput)
	flag.Parse()

	// Print version information and exit
	if *printVersion {
		PrintAppHeader(defaultOutput)
		return
	}

	// If debug is enabled then provide additional debug output
	if *runDebug {
        fmt.Println("[!] Enable debug output")
    }

	// Ensure either client, server or daemon mode is specified
	if !*runServer && !*runClient && !*runDaemon {
        red(defaultOutput, "Error: Must specify either client, server or daemon mode\n")
        flag.Usage()
        return
    }

	if *runDaemon {
        fmt.Println("[!] Running as a service")
		s := state.DaemonState{
			ServerState: state.ServerState{
                ListenAddr: *listenAddr,
				ListenPort: *port,
            },
		}
		err := RunDaemon(s)
		if err != nil {
			panic(err)
		}
    }

	if *runServer {
		fmt.Println("[!] Running in server mode")
        s := state.ServerState{
            ListenAddr: *listenAddr,
            ListenPort: *port,
        }
        err := RunServer(s)
        if err != nil {
            panic(err)
        }
	}

	if *runClient {
		fmt.Println("[!] Running in client mode")
        s := state.ClientState{
            HostAddr: *hostAddr,
            HostPort: *port,
        }
        err := RunClient(s)
        if err != nil {
            panic(err)
        }
	}

	// TODO: Temp absorb flags
	if *hostAddr != "" && *port == 9550 && sectionServer != "" && *listenAddr != "" && *defaults != "" {
        fmt.Println("")
    }
}

func RunDaemon(s state.DaemonState) error {
	return nil
}

func RunServer(s state.ServerState) error {
	return nil
}

func RunClient(s state.ClientState) error {
	return nil
}