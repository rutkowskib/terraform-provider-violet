package violet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type VioletClient struct {
	Username  string
	Password  string
	AppId     string
	AppSecret string
	Token     string
	BaseUrl   string
}

type VioletWebhook struct {
	Id               int64
	AppId            int64
	Event            string
	RemoteEndpoint   string
	Status           string
	DateCreated      string
	DateLastModified string
}

type violetWebhookResponse struct {
	Id               int64  `json:"id"`
	AppId            int64  `json:"app_id"`
	Event            string `json:"event"`
	RemoteEndpoint   string `json:"remote_endpoint"`
	Status           string `json:"status"`
	DateCreated      string `json:"date_created"`
	DateLastModified string `json:"date_last_modified"`
}

func (c *VioletClient) Login(ctx context.Context) {
	var body = []byte(fmt.Sprintf(`{
	"username": "%s",
	"password": "%s"
	}`, c.Username, c.Password))

	res := c.makeRequest(ctx, "POST", "login", body)

	type loginResponse struct {
		Token string `json:"token"`
	}

	var data loginResponse

	err := json.Unmarshal(res, &data)
	if err != nil {
		// Handle error
		return
	}

	c.Token = data.Token
}

func (c *VioletClient) GetWebhook(ctx context.Context, id int64) VioletWebhook {
	path := fmt.Sprintf("events/webhooks/%d", id)
	res := c.makeRequest(ctx, "GET", path, nil)

	var data violetWebhookResponse

	// TODO handle error
	json.Unmarshal(res, &data)

	tflog.Info(ctx, "data", map[string]any{
		"data": data,
		"res":  string(res),
	})

	return VioletWebhook(data)
}

type CreateWebhookInput struct {
	Event          string
	RemoteEndpoint string
}

func (c *VioletClient) CreateWebhook(ctx context.Context, input CreateWebhookInput) VioletWebhook {
	path := fmt.Sprintf("apps/%s/webhooks", c.AppId)
	body := []byte(fmt.Sprintf(`{
		"event": "%s",
		"remote_endpoint": "%s"
	}`, input.Event, input.RemoteEndpoint))

	tflog.Info(ctx, "Making create webhook request", map[string]any{
		"event":           input.Event,
		"remote_endpoint": input.RemoteEndpoint,
	})

	res := c.makeRequest(ctx, "POST", path, body)

	var data violetWebhookResponse

	// TODO handle error
	json.Unmarshal(res, &data)

	tflog.Info(ctx, "data", map[string]any{
		"data": data,
		"res":  string(res),
	})

	return VioletWebhook(data)
}

func (c *VioletClient) DeleteWebhook(ctx context.Context, id int64) {
	tflog.Info(ctx, "Deleting webhook", map[string]any{
		"id": id,
	})

	path := fmt.Sprintf("apps/%s/webhooks/%d", c.AppId, id)
	c.makeRequest(ctx, "DELETE", path, nil)

	tflog.Info(ctx, fmt.Sprint("Webhook %d deleted", id))
}

func (c *VioletClient) makeRequest(ctx context.Context, method string, path string, requestBody []byte) []byte {
	request, err := http.NewRequest(method, c.BaseUrl+path, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Violet-App-Id", c.AppId)
	request.Header.Set("X-Violet-App-Secret", c.AppSecret)

	if c.Token != "" {
		request.Header.Set("X-Violet-Token", c.Token)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	tflog.Info(ctx, "response", map[string]any{
		"status":  response.Status,
		"headers": response.Header,
		"body":    string(body),
	})

	return body
}
