package http

import (
	"effective-mobile-task/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type сlient struct {
	TargetHost string
	TargetPort string
}

func NewClient(host, port string) *сlient {
	return &сlient{host, port}
}

func (c *сlient) GetSongInfo(group, song string) (*models.SongInfo, error) {
	url := fmt.Sprintf("http://%s:%s/info?group=%s&song=%s", c.TargetHost, c.TargetPort, group, song)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get song info: %v", resp.StatusCode)
	}

	var info models.SongInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}
