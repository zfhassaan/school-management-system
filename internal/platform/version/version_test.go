package version

import "testing"

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		local  string
		remote string
		want   bool
	}{
		{"1.0.0", "1.0.1", true},
		{"1.0.0", "v1.0.1", true},
		{"v1.0.0", "v1.0.0", false},
		{"1.0.0", "1.0.0", false},
		{"1.2.0", "1.1.9", false},
		{"1.0.0", "2.0.0", true},
		{"1.0.0", "1.10.0", true},
		{"1.9.0", "1.10.0", true},
		{"1.10.0", "1.9.0", false},
	}

	for _, tt := range tests {
		got := IsNewerVersion(tt.local, tt.remote)
		if got != tt.want {
			t.Errorf("IsNewerVersion(%q, %q) = %v, want %v", tt.local, tt.remote, got, tt.want)
		}
	}
}
