package kms

func (client *Client) NewSecret(request *NewSecretRequest) (*NewSecretResponse, error) {
	var response = NewSecretResponse{}
	err := client.DoAction(request, &response)
	return &response, err
}

type NewSecretRequest struct {
	*baseRequest
	keyId       string
	secretName  string
	secretData  string
	description string
}

type NewSecretResponse struct {
	KeyId       string
	SecretId    string
	AppId       string
	SecretName  string
	Version     string
	Description string
	State       bool
	CreatedBy   string
	Updated     string
	RequestId   string
}

func CreateNewSecretRequest() *NewSecretRequest {
	req := NewSecretRequest{
		baseRequest: defaultBaseRequest("/api/secret/create"),
	}
	return &req
}
func (c *NewSecretRequest) SetKeyId(keyId string) {
	c.keyId = keyId
	c.addQueryParam("keyId", keyId)
}
func (c *NewSecretRequest) GetKeyId() string {
	return c.keyId
}

func (c *NewSecretRequest) SetSecretName(secretName string) {
	c.secretName = secretName
	c.addQueryParam("secretName", secretName)
}
func (c *NewSecretRequest) GetSecretName() string {
	return c.secretName
}

func (c *NewSecretRequest) SetSecretData(secretData string) {
	c.secretData = secretData
	c.addBodyParam("secretData", secretData)
}
func (c *NewSecretRequest) GetSecretData() string {
	return c.secretData
}

func (c *NewSecretRequest) SetDescription(description string) {
	c.description = description
	c.addQueryParam("description", description)
}
func (c *NewSecretRequest) GetDescription() string {
	return c.description
}
