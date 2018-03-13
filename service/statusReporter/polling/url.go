package polling

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/iryonetwork/wwm/status"
)

// type for URL
type URLType string

// URLType values
const (
	TypeInternalURL URLType = "Internal"
	TypeExternalURL URLType = "External"
)

// ExternalURL is type for external service availability endpoint URL
type ExternalURL struct {
	url     string
	timeout time.Duration
}

func (u ExternalURL) FetchStatus() (*status.Response, error) {
	client := http.Client{Timeout: u.timeout}
	resp, err := client.Get(u.url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	statusResp := status.Response{Status: status.OK}

	if resp.StatusCode != http.StatusOK {
		statusResp = status.Response{Status: status.Warning, Msg: fmt.Sprintf("unexpected response code: %d", resp.StatusCode)}
	}

	return &statusResp, nil
}

func (u ExternalURL) String() string {
	return string(u.url)
}

// InternalURL is type for internal service status endpoint URL
type InternalURL struct {
	ExternalURL
}

func (u InternalURL) FetchStatus() (*status.Response, error) {
	client := http.Client{Timeout: u.timeout}
	resp, err := client.Get(u.url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	statusResp := status.Response{}
	if resp.StatusCode != http.StatusOK {
		statusResp = status.Response{Status: status.Error, Msg: fmt.Sprintf("unexpected response code: %d", resp.StatusCode)}
	}

	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		return nil, err
	}

	if statusResp.Status.Int() == -1 {
		// invalid status code, turn into error
		statusResp.Msg = fmt.Sprintf("invalid status value \"%s\"", statusResp.Status)
		statusResp.Status = status.Error
	}

	return &statusResp, nil
}

func NewExternalURL(url string, timeout time.Duration) URLStatusEndpoint {
	return &ExternalURL{url: url, timeout: timeout}
}

func NewInternalURL(url string, timeout time.Duration) URLStatusEndpoint {
	return &InternalURL{ExternalURL{url: url, timeout: timeout}}
}
