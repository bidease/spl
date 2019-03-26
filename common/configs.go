package common

import (
	"fmt"

	"github.com/bidease/spl/tools"
)

// SSHKey ..
type SSHKey struct {
	Fingerprint string
	Name        string
	PublicKey   string `json:"public_key"`
}

// GetSSHKeys ..
func GetSSHKeys() []SSHKey {
	type rawSSHKeys struct {
		Data     []SSHKey
		NumFound uint
	}
	var rawData rawSSHKeys
	tools.Request(fmt.Sprintf("ssh_keys"), &rawData, nil)

	var SSHKeys []SSHKey
	for _, item := range rawData.Data {
		SSHKeys = append(SSHKeys, SSHKey{
			Fingerprint: item.Fingerprint,
			Name:        item.Name,
			PublicKey:   item.PublicKey,
		})
	}
	return SSHKeys
}
