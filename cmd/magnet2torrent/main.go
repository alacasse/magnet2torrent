package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"magnet2torrent/internal/config"
	"magnet2torrent/internal/logging"
)

const version = "0.1.0"

func main() {
	defaultConfigPath := config.GetDefaultConfigPath()

	var (
		configPathFlag = flag.String("config", defaultConfigPath, "path to config file")
		versionFlag    = flag.Bool("version", false, "print version and exit")
		versionShort   = flag.Bool("v", false, "print version and exit (shorthand)")
	)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "magnet2torrent - placeholder CLI for magnet handling\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n  magnet2torrent [flags] [magnet]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nDefault config path: %s\n", defaultConfigPath)
	}

	flag.Parse()

	if *versionFlag || *versionShort {
		fmt.Printf("magnet2torrent %s\n", version)
		return
	}

	logger := logging.NewLogger("info")

	configPath := filepath.Clean(*configPathFlag)
	cfg, usedDefaults, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Errorf("failed to load config: %v", err)
		os.Exit(1)
	}

	if usedDefaults {
		logger.Infof("config not found at %s; using defaults", configPath)
	}

	args := flag.Args()
	magnet := "<none provided>"
	if len(args) > 0 {
		magnet = args[0]
	}

	fmt.Printf("magnet2torrent wired and running\n")
	fmt.Printf("  version     : %s\n", version)
	fmt.Printf("  config path : %s\n", configPath)
	fmt.Printf("  used defaults: %t\n", usedDefaults)
	fmt.Printf("  save dir    : %s\n", cfg.SaveDir)
	fmt.Printf("  log level   : %s\n", cfg.LogLevel)
	fmt.Printf("  magnet arg  : %s\n", magnet)
}
