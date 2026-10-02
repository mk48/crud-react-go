package middleware

import (
	"testing"

	"kfamily/internal/util"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/golang-jwt/jwt/v4"
)

func testMiddleware(t *testing.T) *Middleware {
	t.Helper()
	clients, err := util.NewClientRegistry("web-id", `[{"name":"mobile","clientId":"mobile-id"}]`)
	if err != nil {
		t.Fatal(err)
	}
	return NewMiddleware(nil, &util.AppENV{CasdoorEndpoint: "https://door.example.com/"}, clients)
}

func claims(azp string, aud ...string) *casdoorsdk.Claims {
	c := &casdoorsdk.Claims{TokenType: "access-token", Azp: azp}
	c.Issuer = "https://door.example.com"
	c.Audience = jwt.ClaimStrings(aud)
	return c
}

func TestValidateClaimsResolvesClient(t *testing.T) {
	mw := testMiddleware(t)

	cases := map[string]struct {
		claims *casdoorsdk.Claims
		want   string
	}{
		"by azp":           {claims("mobile-id", "mobile-id"), "mobile"},
		"by audience":      {claims("", "web-id"), util.ClientWeb},
		"first registered": {claims("", "other", "mobile-id"), "mobile"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			client, err := mw.validateClaims(tc.claims)
			if err != nil || client.Name != tc.want {
				t.Errorf("validateClaims = %q, %v; want %q", client.Name, err, tc.want)
			}
		})
	}
}

func TestValidateClaimsRejects(t *testing.T) {
	mw := testMiddleware(t)

	refresh := claims("web-id", "web-id")
	refresh.TokenType = "refresh-token"
	otherIssuer := claims("web-id", "web-id")
	otherIssuer.Issuer = "https://evil.example.com"

	cases := map[string]*casdoorsdk.Claims{
		"refresh token":          refresh,
		"other issuer":           otherIssuer,
		"unregistered azp":       claims("other-app", "web-id"), // azp wins over aud
		"unregistered audience":  claims("", "other-app"),
		"no azp and no audience": claims(""),
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if client, err := mw.validateClaims(c); err == nil {
				t.Errorf("validateClaims accepted the token as %q", client.Name)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("héllo", 2); got != "hé" {
		t.Errorf("truncate = %q, want rune-safe %q", got, "hé")
	}
	if got := truncate("ok", 10); got != "ok" {
		t.Errorf("truncate = %q", got)
	}
}
