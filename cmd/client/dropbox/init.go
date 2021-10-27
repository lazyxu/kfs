package dropbox

import (
	"encoding/json"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"strings"
)

type Client struct {
	c files.Client
}

var LogLevel = dropbox.LogInfo

func New() *Client {
	config := dropbox.Config{
		Token:    "sl.A7AAbH7-z2yF87aByM9B-vNUkmiw4XrTHpkdzah-6S1WSF9fvCrlu9UAw80BSfoyQ8g7T5L2YONQzKnsYnC_46XUe6AMzvvT6xrGi7m8NrwtmTCIL1x189VX6ty3q4IHVa9879_vT7ex",
		LogLevel: LogLevel, // if needed, set the desired logging level. Default is off
	}
	c := files.New(config)
	return &Client{
		c: c,
	}
}

func (client *Client) ListFolder() error {
	resp, err := client.c.ListFolder(files.NewListFolderArg(""))
	if err != nil {
		fmt.Println("ListFolder", err)
		return err
	}
	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Unmarshal", err)
		return err
	}
	fmt.Println("resp", string(res))
	return nil
}

const pathConflictFolderError = "path/conflict/folder/"
const pathNotFoundError = "path/not_found/"

type Lock struct {
	ClientID string `yaml:"client_id"`
	EndTime  string `yaml:"end_time"`
}

func (client *Client) Init(clientID string) error {
	// ignore error: path/conflict/folder/..
	_, err := client.c.CreateFolderV2(files.NewCreateFolderArg("/kfs-root-dir"))
	if err != nil && !strings.Contains(err.Error(), pathConflictFolderError) {
		fmt.Println("/kfs-root-dir", err)
		return err
	}
	// TODO: CreateFolderBatch
	_, err = client.c.CreateFolderV2(files.NewCreateFolderArg("/kfs-root-dir/objects"))
	if err != nil && !strings.Contains(err.Error(), pathConflictFolderError) {
		fmt.Println("/kfs-root-dir/objects", err)
		return err
	}
	_, err = client.c.CreateFolderV2(files.NewCreateFolderArg("/kfs-root-dir/branches"))
	if err != nil && !strings.Contains(err.Error(), pathConflictFolderError) {
		fmt.Println("/kfs-root-dir/branches", err)
		return err
	}
	err = client.withLock(clientID, func() {
		_, err := client.c.Upload(files.NewCommitInfo("/kfs-root-dir/branches/"+clientID+".yaml"), strings.NewReader("empty"))
		if err != nil {
			fmt.Println("Upload.branch", err)
			return
		}
	})
	if err != nil {
		fmt.Println("withLock", err)
		return err
	}
	return nil
}
