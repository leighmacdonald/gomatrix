package gomatrix

import (
	"github.com/leighmacdonald/golm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestFullSession(t *testing.T) {
	// Create device & id (private) keys
	act, err := golm.InitDevice()
	assert.NoError(t, err)
	act2, err := store.LoadAccount(deviceId)
	assert.True(t, assert.ObjectsAreEqual(act, act2))

	// Matrix: Should be sent by client to /keys/upload ['device_keys']
	// signed with Account.Sign()
	idKeys, err := act.GetIdentityKeys()
	assert.NoError(t, err)
	assert.Len(t, idKeys.Curve25519, 43)
	assert.Len(t, idKeys.Ed25519, 43)
	idKeys2, _ := act.GetIdentityKeys()
	assert.True(t, assert.ObjectsAreEqual(idKeys, idKeys2))

	log.Println(deviceId)
	act.GenerateOneTimeKeys(1)
	keys, err := act.GetOneTimeKeys()
	assert.NoError(t, err)
	act.MarkKeysAsPublished()
	log.Println(keys.Curve25519)
}
