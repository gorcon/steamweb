package steamweb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newConfig(uri string) *Config {
	return &Config{
		Disabled:           false,
		Key:                "R1jamSsz17LHA9WgDW099YGfCs4fn0m0",
		URL:                uri,
		DefaultServerNames: []string{"My PZ Server"},
	}
}

func TestClient_GetPlayerBans(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintln(w, `{"players":[{"SteamId":"7656119","CommunityBanned":false,"VACBanned":true,"NumberOfVACBans":1,"DaysSinceLastBan":1530,"NumberOfGameBans":0,"EconomyBan":"none"}]}`)
	}))
	defer ts.Close()

	cfg := newConfig(ts.URL)

	steamID := "7656119"
	expect := []PlayerBans{
		{
			SteamID:          "7656119",
			CommunityBanned:  false,
			VACBanned:        true,
			NumberOfVACBans:  1,
			DaysSinceLastBan: 1530,
			NumberOfGameBans: 0,
			EconomyBan:       "none",
		},
	}

	if assert.Nil(t, cfg.Validate()) {
		client := NewClient(cfg)
		response, err := client.GetPlayerBans(steamID)

		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, expect, response)
		}
	}
}

func TestClient_GetServerList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintln(w, `{"response":{"servers":[{"addr":"127.0.0.1:16261","gameport":16261,"steamid":"90268762852129810","name":"My PZ Server","appid":108600,"gamedir":"zomboid","version":"1.0.0.0","product":"zomboid","region":-1,"players":2,"max_players":32,"bots":0,"secure":true,"dedicated":true,"os":"w","gametype":"hidden;hosted"},{"addr":"127.0.0.2:16267","gameport":16267,"steamid":"90268762793969688","name":"Super Server","appid":108600,"gamedir":"zomboid","version":"1.0.0.0","product":"zomboid","region":-1,"players":0,"max_players":10,"bots":0,"map":"vehicle_interior;SecretZ_v4;InG","secure":false,"dedicated":true,"os":"w"},{"addr":"127.0.0.3:16260","gameport":16260,"steamid":"90268200350011416","name":"Best Server","appid":108600,"gamedir":"zomboid","version":"1.0.0.0","product":"zomboid","region":-1,"players":9,"max_players":30,"bots":0,"map":"Muldraugh, KY","secure":true,"dedicated":true,"os":"l"},{"addr":"127.0.0.4:16261","gameport":16261,"steamid":"90268799310246930","name":"My PZ Server","appid":108600,"gamedir":"zomboid","version":"1.0.0.0","product":"zomboid","region":-1,"players":13,"max_players":32,"bots":0,"map":"Muldraugh, KY","secure":true,"dedicated":true,"os":"w","gametype":"hidden;hosted"},{"addr":"127.0.0.5:16261","gameport":16261,"steamid":"90268762155669518","name":"My PZ Server","appid":108600,"gamedir":"zomboid","version":"1.0.0.0","product":"zomboid","region":-1,"players":2,"max_players":2,"bots":0,"map":"Muldraugh, KY","secure":true,"dedicated":true,"os":"w","gametype":"hidden;hosted"}]}}`)
	}))
	defer ts.Close()

	cfg := newConfig(ts.URL)

	tests := []struct {
		name    string
		filter  *GetServerListFilter
		want    []Server
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "native filters",
			filter: &GetServerListFilter{},
			want: []Server{
				{Addr: "127.0.0.4:16261", GamePort: 16261, SteamID: "90268799310246930", Name: "My PZ Server", AppID: 108600, GameDir: "zomboid", Version: "1.0.0.0", Product: "zomboid", Region: -1, Players: 13, MaxPlayers: 32, Bots: 0, Map: "Muldraugh, KY", Secure: true, Dedicated: true, OS: "w", GameType: "hidden;hosted"},
				{Addr: "127.0.0.3:16260", GamePort: 16260, SteamID: "90268200350011416", Name: "Best Server", AppID: 108600, GameDir: "zomboid", Version: "1.0.0.0", Product: "zomboid", Region: -1, Players: 9, MaxPlayers: 30, Bots: 0, Map: "Muldraugh, KY", Secure: true, Dedicated: true, OS: "l", GameType: ""},
				{Addr: "127.0.0.1:16261", GamePort: 16261, SteamID: "90268762852129810", Name: "My PZ Server", AppID: 108600, GameDir: "zomboid", Version: "1.0.0.0", Product: "zomboid", Region: -1, Players: 2, MaxPlayers: 32, Bots: 0, Map: "", Secure: true, Dedicated: true, OS: "w", GameType: "hidden;hosted"},
				{Addr: "127.0.0.5:16261", GamePort: 16261, SteamID: "90268762155669518", Name: "My PZ Server", AppID: 108600, GameDir: "zomboid", Version: "1.0.0.0", Product: "zomboid", Region: -1, Players: 2, MaxPlayers: 2, Bots: 0, Map: "Muldraugh, KY", Secure: true, Dedicated: true, OS: "w", GameType: "hidden;hosted"},
				{Addr: "127.0.0.2:16267", GamePort: 16267, SteamID: "90268762793969688", Name: "Super Server", AppID: 108600, GameDir: "zomboid", Version: "1.0.0.0", Product: "zomboid", Region: -1, Players: 0, MaxPlayers: 10, Bots: 0, Map: "vehicle_interior;SecretZ_v4;InG", Secure: false, Dedicated: true, OS: "w", GameType: ""}},
			wantErr: assert.NoError,
		},
		{
			name:   "custom filters",
			filter: &GetServerListFilter{NoDefaultServers: true},
			want: []Server{
				{Addr: "127.0.0.3:16260", GamePort: 16260, SteamID: "90268200350011416", Name: "Best Server", AppID: 108600, GameDir: "zomboid", Version: "1.0.0.0", Product: "zomboid", Region: -1, Players: 9, MaxPlayers: 30, Bots: 0, Map: "Muldraugh, KY", Secure: true, Dedicated: true, OS: "l", GameType: ""},
				{Addr: "127.0.0.2:16267", GamePort: 16267, SteamID: "90268762793969688", Name: "Super Server", AppID: 108600, GameDir: "zomboid", Version: "1.0.0.0", Product: "zomboid", Region: -1, Players: 0, MaxPlayers: 10, Bots: 0, Map: "vehicle_interior;SecretZ_v4;InG", Secure: false, Dedicated: true, OS: "w", GameType: ""}},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(cfg)

			got, err := client.GetServerList(tt.filter)
			if !tt.wantErr(t, err, fmt.Sprintf("GetServerList(%v)", tt.filter)) {
				return
			}

			assert.Equalf(t, tt.want, got, "GetServerList(%v)", tt.filter)
		})
	}
}
