package dropbox

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

const lockEndTimeLayout = "2006-01-02 15:04:05 -0700"
const lockFileName = "/kfs-root-dir/lock.yaml"

func (client *Client) withLock(clientID string, fn func()) error {
	err := client.lock(clientID)
	if err != nil {
		return err
	}
	err = client.checkLock(clientID)
	if err != nil {
		return err
	}
	defer client.unlock()
	fn()
	return nil
}

func (client *Client) lock(clientID string) error {
start:
	_, content, err := client.c.Download(files.NewDownloadArg(lockFileName))
	if err == nil {
		lock := &Lock{}
		err = yaml.NewDecoder(content).Decode(lock)
		if err != nil {
			panic(err)
		}
		endTime, err := time.Parse(lockEndTimeLayout, lock.EndTime)
		if err != nil {
			panic(err)
		}
		if endTime.After(time.Now()) {
			fmt.Println(clientID, "wait for lock in client", lock.ClientID)
			time.Sleep(time.Second)
			goto start
		}
	} else if !strings.Contains(err.Error(), pathNotFoundError) {
		fmt.Println("GetMetadata.lock", err)
		return err
	}

	lock, err := yaml.Marshal(&Lock{
		ClientID: clientID,
		EndTime:  time.Now().Add(time.Second * 10).Format(lockEndTimeLayout),
	})
	if err != nil {
		fmt.Println("yaml.Marshal", err)
		return err
	}
	_, err = client.c.Upload(files.NewCommitInfo(lockFileName), bytes.NewBuffer(lock))
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) checkLock(clientID string) error {
	_, content, err := client.c.Download(files.NewDownloadArg(lockFileName))
	if err != nil {
		return err
	}
	lock := &Lock{}
	err = yaml.NewDecoder(content).Decode(lock)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(lockEndTimeLayout, lock.EndTime)
	if err != nil {
		panic(err)
	}
	if clientID != lock.ClientID {
		return fmt.Errorf("[%s]locked by another client(%s)", clientID, lock.ClientID)
	}
	if endTime.Before(time.Now()) {
		return fmt.Errorf("lock time(%s) is before now(%s)", lock.EndTime, time.Now().Format(lockEndTimeLayout))
	}
	return nil
}

func (client *Client) unlock() error {
	_, err := client.c.DeleteV2(files.NewDeleteArg(lockFileName))
	return err
}
