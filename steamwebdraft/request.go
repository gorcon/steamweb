package steamweb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrRequiredParam = errors.New("param is required")

// GetServerListFilter represents the filter parameters used when querying game servers
// from the Steam server browser. Each field corresponds to a specific filter that can
// be applied to narrow down the server list results.
// The JSON tags match the actual query parameters used in the Steam server browser protocol.
//
// See: https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol.
type GetServerListFilter struct {
	// NotOr is a special filter, specifies that servers matching any of the following [x]
	// conditions should not be returned.
	// Usage: \nor\[x].
	NotOr []string `json:"nor,omitempty"`
	// NotAnd is a special filter, specifies that servers matching all of the following [x]
	// conditions should not be returned.
	// Usage: \nand\[x].
	NotAnd []string `json:"nand,omitempty"`
	// Dedicated is a filter for servers running dedicated.
	// Usage: \dedicated\1.
	Dedicated bool `json:"dedicated,omitempty"`
	// Secure is a filter for servers using anti-cheat technology (VAC, but potentially others as well).
	// Usage: \secure\1.
	Secure bool `json:"secure,omitempty"`
	// GameDir is a filter for servers running the specified modification (ex. cstrike).
	// Usage: \gamedir\[mod].
	GameDir string `json:"gamedir,omitempty"`
	// Map is a filter for servers running the specified map (ex. cs_italy).
	// Usage: \map\[map].
	Map string `json:"map,omitempty"`
	// Linux is a filter for servers running on a Linux platform.
	// Usage: \linux\1.
	Linux bool `json:"linux,omitempty"`
	// NoPassword is a filer for servers that are not password protected.
	// Usage: \password\0.
	NoPassword bool `json:"password,omitempty"`
	// NotEmpty is a filter for servers that are not empty.
	// Usage: \empty\1.
	NotEmpty bool `json:"empty,omitempty"`
	// NotFull is a filter for servers that are not full.
	// Usage: \full\1.
	NotFull bool `json:"full,omitempty"`
	// Proxy is a filter for servers that are spectator proxies.
	// Usage: \proxy\1.
	Proxy bool `json:"proxy,omitempty"`
	// AppID is a filer for servers that are running game [appid].
	// Usage: \appid\[appid].
	AppID int `json:"appid,omitempty"`
	// NotAppID is a filter for servers that are NOT running game [appid].
	// This was introduced to block Left 4 Dead games from the Steam Server Browser.
	// Usage: \napp\[appid].
	NotAppID int `json:"napp,omitempty"`
	// NoPlayers is a filer for servers that are empty.
	// Usage: \noplayers\1.
	NoPlayers bool `json:"noplayers,omitempty"`
	// Whitelisted is a filter for servers that are whitelisted.
	// Usage: \white\1.
	Whitelisted bool `json:"white,omitempty"`
	// GameTypeTags is a filer for servers with all of the given tag(s) in sv_tags.
	// Usage: \gametype\[tag,…].
	GameTypeTags []string `json:"gametype,omitempty"`
	// GameDataTags is a filer for servers with all of the given tag(s) in their ‘hidden’ tags (L4D2).
	// Usage: \gamedata\[tag,…].
	GameDataTags []string `json:"gamedata,omitempty"`
	// GameDataOrTags is a filer for servers with any of the given tag(s) in their ‘hidden’ tags (L4D2).
	// Usage: \gamedataor\[tag,…]
	GameDataOrTags []string `json:"gamedataor,omitempty"`
	// NameMatch is a filer for servers with their hostname matching [hostname] (can use * as a wildcard).
	// Usage: \name_match\[hostname].
	NameMatch string `json:"name_match,omitempty"`
	// VersionMatch is a filer for servers running version [version] (can use * as a wildcard).
	// Usage: \version_match\[version].
	VersionMatch string `json:"version_match,omitempty"`
	// CollapseAddrHash is a filer for that returns only one server for each unique IP address matched.
	// Usage: \collapse_addr_hash\1.
	CollapseAddrHash bool `json:"collapse_addr_hash,omitempty"`
	// GameAddr is a filer for that returns only servers on the specified IP address (port supported and optional).
	// Usage: \gameaddr\[ip].
	GameAddr string `json:"gameaddr,omitempty"`

	// Custom filters - not provided by Steam.

	// NoHidden is a custom filer for servers that are not hidden.
	NoHidden bool `json:"nohidden,omitempty"`
	// NoDefaultServers is a custom filer for servers that has a name that differs from the default name.
	NoDefaultServers bool `json:"no_default_servers,omitempty"`
	// Limit limits response.
	Limit int `json:"limit,omitempty"`
}

// String converts fields to url part with params.
func (g *GetServerListFilter) String() string { //nolint:funlen,cyclop // I don't care
	query := `\appid\` + strconv.Itoa(g.AppID)

	if g.Dedicated {
		query += `\dedicated\1`
	}

	if g.Secure {
		query += `\secure\1`
	}

	if g.GameDir != "" {
		query += `\gamedir\` + g.GameDir
	}

	// Not working in PZ.
	if g.Map != "" {
		query += `\map\` + g.Map
	}

	if g.Linux {
		query += `\linux\1`
	}

	// Not working in PZ.
	if g.NoPassword {
		query += `\password\0`
	}

	if g.NotEmpty {
		query += `\empty\1`
	}

	if g.NotFull {
		query += `\full\1`
	}

	// Not working in PZ.
	if g.Proxy {
		query += `\proxy\1`
	}

	// Not working in PZ.
	if g.NotAppID != 0 {
		query += `\napp\` + strconv.Itoa(g.NotAppID)
	}

	if g.NoPlayers {
		query += `\noplayers\1`
	}

	// Not working in PZ.
	if g.Whitelisted {
		query += `\white\1`
	}

	if len(g.GameTypeTags) != 0 {
		query += `\gametype\` + strings.Join(g.GameTypeTags, `;`)
	}

	// Not working in PZ.
	if len(g.GameDataTags) != 0 {
		query += `\gamedata\` + strings.Join(g.GameDataTags, `,`)
	}

	// Not working in PZ.
	if len(g.GameDataOrTags) != 0 {
		query += `\gamedataor\` + strings.Join(g.GameDataOrTags, `,`)
	}

	if g.NameMatch != "" {
		query += `\name_match\*` + g.NameMatch + `*`
	}

	return query
}

func (g *GetServerListFilter) Validate() error {
	if g.AppID == 0 {
		return fmt.Errorf("%w: %s", ErrRequiredParam, "appid")
	}

	return nil
}
