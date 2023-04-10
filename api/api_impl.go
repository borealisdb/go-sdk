package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Impl struct {
	httpClient   *http.Client
	baseUrl      string
	getTokenFunc TokenFunc
}

func New(baseUrl string, getTokenFunc TokenFunc, config Config) (*Impl, error) {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 200,
		IdleConnTimeout:     90 * time.Second,
	}

	if config.TlsCaLocation != "" {
		caCert, err := ioutil.ReadFile(config.TlsCaLocation)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tr.TLSClientConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	return &Impl{httpClient: client, baseUrl: baseUrl, getTokenFunc: getTokenFunc}, nil
}

func (i *Impl) GenerateClusterToken(request GenerateClusterTokenRequest) (GenerateClusterTokenResponse, error) {
	return requestWithBody[GenerateClusterTokenResponse](
		i.httpClient,
		request,
		fmt.Sprintf("%v/clusters/oauth2/token", i.baseUrl),
		http.MethodPost,
		i.getTokenFunc,
	)
}

func (i *Impl) IntrospectClusterToken() {

}

func (i *Impl) AttachRolesToUser(request AttachRolesToUserRequest) (AttachRolesToUserResponse, error) {
	return requestWithBody[AttachRolesToUserResponse](
		i.httpClient,
		request,
		fmt.Sprintf("%v/users/roles", i.baseUrl),
		http.MethodPost,
		i.getTokenFunc,
	)
}

func requestWithBody[Output any](client *http.Client, body any, path, method string, tokenFunc TokenFunc) (Output, error) {
	var emptyOut Output
	var reader io.Reader
	if body != nil {
		requestBodyBytes, err := json.Marshal(body)
		if err != nil {
			return emptyOut, err
		}
		reader = bytes.NewReader(requestBodyBytes)
	}
	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		return emptyOut, err
	}

	token, err := tokenFunc()
	if err != nil {
		return emptyOut, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return emptyOut, err
	}

	if resp.StatusCode == http.StatusOK {
		responseBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return emptyOut, err
		}

		defer resp.Body.Close()
		var out Output
		if err := json.Unmarshal(responseBodyBytes, &out); err != nil {
			return emptyOut, err
		}
		return out, nil
	}

	if resp.StatusCode == http.StatusCreated {
		return emptyOut, nil
	}

	return emptyOut, fmt.Errorf("unexpected status from api: %v", resp.StatusCode)
}
