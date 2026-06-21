package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var dockerClient = &http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", "/var/run/docker.sock")
		},
	},
}

type dockerContainer struct {
	Names  []string `json:"Names"`
	State  string   `json:"State"`
	Status string   `json:"Status"`
}

func collectContainers(project string) []ContainerStatus {
	endpoint := "http://localhost/containers/json?all=1"
	if project != "" {
		f := fmt.Sprintf(`{"label":[%q]}`, "com.docker.compose.project="+project)
		endpoint += "&filters=" + url.QueryEscape(f)
	}

	resp, err := dockerClient.Get(endpoint)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var containers []dockerContainer
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil
	}

	result := make([]ContainerStatus, len(containers))
	for i, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		result[i] = ContainerStatus{
			Name:   name,
			State:  c.State,
			Status: c.Status,
		}
	}
	return result
}
