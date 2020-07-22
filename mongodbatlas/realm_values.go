package mongodbatlas

import (
        "context"
        "fmt"
        "encoding/json"
        "net/http"
        "net/url"
        "log"
)
const (
    realmValuesPath = "groups/%s/apps/%s/values"
)


func (c *Client) RealmValueFromString(key string, value string) (*RealmValue, error) {
    t := make(map[string]interface{})
    err := json.Unmarshal([]byte(value), &t)
    if err != nil {
        log.Printf("Unmarshal: %v", err)
        return nil, err
    }

    v := &RealmValue{
        Name: key,
        Value: t,
    }
    return v, nil
}

// RealmService is an interface for interfacing with the Realm

// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys/
type RealmValuesService interface {
        List(context.Context, string, string, *ListOptions) ([]RealmValue, *Response, error)
        Get(context.Context, string, string, string) (*RealmValue, *Response, error)
        Create(context.Context, string, string, *RealmValue) (*RealmValue, *Response, error)
        Update(context.Context, string, string, string, *RealmValue) (*RealmValue, *Response, error)
        Delete(context.Context, string, string, string) (*Response, error)
}

// RealmValuesServiceOp handles communication with the RealmValue related methods
// of the MongoDB Atlas API
type RealmValuesServiceOp service

var _ RealmValuesService = &RealmValuesServiceOp{}


// RealmValuesInput represents MongoDB API key input request for Create.
type RealmValue struct {
    ID          string                      `json:"_id,omitempty"`
    Name        string                      `json:"name,omitempty"`
    Value       map[string]interface{}      `json:"value,omitempty"`
    Private     bool                        `json:"private,omitempty"`
}


// realmValuesResponse is the response from the RealmValuesService.List.
type realmValuesResponse struct {
        Apps []RealmValue  
}

func (s *RealmValuesServiceOp) AddRealmAuthToRequest(ctx context.Context,request *http.Request) (error) {

        path := fmt.Sprintf("%s%s",realmDefaultBaseURL,realmLoginPath)
        realmAuth := currentRealmAuth
        log.Printf("Got realmAuth=%+v",realmAuth)
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

        log.Printf("REALM AUTH root: %v", root)
        currentRealmAuth = root
        token := fmt.Sprintf("Bearer %s", currentRealmAuth.AccessToken)
        log.Printf("REALM AUTH token: %s", token)

        request.Header.Add("Authorization", token )
    return nil


}
// List all API-KEY in the organization associated to {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-all/
func (s *RealmValuesServiceOp) List(ctx context.Context, groupID string, appID string, listOptions *ListOptions) ([]RealmValue, *Response, error) {
    path := fmt.Sprintf(realmValuesPath, groupID, appID)


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
    
    //log.Printf("values List path=%+v",path)

    //root := new(realmValuesResponse)
    root := make([]RealmValue,0)
    resp, err := s.Client.Do(ctx, req, &root)
    if err != nil {
            return nil, resp, err
    }

    log.Printf("root %+v",root)
    log.Printf("resp ------>>> %+v",resp)

    return root, resp, nil
}

// Get gets the RealmValue specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-one/
func (s *RealmValuesServiceOp) Get(ctx context.Context, groupID string, appID string, valueID string) (*RealmValue, *Response, error) {
        if appID == "" {
                return nil, nil, NewArgError("appID", "must be set")
        }

        basePath := fmt.Sprintf(realmValuesPath, groupID, appID)
        escapedEntry := url.PathEscape(valueID)
        path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

        path = fmt.Sprintf("%s%s", realmDefaultBaseURL, path)

        req, err := s.Client.NewRequest(ctx, http.MethodGet,path, nil)
        if err != nil {
                return nil, nil, err
        }

        s.AddRealmAuthToRequest(ctx,req)
        root := new(RealmValue)
        resp, err := s.Client.Do(ctx, req, root)
        if err != nil {
                return nil, resp, err
        }

        return root, resp, err
}

// Create an API Key by the {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-create-one/
func (s *RealmValuesServiceOp) Create(ctx context.Context, groupID string, appID string, createRequest *RealmValue) (*RealmValue, *Response, error) {
        if createRequest == nil {
                return nil, nil, NewArgError("createRequest", "cannot be nil")
        }

        path := fmt.Sprintf(realmValuesPath, groupID, appID)

        path = fmt.Sprintf("%s%s", realmDefaultBaseURL, path)

        req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
        if err != nil {
                return nil, nil, err
        }

        s.AddRealmAuthToRequest(ctx,req)
        root := new(RealmValue)
        resp, err := s.Client.Do(ctx, req, root)
        if err != nil {
                return nil, resp, err
        }

        return root, resp, err
}

// Update a API Key in the organization associated to {ORG-ID}
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-update-one/
func (s *RealmValuesServiceOp) Update(ctx context.Context, groupID string, appID string, keyID string, updateRequest *RealmValue) (*RealmValue, *Response, error) {
        if updateRequest == nil {
                return nil, nil, NewArgError("updateRequest", "cannot be nil")
        }

        basePath := fmt.Sprintf(realmValuesPath, groupID, appID)
        escapedEntry := url.PathEscape(keyID)
        path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

        path = fmt.Sprintf("%s%s", realmDefaultBaseURL, path)

        req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
        if err != nil {
                return nil, nil, err
        }

        s.AddRealmAuthToRequest(ctx,req)
        root := new(RealmValue)
        resp, err := s.Client.Do(ctx, req, root)
        if err != nil {
                return nil, resp, err
        }

        return root, resp, err
}

// Delete the API Key specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/apiKey-delete-one-apiKey/
func (s *RealmValuesServiceOp) Delete(ctx context.Context, groupID, appID string, keyID string) (*Response, error) {
        if appID == "" {
                return nil, NewArgError("appID", "must be set")
        }

        basePath := fmt.Sprintf(realmValuesPath, groupID, appID)
        escapedEntry := url.PathEscape(keyID)
        path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

        path = fmt.Sprintf("%s%s", realmDefaultBaseURL, path)

        req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
        if err != nil {
                return nil, err
        }

    s.AddRealmAuthToRequest(ctx,req)
        resp, err := s.Client.Do(ctx, req, nil)

        return resp, err
}

