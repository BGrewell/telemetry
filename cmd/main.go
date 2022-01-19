package main

import (
	"flag"
	"fmt"
	"github.com/BGrewell/tgams/internal/engine"
	"github.com/BGrewell/tgams/internal/grpc"
	log "github.com/BGrewell/tgams/internal/logging"
	"github.com/BGrewell/tgams/internal/state"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	version   = "debug"
	tag       = "debug"
	commit    = "debug"
	branch    = "debug"
	showlines = true

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

	//TODO: Figure out how to allow colors to be used as parameter inputs. Currently the length of the ANSI control
	//      sequence is added to the length of the string and messes up the padding since the characters are counted
	//      but aren't printed to the screen.
	parameterFormat = "%s%-22s %-20v %s\n"

	flagKeys = []string{} // used to control the order of the sections when help is printed
	flagMap  = map[string][]*CommandLineFlag{}
)

// init sets up the various variables and structures used by the program.
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

// setupLogging sets up the logging for the program.
func setupLogging(logLevel *string, useJson *bool, defaultOutput *os.File, disableLines *bool) {
	// Parse log level
	level := logrus.ErrorLevel
	switch strings.ToLower(*logLevel) {
	case "panic":
		level = logrus.PanicLevel
	case "fatal":
		level = logrus.FatalLevel
	case "error":
		level = logrus.ErrorLevel
	case "warn":
	case "warning":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	case "trace":
		level = logrus.TraceLevel
	}

	// Setup formatter
	var formatter logrus.Formatter
	formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	if *useJson {
		formatter = &logrus.JSONFormatter{}
	}

	// Setup logging
	lines := false
	if showlines {
		lines = !*disableLines
	}
	log.Setup(level, defaultOutput, formatter, lines)
}

// RunDaemon executes the program in daemon mode
func RunDaemon(s state.DaemonState) error {
	// do any daemon specific stuff here
	return RunServer(s.ServerState)
}

// RunServer executes the program as a server
func RunServer(s state.ServerState) error {
	server := grpc.GetControlServer(s.ListenAddr, s.ListenPort)
	server.ServeAsync()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	server.Shutdown()

	return nil
}

// RunClient executes the program as a client
func RunClient(s state.ClientState) error {

	core, err := engine.NewCoreEngine(s.HostAddr, s.HostPort, s.Timeout)
	if err != nil {
		return err
	}

	for core.IsRunning() {
		time.Sleep(time.Second)
	}

	return nil
}

func main() {

	// Setup command line parsing
	defaultOutput := os.Stdout
	sectionGeneral := "General"
	sectionServer := "Server"
	sectionClient := "Client"
	sectionController := "Controller"

	// Add flags
	AddFlagBool("help", sectionGeneral, false, "Print this help message", &[]string{"h"})
	printVersion := AddFlagBool("version", sectionGeneral, false, "Print the version information", &[]string{"v"})
	runServer := AddFlagBool("server", sectionGeneral, false, "Run in server mode", &[]string{"s"})
	hostAddr := AddFlagString("client", sectionGeneral, "", "Run as a client and connect to the service specified", &[]string{"c"})
	runDaemon := AddFlagBool("daemon", sectionGeneral, false, "Run as a service", &[]string{"d"})
	logLevel := AddFlagString("log-level", sectionGeneral, "info", "Set the logging level", &[]string{"ll"})
	useJson := AddFlagBool("json", sectionGeneral, false, "Use JSON log output", &[]string{"j"})
	port := AddFlagInt("port", sectionGeneral, 9550, "Control port", &[]string{"p"})
	defaults := AddFlagString("defaults", sectionGeneral, "~/.tgams.yaml", "Load defaults from file", nil)
	linesDisabled := AddFlagBool("hide-lines", sectionGeneral, false, "Hide lines in output", &[]string{"hl"})

	// Client specific flags
	timeout := AddFlagInt("timeout", []string{sectionClient, sectionController}, 5, "Connection Timeout in seconds", &[]string{"t"})

	// Server specific flags
	listenAddr := AddFlagString("listen", sectionServer, "0.0.0.0", "Listen address", &[]string{"l"})

	// Controller specific flags
	remoteAddr := AddFlagString("remote", sectionController, "", "Daemon ip address", &[]string{"r"})

	// Parse command line
	flag.Usage = Usage(defaultOutput)
	flag.Parse()

	// Print version information and exit
	if *printVersion {
		PrintAppHeader(defaultOutput)
		return
	}

	// Setup logging
	setupLogging(logLevel, useJson, defaultOutput, linesDisabled)

	// Ensure either client, server or daemon mode is specified
	if !*runServer && *hostAddr == "" && !*runDaemon {
		log.ErrorWithFields(map[string]interface{}{"server": *runServer, "client": *hostAddr, "daemon": *runDaemon}, "must specify either server, client or daemon mode")
		flag.Usage()
		return
	}

	if *runDaemon {
		log.Info("running in daemon mode")
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
		log.Info("running in server mode")
		s := state.ServerState{
			ListenAddr: *listenAddr,
			ListenPort: *port,
		}
		err := RunServer(s)
		if err != nil {
			panic(err)
		}
	}

	if *hostAddr != "" {
		log.Info("running in client mode")
		s := state.ClientState{
			HostAddr: *hostAddr,
			HostPort: *port,
			Timeout:  *timeout,
		}
		err := RunClient(s)
		if err != nil {
			panic(err)
		}
	}

	// TODO: Temp absorb flags
	if *hostAddr != "" && *port == 9550 && sectionServer != "" && sectionClient != "" && *listenAddr != "" && *defaults != "" &&
		*remoteAddr != "" {
		fmt.Println("")
	}

}
