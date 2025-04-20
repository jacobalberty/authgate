package peers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	peersEndpoint = "/api/peers"
)

type Client interface {
	GetPeersByIP(ctx context.Context, ip string) ([]Peer, error)
	IsPeerInGroup(ctx context.Context, group, ip string) (bool, error)
}

func New(rootEP, token string) (Client, error) {
	// Placeholder for client initialization logic
	// This could involve setting up HTTP clients, authentication, etc.
	return &client{
		rootEndpoint: rootEP,
		token:        token,
		client:       &http.Client{},
	}, nil
}

type client struct {
	rootEndpoint string
	token        string
	client       *http.Client
}

func (c *client) GetPeersByIP(ctx context.Context, ip string) ([]Peer, error) {
	u, err := url.Parse(c.rootEndpoint + peersEndpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("ip", ip)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Accept", "application/json")
	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get peers: %s", resp.Status)
	}
	var peers []Peer
	if err := json.NewDecoder(resp.Body).Decode(&peers); err != nil {
		return nil, err
	}
	return peers, nil

}

func (c *client) IsPeerInGroup(ctx context.Context, group, ip string) (bool, error) {
	peers, err := c.GetPeersByIP(ctx, ip)
	if err != nil {
		return false, err
	}
	if len(peers) != 1 {
		return false, nil
	}
	for _, g := range peers[0].Groups {
		if g.Name == group {
			return true, nil
		}
	}
	return false, nil
}
