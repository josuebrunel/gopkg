package pbc

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/davesavic/clink"
)

const (
	EndpointHealth    = "/api/health"
	EndpointAuthAdmin = "/api/admins/auth-with-password"
	EndpointAuthUser  = "/api/collections/users/auth-with-password"
	EndpointRecords   = "/api/collections/%s/records"
	EndpointRecordID  = "/api/collections/%s/records/%s"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

func getEndpoint(endpoint string, params ...any) string {
	return fmt.Sprintf(endpoint, params...)
}

type Client struct {
	BaseURL string
	Token   string
	Client  *clink.Client
}

func New(baseURL string) Client {
	c := clink.NewClient()
	c.Headers["Content-Type"] = "application/json"
	return Client{BaseURL: baseURL, Client: c}
}

func (c Client) buildUrl(path string, query *Query) string {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		slog.Error("failed to parse base url", "baseURL", c.BaseURL)
		return c.BaseURL
	}

	parsedPath, err := url.Parse(getEndpoint(path, query.Params...))
	if err != nil {
		slog.Error("failed to parse path", "path", path)
		return c.BaseURL
	}
	u := base.ResolveReference(parsedPath)
	q := u.Query()

	if len(query.Expand) > 0 {
		q.Set("expand", QmListString(query.Expand))
	}
	if len(query.Fields) > 0 {
		q.Set("fields", QmListString(query.Fields))
	}
	if len(query.Sort) > 0 {
		q.Set("sort", QmListString(query.Sort))
	}
	if query.Filters != "" {
		q.Set("filters", query.Filters)
	}
	if query.Page > 0 {
		q.Set("page", strconv.Itoa(query.Page))
		q.Set("perPage", strconv.Itoa(query.PerPage))
		if query.SkipTotal {
			q.Set("skipTotal", "true")
		}
	}

	u.RawQuery = q.Encode()
	slog.Debug("built url", "url", u.String())
	return u.String()
}

func (c *Client) Request(method, url string, opts ...QueryOption) (*http.Response, error) {
	query := NewQuery(opts...)
	url = c.buildUrl(url, query)

	slog.Debug("request attr", "url", url, "headers", query.Headers, "payload", query.Data)
	req, err := http.NewRequest(method, url, query.DataBytes)
	if err != nil {
		slog.Error("failed to prepare request", "url", url)
		return nil, err
	}

	for key, val := range query.Headers {
		req.Header.Set(key, val)
	}

	slog.Debug("calling", "url", req.URL.String())
	resp, err := c.Client.Do(req)
	if err != nil {
		slog.Error("error while trying to fetch url", "url", url, "status-code", resp.StatusCode)
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, ErrInvalidStatusCode
	}
	return resp, nil
}

func (c *Client) Health() (HealthResponse, error) {
	resp, err := c.Request(http.MethodGet, EndpointHealth)
	return ResponseTo[HealthResponse](resp), err
}

func (c *Client) Auth(endpoint, username, password string) (*http.Response, error) {
	payload := RequestAuth{
		Identity: username,
		Password: password,
	}
	return c.Request(http.MethodPost, endpoint, WithData(payload))
}

func (c *Client) AdminAuth(username, password string) (*http.Response, error) {
	return c.Auth(EndpointAuthAdmin, username, password)
}

func (c *Client) UserAuth(username, password string) (*http.Response, error) {
	return c.Auth(EndpointAuthUser, username, password)
}

func (c *Client) RecordCreate(name string, opts ...QueryOption) (*http.Response, error) {
	opts = append(opts, WithParams(name))
	return c.Request(http.MethodPost, EndpointRecords, opts...)
}

func (c *Client) RecordGet(name string, id string, opts ...QueryOption) (*http.Response, error) {
	opts = append(opts, WithParams(name, id))
	return c.Request(http.MethodGet, EndpointRecordID, opts...)
}

func (c *Client) RecordList(name string, opts ...QueryOption) (*http.Response, error) {
	opts = append(opts, WithParams(name))
	return c.Request(http.MethodGet, EndpointRecords, opts...)
}

func (c *Client) RecordUpdate(name string, id string, opts ...QueryOption) (*http.Response, error) {
	opts = append(opts, WithParams(name, id))
	return c.Request(http.MethodPatch, EndpointRecordID, opts...)
}

func (c *Client) RecordDelete(name string, id string, opts ...QueryOption) (*http.Response, error) {
	opts = append(opts, WithParams(name, id))
	return c.Request(http.MethodDelete, EndpointRecordID, opts...)
}

func ResponseTo[T any](resp *http.Response) T {
	var t T
	if err := clink.ResponseToJson(resp, &t); err != nil {
		slog.Error("failed to unmarshal response", "t", t)
	}
	return t
}
