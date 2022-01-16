package pastefy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Success bool `json:"success"`
}

type PastefyApiClient struct {
	baseUrl    string
	apiToken   string
	httpClient http.Client
}

var PasteType = struct {
	MULTI_PASTE string
	PASTE       string
}{
	MULTI_PASTE: "MULTI_PASTE",
	PASTE:       "PASTE",
}

type User struct {
	AuthTypes      []string `json:"auth_types"`
	AuthType       string   `json:"auth_type"`
	Color          string   `json:"color"`
	LoggedIn       bool     `json:"logged_in"`
	Name           string   `json:"name"`
	ProfilePicture string   `json:"profile_picture"`
	Id             string   `json:"id"`
}

type Overview struct {
	Folder []Folder `json:"folder"`
	Pastes []Paste  `json:"pastes"`
}

type Error struct {
	Exception string `json:"exception"`
	Success   bool   `json:"success"`
	Exists    bool   `json:"exists"`
	Error     bool   `json:"error"`
}

func NewClient() PastefyApiClient {
	api := PastefyApiClient{
		baseUrl:    "https://pastefy.ga/api/v2",
		httpClient: http.Client{},
	}
	return api
}
func NewClientWithBaseURL(baseURL string) PastefyApiClient {
	api := NewClient()
	api.SetBaseURL(baseURL)
	return api
}

func (apiClient *PastefyApiClient) SetApiToken(token string) {
	apiClient.apiToken = token
}

func (apiClient *PastefyApiClient) SetBaseURL(baseURL string) {
	apiClient.baseUrl = baseURL
}

func (apiClient PastefyApiClient) Request(method string, url string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyJson, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyJson)
	}

	req, err := http.NewRequest(method, apiClient.baseUrl+url, bodyReader)

	if err != nil {
		return nil, err
	}
	if apiClient.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+apiClient.apiToken)
	}
	res, err := apiClient.httpClient.Do(req)

	if res.StatusCode > 305 {
		errorResponse := Error{}
		all, _ := ioutil.ReadAll(res.Body)
		err := json.Unmarshal(all, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("error while request. Exception: " + errorResponse.Exception + ".")
	}
	return res, err
}

func (apiClient PastefyApiClient) RequestMap(method string, url string, body interface{}, ma interface{}) (*http.Response, error) {
	response, err := apiClient.Request(method, url, body)
	if err != nil {
		return response, err
	}
	all, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	err2 := json.Unmarshal(all, &ma)
	if err2 != nil {
		return nil, err2
	}
	return response, err
}

func (paste Paste) GetMultiPasteParts() []MultiPastePart {
	if paste.Type == PasteType.MULTI_PASTE {
		list := []MultiPastePart{}
		err := json.Unmarshal([]byte(paste.Content), &list)
		if err != nil {
			return nil
		}
		return list
	}
	return nil
}

func (apiClient PastefyApiClient) GetUser() (User, error) {
	user := User{}
	_, err := apiClient.RequestMap("GET", "/user", nil, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (apiClient PastefyApiClient) GetOverview() (Overview, error) {
	overview := Overview{}
	_, err := apiClient.RequestMap("GET", "/user/overview", nil, &overview)
	if err != nil {
		return Overview{}, err
	}
	return overview, nil
}
