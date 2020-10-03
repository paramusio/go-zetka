package zetka

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
	"io"
	"net/http"
	"net/url"
)

type Encoding string

func (e Encoding) String() string {
	return string(e)
}

var (
	JSON Encoding = "json"
	// ETF Not yet supported
	//ETF Encoding = "etf"
)

func gatewayURI(baseURL, token, version string, encoding Encoding) (string, error) {
	// Validate our encoding
	switch encoding {
	case JSON:
		break
	default:
		return "", errors.New("invalid encoding")
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+GetGatewayPath, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed constructing request")
	}

	headers := req.Header.Clone()
	headers.Set("Authorization", fmt.Sprintf("Bot %s", token))
	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed performing request")
	}

	defer resp.Body.Close()
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", err
	}

	uri := fastjson.GetString(buf.Bytes(), "url")
	if len(uri) == 0 {
		return "", fmt.Errorf("something fucked up")
	}

	u, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrap(err, "failed parsing gateway uri from discord")
	}

	vals := url.Values{}
	vals.Set("v", version)
	vals.Set("encoding", encoding.String())
	u.RawQuery = vals.Encode()

	return u.String(), nil
}
