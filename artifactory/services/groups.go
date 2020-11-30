package services

import (
	"encoding/json"
	"fmt"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// application/vnd.org.jfrog.artifactory.security.Group+json
type Group struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	AutoJoin        bool   `json:"autoJoin"`
	AdminPrivileges bool   `json:"adminPrivileges"`
	Realm           string `json:"realm"`
	RealmAttributes string `json:"realmAttributes"`
}
type GroupService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewGroupService(client *rthttpclient.ArtifactoryHttpClient) *GroupService {
	return &GroupService{client: client}
}
func (gs *GroupService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	gs.ArtDetails = rt
}

func (gs *GroupService) GetGroup(name string) (*Group, error) {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/groups/%s", gs.ArtDetails.GetUrl(), name)
	_, body, _, err := gs.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	var group Group
	if err := json.Unmarshal(body, &group); err != nil {
		return nil, errorutils.CheckError(err)
	}
	return &group, nil
}

func (gs *GroupService) CreateGroup(group Group) error {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()

	if content, err := json.Marshal(group); err != nil {
		url := fmt.Sprintf("%sapi/groups/%s", gs.ArtDetails.GetUrl(), group.Name)
		if _, _, err := gs.client.SendPut(url, content, &httpDetails); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GroupService) DeleteGroup(name string) error {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/groups/%s", gs.ArtDetails.GetUrl(), name)
	_, _, err := gs.client.SendDelete(url, nil, &httpDetails)
	return err
}

func (gs *GroupService) GroupExits(name string) (bool, error) {
	log.Debug("Hello, I am logging")
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/groups/%s", gs.ArtDetails.GetUrl(), name)
	res, _, err := gs.client.SendHead(url, &httpDetails)
	exists := res.StatusCode == 200
	if err != nil || !exists {
		return false, err
	}

	return true, nil
}
