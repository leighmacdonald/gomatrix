package gomatrix

import (
	"encoding/json"
	"fmt"
	"github.com/leighmacdonald/golm"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func signObject(userId string, deviceId string, act golm.Account, o map[string]interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return o, err
	}
	sig := act.Sign(string(b))
	o["signatures"] = map[string]interface{}{
		userId: map[string]string{
			fmt.Sprintf("ed25519:%s", deviceId): sig,
		},
	}
	return o, nil
}

// InitDevice creates a new device and stores the keys and account into
// the configured Store.
func InitDeviceCrypto(cli *Client, deviceId string) (act golm.Account, err error) {
	act, err = cli.Store.LoadAccount(deviceId)
	if err != nil {
		act = golm.CreateNewAccount()
		if err = cli.Store.SaveAccount(deviceId, act); err != nil {
			return act, err
		}
		return act, err
	}

	return act, nil
}

// NewEncryptedClient creates a new Matrix Client ready for syncing which understands olm & megolm protocols
func NewEncryptedClient(homeServerURL, userID, accessToken string, deviceId string) (*Client, error) {
	hsURL, err := url.Parse(homeServerURL)
	if err != nil {
		return nil, err
	}
	// By default, use an in-memory store which will never save filter ids / next batch tokens to disk.
	// The client will work with this storer: it just won't remember across restarts.
	// In practice, a database backend should be used.
	store := NewInMemoryStore()
	cli := Client{
		AccessToken:   accessToken,
		HomeserverURL: hsURL,
		UserID:        userID,
		Prefix:        "/_matrix/client/r0",
		Syncer:        NewDefaultSyncer(userID, store),
		Store:         store,
	}
	// By default, use the default HTTP client.
	cli.Client = http.DefaultClient

	if deviceId == "" {
		deviceId = RandStringRunes(20)
		fmt.Print("Created new device ID")
	}
	_, err = InitDeviceCrypto(&cli, deviceId)
	if err != nil {
		log.Println("Failed to initialize client crypto")
	}

	return &cli, nil
}
