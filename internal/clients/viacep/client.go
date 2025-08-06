package viacep

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

type Client interface {
	GetAddressByCEP(cep string) (*Response, error)
}

type client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://viacep.com.br/ws",
	}
}

func (c *client) GetAddressByCEP(cep string) (*Response, error) {
	url := fmt.Sprintf("%s/%s/json/", c.baseURL, cep)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP returned status code: %d", resp.StatusCode)
	}

	var viaCEPResp Response
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResp); err != nil {
		return nil, fmt.Errorf("error decoding ViaCEP response: %w", err)
	}

	if viaCEPResp.Erro {
		return nil, fmt.Errorf("CEP not found")
	}

	return &viaCEPResp, nil
}
