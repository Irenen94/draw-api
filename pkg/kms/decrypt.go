package kms

func (client *Client) Decrypt(request *DecryptRequest) (*DecryptResponse, error) {
	var response = DecryptResponse{}
	err := client.DoAction(request, &response)
	return &response, err
}

type DecryptRequest struct {
	*baseRequest
	ciphertextBlob string `json:"ciphertextBlob"`
}

type DecryptResponse struct {
	Plaintext    string `json:"plaintext"`
	KeyId        string `json:"keyId"`
	KeyVersionId string `json:"keyVersionId"`
	RequestId    string `json:"requestId"`
}

func (c *DecryptRequest) SetCiphertextBlob(ciphertextBlob string) {
	c.ciphertextBlob = ciphertextBlob
	c.addQueryParam("ciphertextBlob", ciphertextBlob)
}
func (c *DecryptRequest) getCiphertextBlob() string {
	return c.ciphertextBlob
}

// CreateDecryptRequest creates a request to invoke Decrypt API
func CreateDecryptRequest() (request *DecryptRequest) {
	req := DecryptRequest{
		baseRequest: defaultBaseRequest("/api/key/decrypt"),
	}
	return &req
}
