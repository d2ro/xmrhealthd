package xmrhealthd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// Response data format
type State struct {
	CryptoCode string
	Synced     bool
}

// Subset of monerod RPC info data
type MoneroResult struct {
	Result struct {
		Synchronized bool `json:"synchronized"`
	} `json:"result"`
}

// QueryMonerod queries the monero daemon at the given IP address and returns the sync state.
func QueryMonerod(ipAddress string) (bool, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:18081/json_rpc", ipAddress), strings.NewReader(`{"jsonrpc":"2.0","id":"0","method":"get_info"}`))
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := (&http.Client{
		Timeout: 5 * time.Second,
	}).Do(req)
	if err != nil {
		return false, fmt.Errorf("doing request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("got status: %s", resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading body: %w", err)
	}
	res := &MoneroResult{}
	if err = json.Unmarshal(body, &res); err != nil {
		return false, fmt.Errorf("unmarshaling json: %w", err)
	}
	return res.Result.Synchronized, nil
}

// Run runs the xmrhealthd, querying the sync status from the given IP address every 10 seconds.
// Run blocks.
func Run(ipAddress string) {

	if net.ParseIP(ipAddress) == nil {
		log.Println("invalid IP address: %s", ipAddress)
		return
	}

	// state cache, update every 10 seconds

	var responseData []byte // marshaled JSON data
	var responseStatusCode = http.StatusOK

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for ; true; <-ticker.C {
			state, err := QueryMonerod(ipAddress)
			if err == nil {
				responseData, _ = json.Marshal(State{"XMR", state})
				responseStatusCode = http.StatusOK
			} else {
				responseData = nil
				responseStatusCode = http.StatusInternalServerError
				log.Println(err)
			}
		}
	}()

	// serve state

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(responseStatusCode) // call after adding headers
		w.Write(responseData)
	})
	log.Fatal(http.ListenAndServe(":64325", nil))
}
