package kms

import (
	"strconv"
)

func (client *Client) Encrypt(request *EncryptRequest) (*EncryptResponse, error) {
	var response = EncryptResponse{}
	err := client.DoAction(request, &response)
	return &response, err
}

type EncryptRequest struct {
	*baseRequest
	keyId     string `json:"keyId"`
	plaintext string `json:"plaintext"`
	random    bool   `json:"random"`
}
type EncryptResponse struct {
	CiphertextBlob string `json:"ciphertextBlob"`
	KeyId          string `json:"keyId"`
	KeyVersionId   string `json:"keyVersionId"`
	RequestId      string `json:"requestId"`
}

// CreateEncryptRequest creates a request to invoke Encrypt API
func CreateEncryptRequest() (request *EncryptRequest) {
	req := EncryptRequest{
		baseRequest: defaultBaseRequest("/api/key/encrypt"),
	}
	req.SetRandom(true)
	return &req
}
func (c *EncryptRequest) SetKeyId(keyid string) {
	c.keyId = keyid
	c.addQueryParam("keyId", keyid)
}
func (c *EncryptRequest) GetKeyId() string {
	return c.keyId
}
func (c *EncryptRequest) SetPlaintext(plaintext string) {
	c.plaintext = plaintext
	c.addQueryParam("plaintext", plaintext)
}
func (c *EncryptRequest) GetPlaintext() string {
	return c.plaintext
}

func (c *EncryptRequest) SetRandom(random bool) {
	c.random = random
	c.addQueryParam("random", strconv.FormatBool(random))
}
func (c *EncryptRequest) GetRandom() bool {
	return c.random
}
