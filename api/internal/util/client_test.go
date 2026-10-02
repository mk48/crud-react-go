package util

import (
	"strings"
	"testing"
)

func TestNewClientRegistry(t *testing.T) {
	r, err := NewClientRegistry("web-id", `[
		{"name":"mobile","clientId":"mobile-id"},
		{"name":"batch:nightly-sync","clientId":"batch-id","serviceUserId":"00000000-0000-0000-0000-000000000002"}
	]`)
	if err != nil {
		t.Fatal(err)
	}

	for clientID, want := range map[string]string{"web-id": ClientWeb, "mobile-id": "mobile", "batch-id": "batch:nightly-sync"} {
		c, ok := r.Lookup(clientID)
		if !ok || c.Name != want {
			t.Errorf("Lookup(%q) = %q, %v; want %q", clientID, c.Name, ok, want)
		}
	}

	if c, _ := r.Lookup("mobile-id"); c.IsService() {
		t.Error("mobile should be a user client")
	}
	if c, _ := r.Lookup("batch-id"); !c.IsService() {
		t.Error("batch:nightly-sync should be a service client")
	}
	if _, ok := r.Lookup("unknown"); ok {
		t.Error("unknown client id must not resolve")
	}
}

func TestNewClientRegistryWebOnly(t *testing.T) {
	r, err := NewClientRegistry("web-id", "")
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := r.Lookup("web-id"); !ok || c.Name != ClientWeb {
		t.Errorf("web client not registered: %+v %v", c, ok)
	}
}

func TestNewClientRegistryRejects(t *testing.T) {
	cases := map[string]string{
		"invalid json":      `{"name":"mobile"}`,
		"reserved name":     `[{"name":"system","clientId":"x"}]`,
		"duplicate name":    `[{"name":"web","clientId":"x"}]`,
		"duplicate id":      `[{"name":"mobile","clientId":"web-id"}]`,
		"missing client id": `[{"name":"mobile"}]`,
		"bad name":          `[{"name":"Mobile App","clientId":"x"}]`,
		"empty name":        `[{"name":"","clientId":"x"}]`,
		"bad service user":  `[{"name":"batch:x","clientId":"x","serviceUserId":"not-a-uuid"}]`,
	}
	for name, apiClients := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := NewClientRegistry("web-id", apiClients); err == nil {
				t.Errorf("NewClientRegistry(%s) succeeded, want error", apiClients)
			}
		})
	}
}

func TestClientNameLength(t *testing.T) {
	if !clientNamePattern.MatchString(strings.Repeat("a", 100)) {
		t.Error("100 characters should fit operation.client")
	}
	if clientNamePattern.MatchString(strings.Repeat("a", 101)) {
		t.Error("101 characters must not fit operation.client varchar(100)")
	}
}
