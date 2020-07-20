package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
    "log"
)
const (
    realmAppsPath = "/groups/%s/apps"
	realmDefaultBaseURL = "https://realm.mongodb.com/api/admin/v3.0"
    realmLoginPath = "/auth/providers/mongodb-cloud/login"
)

type RealmAtlasApiKey struct {
    Username string
    Password string
}
type RealmAuth struct {
	AccessToken string   `json:"access_token,omitempty"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	UserID         string      `json:"user_id,omitempty"`
	DeviceID       string      `json:"device_id,omitempty"`
}

// RealmService is an interface for interfacing with the Realm
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys/
type RealmAppsService interface {
	List(context.Context, string, *ListOptions) ([]RealmApp, *Response, error)
	Get(context.Context, string, string) (*RealmApp, *Response, error)
	Create(context.Context, string, *RealmAppsInput) (*RealmApp, *Response, error)
	Update(context.Context, string, string, *RealmAppsInput) (*RealmApp, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// RealmAppsServiceOp handles communication with the RealmApp related methods
// of the MongoDB Atlas API
type RealmAppsServiceOp service

var _ RealmAppsService = &RealmAppsServiceOp{}

var RealmAccessToken = ""

// RealmAppsInput represents MongoDB API key input request for Create.
type RealmAppsInput struct {
	Name string   `json:"name,omitempty"`
}

// RealmApp represents MongoDB API Key.
//{"_id":"5f12de8c15049be9464eb269","client_app_id":"mad-elion-arays","name":"mad-elion","location":"US-VA","deployment_model":"GLOBAL","domain_id":"5f12de8c15049be9464eb26a","group_id":"5f12d8cc6c2bfd1e0c670f4a","last_used":1595072140,"last_modified":1595072140,"product":"standard"}
type RealmApp struct {
	Name string   `json:"name,omitempty"`
	ID         string      `json:"_id,omitempty"`
	ClientAppID       string      `json:"client_app_id,omitempty"`
	Location string      `json:"location,omitempty"`
	DeploymentModel string      `json:"deployment_model,omitempty"`
	GroupID string      `json:"group_id,omitempty"`
	Product string      `json:"product,omitempty"`
	DomainID string      `json:"domain_id,omitempty"`
}

// realmAppsResponse is the response from the RealmAppsService.List.
type realmAppsResponse struct {
	Apps []RealmApp  
}

var currentRealmAuth *RealmAuth
var currentRealmAtlasApiKey *RealmAtlasApiKey

func (c *Client) SetCurrentRealmAtlasApiKey(rk *RealmAtlasApiKey) {
    currentRealmAtlasApiKey = rk
}
func (c *Client) GetCurrentRealmAtlasApiKey() (*RealmAtlasApiKey) {
    return currentRealmAtlasApiKey
}

func (s *RealmAppsServiceOp) AddRealmAuthToRequest(ctx context.Context,request *http.Request) (error) {

	path := fmt.Sprintf("%s%s",realmDefaultBaseURL,realmLoginPath)
    data := map[string]interface{}{
		"username": currentRealmAtlasApiKey.Username,
		"apiKey":   currentRealmAtlasApiKey.Password,
	}

	loginReq, err := s.Client.NewRequest(ctx, http.MethodPost, path, &data)
	if err != nil {
		return err
	}

    root := &RealmAuth{}
	_, err = s.Client.Do(ctx, loginReq, root)
	if err != nil {
	    log.Printf("REALM AUTH error: %s", err)
		return err
	}

	//log.Printf("REALM AUTH root: %v", root)
    currentRealmAuth = root
    token := fmt.Sprintf("Bearer %s", currentRealmAuth.AccessToken)
	//log.Printf("REALM AUTH token: %s", token)

	request.Header.Add("Authorization", token )
    return nil


}
// List all API-KEY in the organization associated to {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-all/
func (s *RealmAppsServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]RealmApp, *Response, error) {
	path := fmt.Sprintf("/groups/%s/apps", groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	path = fmt.Sprintf("%s%s",realmDefaultBaseURL,path)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
    
    s.AddRealmAuthToRequest(ctx,req)
    //log.Printf("REALM - check token in header %v", req.Header)
    
    //root := new(realmAppsResponse)
    root := make([]RealmApp,0)	
    resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}


	return root, resp, nil
}

// Get gets the RealmApp specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-one/
func (s *RealmAppsServiceOp) Get(ctx context.Context, groupID string, appID string) (*RealmApp, *Response, error) {
	if appID == "" {
		return nil, nil, NewArgError("appID", "must be set")
	}

	basePath := fmt.Sprintf(realmAppsPath, groupID)
	escapedEntry := url.PathEscape(appID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	path = fmt.Sprintf("%s/%s", realmDefaultBaseURL, path)

	req, err := s.Client.NewRequest(ctx, http.MethodGet,path, nil)
	if err != nil {
		return nil, nil, err
	}

    s.AddRealmAuthToRequest(ctx,req)
	root := new(RealmApp)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create an API Key by the {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-create-one/
func (s *RealmAppsServiceOp) Create(ctx context.Context, groupID string, createRequest *RealmAppsInput) (*RealmApp, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

    path := fmt.Sprintf(realmAppsPath, groupID)

	path = fmt.Sprintf("%s/%s", realmDefaultBaseURL, path)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

    s.AddRealmAuthToRequest(ctx,req)
	root := new(RealmApp)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update a API Key in the organization associated to {ORG-ID}
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-update-one/
func (s *RealmAppsServiceOp) Update(ctx context.Context, groupID, appID string, updateRequest *RealmAppsInput) (*RealmApp, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(realmAppsPath, groupID)
	escapedEntry := url.PathEscape(appID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	path = fmt.Sprintf("%s/%s", realmDefaultBaseURL, path)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

    s.AddRealmAuthToRequest(ctx,req)
	root := new(RealmApp)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the API Key specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKey-delete-one-apiKey/
func (s *RealmAppsServiceOp) Delete(ctx context.Context, groupID, appID string) (*Response, error) {
	if appID == "" {
		return nil, NewArgError("appID", "must be set")
	}

	basePath := fmt.Sprintf(realmAppsPath, groupID)
	escapedEntry := url.PathEscape(appID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	path = fmt.Sprintf("%s/%s", realmDefaultBaseURL, path)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

    s.AddRealmAuthToRequest(ctx,req)
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
