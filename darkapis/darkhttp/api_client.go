package darkhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

type HttpClient struct {
	httpc   http.Client
	api_url string
}

func NewClient(
	api_url string, // for example https://darkstat.dd84ai.com
) *HttpClient {
	c := &HttpClient{
		httpc:   http.Client{},
		api_url: api_url,
	}

	return c
}

const ApplicationJson = "application/json"

type EmptyInput struct{}

func make_request[IN any, OUT any](c *HttpClient, endpoint_url core_types.Url, input IN) (OUT, error) {
	var output OUT

	post_body, err := json.Marshal(input)

	if logus.Log.CheckError(err, "failed to marshal input") {
		return output, err
	}

	res, err := c.httpc.Post(c.api_url+endpoint_url.ToStr(), ApplicationJson, bytes.NewBuffer(post_body))

	if logus.Log.CheckError(err, "failed to request to get pobs") {
		return output, err
	}

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return output, errors.New(fmt.Sprintln("not 200 status code, received status code=", res.StatusCode, " boby=", string(resBody)))
	}

	resBody, err := io.ReadAll(res.Body)
	if logus.Log.CheckError(err, "client: could not read response body\n") {
		return output, err
	}

	err = json.Unmarshal(resBody, &output)
	if logus.Log.CheckError(err, "failed to unmarshal output") {
		return output, err
	}

	return output, nil
}
