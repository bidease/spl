package spl

// SSHKey ..
type SSHKey struct {
	Name        string
	Fingerprint string
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

// GetSSHKeys ..
func GetSSHKeys() (*[]SSHKey, error) {
	var sshkeys []SSHKey

	_, err := RequestGet("ssh_keys", &sshkeys)
	if err != nil {
		return nil, err
	}

	return &sshkeys, nil
}
