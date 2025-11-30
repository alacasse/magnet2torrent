package main

import (
	"errors"
	"strings"
	"testing"

	"magnet2torrent/internal/config"
	"magnet2torrent/internal/logging"
)

type stubQBClient struct {
	loginErr   error
	addErr     error
	lastMagnet string
}

func (s *stubQBClient) Login() error {
	return s.loginErr
}

func (s *stubQBClient) AddMagnet(magnet string) error {
	s.lastMagnet = magnet
	return s.addErr
}

func TestProcessMagnetSuccess(t *testing.T) {
	origFactory := qbClientFactory
	defer func() { qbClientFactory = origFactory }()

	stub := &stubQBClient{}
	qbClientFactory = func(cfg *config.Config) qbClient {
		return stub
	}

	cfg := &config.Config{
		QbHost:     "http://example.test",
		QbUsername: "admin",
		QbPassword: "password",
	}
	magnet := "magnet:?xt=urn:btih:example"

	logger := logging.NewLogger("info")
	if err := processMagnet(magnet, cfg, logger); err != nil {
		t.Fatalf("processMagnet returned error: %v", err)
	}

	if stub.lastMagnet != magnet {
		t.Fatalf("expected magnet %s, got %s", magnet, stub.lastMagnet)
	}
}

func TestProcessMagnetLoginError(t *testing.T) {
	origFactory := qbClientFactory
	defer func() { qbClientFactory = origFactory }()

	stub := &stubQBClient{loginErr: errors.New("login failed")}
	qbClientFactory = func(cfg *config.Config) qbClient { return stub }

	logger := logging.NewLogger("info")
	err := processMagnet("magnet:?xt=urn:btih:example", &config.Config{
		QbHost:     "http://example.test",
		QbUsername: "admin",
		QbPassword: "password",
	}, logger)
	if err == nil || !strings.Contains(err.Error(), "login failed") {
		t.Fatalf("expected login error, got %v", err)
	}
}

func TestProcessMagnetAddError(t *testing.T) {
	origFactory := qbClientFactory
	defer func() { qbClientFactory = origFactory }()

	stub := &stubQBClient{addErr: errors.New("add failed")}
	qbClientFactory = func(cfg *config.Config) qbClient { return stub }

	logger := logging.NewLogger("info")
	err := processMagnet("magnet:?xt=urn:btih:example", &config.Config{
		QbHost:     "http://example.test",
		QbUsername: "admin",
		QbPassword: "password",
	}, logger)
	if err == nil || !strings.Contains(err.Error(), "add failed") {
		t.Fatalf("expected add error, got %v", err)
	}
}

func TestValidateQBConfig(t *testing.T) {
	cases := []struct {
		name    string
		cfg     config.Config
		wantErr string
	}{
		{name: "missing host", cfg: config.Config{QbUsername: "u", QbPassword: "p"}, wantErr: "qbittorrent host is empty"},
		{name: "missing user", cfg: config.Config{QbHost: "http://h", QbPassword: "p"}, wantErr: "qbittorrent username is empty"},
		{name: "missing pass", cfg: config.Config{QbHost: "http://h", QbUsername: "u"}, wantErr: "qbittorrent password is empty"},
		{name: "ok", cfg: config.Config{QbHost: "http://h", QbUsername: "u", QbPassword: "p"}, wantErr: ""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validateQBConfig(&tc.cfg)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
				}
			}
		})
	}
}
