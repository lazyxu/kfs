package onedrive

import (
	"context"
	"fmt"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
)

func Login() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "..."},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := onedrive.NewClient(tc)

	// list all OneDrive drives for the current logged in user
	drives, err := client.Drives.List(ctx)
	if err != nil {
		return err
	}
	fmt.Println(drives)
	return nil
}
