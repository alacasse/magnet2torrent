package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"magnet2torrent/internal/config"
	"magnet2torrent/internal/logging"
	"magnet2torrent/internal/qbclient"
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

	if usedDefaults || needsQBConfig(cfg) {
		if !isInteractive() {
			logger.Errorf("config missing and no TTY available; create %s manually with qbHost/qbUsername/qbPassword", configPath)
			os.Exit(1)
		}
		if err := promptAndSaveConfig(configPath, cfg, logger); err != nil {
			logger.Errorf("failed to save config: %v", err)
			os.Exit(1)
		}
	}

	args := flag.Args()
	magnet := "<none provided>"
	if len(args) > 0 {
		magnet = args[0]
		if err := processMagnet(magnet, cfg, logger); err != nil {
			logger.Errorf("failed to process magnet: %v", err)
			os.Exit(1)
		}
	}

	fmt.Printf("magnet2torrent wired and running\n")
	fmt.Printf("  version     : %s\n", version)
	fmt.Printf("  config path : %s\n", configPath)
	fmt.Printf("  used defaults: %t\n", usedDefaults)
	fmt.Printf("  save dir    : %s\n", cfg.SaveDir)
	fmt.Printf("  log level   : %s\n", cfg.LogLevel)
	fmt.Printf("  magnet arg  : %s\n", magnet)
}

type qbClient interface {
	Login() error
	AddMagnet(string) error
}

var qbClientFactory = func(cfg *config.Config) qbClient {
	return qbclient.New(cfg.QbHost, cfg.QbUsername, cfg.QbPassword)
}

func processMagnet(magnetLink string, cfg *config.Config, logger *logging.Logger) error {
	if err := validateQBConfig(cfg); err != nil {
		return err
	}

	qb := qbClientFactory(cfg)

	if err := qb.Login(); err != nil {
		return fmt.Errorf("qbittorrent login failed: %w", err)
	}

	if err := qb.AddMagnet(magnetLink); err != nil {
		return fmt.Errorf("could not send magnet to qbittorrent: %w", err)
	}

	logger.Infof("magnet forwarded to qBittorrent at %s", cfg.QbHost)
	return nil
}

func validateQBConfig(cfg *config.Config) error {
	if cfg.QbHost == "" {
		return errors.New("qbittorrent host is empty; set qbHost in config")
	}
	if cfg.QbUsername == "" {
		return errors.New("qbittorrent username is empty; set qbUsername in config")
	}
	if cfg.QbPassword == "" {
		return errors.New("qbittorrent password is empty; set qbPassword in config")
	}
	return nil
}

func needsQBConfig(cfg *config.Config) bool {
	return cfg.QbHost == "" || cfg.QbUsername == "" || cfg.QbPassword == ""
}

func promptAndSaveConfig(configPath string, cfg *config.Config, logger *logging.Logger) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Config not found or incomplete. Please provide qBittorrent settings.\n")
	cfg.QbHost = promptValue(reader, "qBittorrent host (e.g. http://localhost:8080)", cfg.QbHost)
	cfg.QbUsername = promptValue(reader, "qBittorrent username", cfg.QbUsername)
	cfg.QbPassword = promptValue(reader, "qBittorrent password", cfg.QbPassword)

	if err := validateQBConfig(cfg); err != nil {
		return err
	}

	if err := config.SaveConfig(configPath, cfg); err != nil {
		return err
	}

	logger.Infof("config written to %s", configPath)
	return nil
}

func promptValue(reader *bufio.Reader, label string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultVal
	}
	return text
}

func isInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
