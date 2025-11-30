package config

import (
	"path/filepath"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		goos    string
		home    string
		appdata string
		want    string
	}{
		{
			name: "linux_default",
			goos: "linux",
			home: "/home/alice",
			want: filepath.Join("/home/alice", ".config", "magnet2torrent", "config.json"),
		},
		{
			name:    "windows_appdata_present",
			goos:    "windows",
			appdata: `C:\Users\alice\AppData\Roaming`,
			want:    filepath.Join(`C:\Users\alice\AppData\Roaming`, "magnet2torrent", "config.json"),
		},
		{
			name: "windows_appdata_missing_fallback_home",
			goos: "windows",
			home: `C:\Users\alice`,
			want: filepath.Join(`C:\Users\alice`, "AppData", "Roaming", "magnet2torrent", "config.json"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := defaultConfigPath(tc.goos, tc.home, tc.appdata)
			if got != tc.want {
				t.Fatalf("defaultConfigPath(%q, %q, %q) = %q, want %q", tc.goos, tc.home, tc.appdata, got, tc.want)
			}
		})
	}
}
