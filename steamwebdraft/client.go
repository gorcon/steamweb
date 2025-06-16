package steamweb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
)

const (
	GetPlayerBansURL = "/ISteamUser/GetPlayerBans/v1?key=%s&steamids=%s"
	GetServerListURL = "/IGameServersService/GetServerList/v1?key=%s&limit=%d&filter=%s"
)

var (
	ErrWrongStatusCode = errors.New("wrong status code")
	ErrEmptyResponse   = errors.New("empty response")
)

// Client is http client for getting requests to ISteamUser api.
type Client struct {
	config *Config
	http   *http.Client
}

// NewClient creates and returns a new Client instance initialized with the provided configuration.
func NewClient(cfg *Config) *Client {
	cfg.SetDefaults()

	return &Client{
		config: cfg,
		http: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:       cfg.Transport.Dialer.Timeout,
					Deadline:      cfg.Transport.Dialer.Deadline,
					FallbackDelay: cfg.Transport.Dialer.FallbackDelay,
					KeepAlive:     cfg.Transport.Dialer.KeepAlive,
				}).DialContext,
				TLSHandshakeTimeout: cfg.Transport.TLSHandshakeTimeout,
			},
		},
	}
}

// GetPlayerBans returns Community, VAC, and Economy ban statuses for given players.
// Example URL: http://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key=XXXXXXXXXXXXXXXXX&steamids=XXXXXXXX,YYYYY
func (c *Client) GetPlayerBans(steamIDs ...string) ([]PlayerBans, error) {
	response := GetPlayerBansResponse{}

	// Return empty ban history with disabled client.
	if c.config.Disabled {
		return response.Players, nil
	}

	uri := c.config.URL + fmt.Sprintf(GetPlayerBansURL, c.config.Key, strings.Join(steamIDs, ","))

	body, err := c.sendRequest(context.Background(), http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Players, nil
}

// GetServerList returns Steam servers from filter query.
// Example URL: http://api.steampowered.com/IGameServersService/GetServerList/v1/?key=XXXXXXXXXXXXXXXXX&limit=X&filter=F
func (c *Client) GetServerList(filter *GetServerListFilter) ([]Server, error) {
	response := GetServerListResponse{}

	// Return empty servers list with disabled client.
	if c.config.Disabled {
		return response.Response.Servers, nil
	}

	limit := filter.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	uri := c.config.URL + fmt.Sprintf(GetServerListURL, c.config.Key, limit, filter.String())

	body, err := c.sendRequest(context.Background(), http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return c.filterServers(response.Response.Servers, filter), nil
}

func (c *Client) sendRequest(ctx context.Context, method, uri string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if res != nil {
		defer res.Body.Close()
	} else {
		return nil, ErrEmptyResponse
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d %s", ErrWrongStatusCode, res.StatusCode, res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func (c *Client) filterServers(servers []Server, filter *GetServerListFilter) []Server {
	if filter.NoHidden || filter.NoDefaultServers {
		removeAddrs := make(map[string]bool)

		for i := range servers {
			server := &servers[i]

			if filter.NoHidden && strings.Contains(server.GameType, "hidden") {
				removeAddrs[server.Addr] = true

				continue
			}

			if filter.NoDefaultServers {
				for _, name := range c.config.DefaultServerNames {
					if server.Name == name {
						removeAddrs[server.Addr] = true

						continue
					}
				}
			}
		}

		servers = c.removeFilteredServers(servers, removeAddrs)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Players > servers[j].Players
	})

	return servers
}

func (c *Client) removeFilteredServers(servers []Server, addrs map[string]bool) []Server {
	result := make([]Server, 0, len(servers))

	for i := range servers {
		if _, ok := addrs[servers[i].Addr]; !ok {
			result = append(result, servers[i])
		}
	}

	return result
}
