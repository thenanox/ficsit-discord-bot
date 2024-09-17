package satisfactory

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func QueryServerState() (*QueryServerStateOutput, error) {
	response := &QueryServerStateOutput{}
	err := makeRequest("QueryServerState", "{}", response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type QueryServerStateOutput struct {
	Data ServerStateDataOutput `json:"data"`
}

type ServerStateDataOutput struct {
	ServerGameState ServerGameStateOutput `json:"serverGameState"`
}

type ServerGameStateOutput struct {
	ActiveSessionName   string  `json:"activeSessionName"`
	NumConnectedPlayers int     `json:"numConnectedPlayers"`
	PlayerLimit         int     `json:"playerLimit"`
	TechTier            int     `json:"techTier"`
	ActiveSchematic     string  `json:"activeSchematic"`
	GamePhase           string  `json:"gamePhase"`
	IsGameRunning       bool    `json:"isGameRunning"`
	TotalGameDuration   int     `json:"totalGameDuration"`
	IsGamePaused        bool    `json:"isGamePaused"`
	AverageTickRate     float64 `json:"averageTickRate"`
	AutoLoadSessionName string  `json:"autoLoadSessionName"`
}

func makeRequest(function string, data string, response interface{}) error {
	token := os.Getenv("SATISFACTORY_TOKEN")
	ip := os.Getenv("SATISFACTORY_IP")
	port := os.Getenv("SATISFACTORY_PORT")

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s:%s/api/v1", ip, port), strings.NewReader(fmt.Sprintf(`{"function":"%s","data":%s}`, function, data)))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Skip SSL verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	// Process the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	if len(body) > 0 {
		return json.Unmarshal(body, response)
	}
	return nil
}
