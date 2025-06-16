package steamweb

type (
	// GetPlayerBansResponse describes response for Steam GetPlayerBans request.
	GetPlayerBansResponse struct {
		Players []PlayerBans `json:"players"`
	}

	// PlayerBans is list of player ban objects for each 64 bit ID requested.
	PlayerBans struct {
		// SteamId (string) The player's 64 bit ID.
		SteamID string `json:"SteamId"`

		// CommunityBanned (bool) Indicates whether or not the player is banned from Steam Community.
		CommunityBanned bool `json:"CommunityBanned"`

		// VACBanned (bool) Indicates whether or not the player has VAC bans on record.
		VACBanned bool `json:"VACBanned"`

		// NumberOfVACBans (int) Number of VAC bans on record.
		NumberOfVACBans int `json:"NumberOfVACBans"`

		// DaysSinceLastBan (int) Number of days since the last ban.
		DaysSinceLastBan int `json:"DaysSinceLastBan"`

		// NumberOfGameBans (int) Number of bans in games, this includes CS:GO Overwatch bans.
		NumberOfGameBans int `json:"NumberOfGameBans"`

		// EconomyBan (string) The player's ban status in the economy.
		// If the player has no bans on record the string will be "none",
		// if the player is on probation it will say "probation", etc.
		EconomyBan string `json:"EconomyBan"`
	}
)

type (
	GetServerListResponse struct {
		Response struct {
			Servers []Server `json:"servers"`
		} `json:"response"`
	}

	Server struct {
		Addr       string `json:"addr"`
		GamePort   int    `json:"gameport"`
		SteamID    string `json:"steamid"`
		Name       string `json:"name"`
		AppID      int    `json:"appid"`
		GameDir    string `json:"gamedir"`
		Version    string `json:"version"`
		Product    string `json:"product"`
		Region     int    `json:"region"`
		Players    int    `json:"players"`
		MaxPlayers int    `json:"max_players"`
		Bots       int    `json:"bots"`
		Map        string `json:"map"`
		Secure     bool   `json:"secure"`
		Dedicated  bool   `json:"dedicated"`
		OS         string `json:"os"`
		GameType   string `json:"gametype"`
	}
)
