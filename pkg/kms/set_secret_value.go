package kms

func (client *Client) SetSecretValue(request *SetSecretValueRequest) (*SetSecretValueResponse, error) {
	var response = SetSecretValueResponse{}
	err := client.DoAction(request, &response)
	return &response, err
}

type SetSecretValueRequest struct {
	*baseRequest
	secretName string
	secretData string
}

type SetSecretValueResponse struct {
	SecretId    string `json:"secretId"`
	KeyId       string `json:"keyId"`
	AppId       string `json:"appId"`
	SecretName  string `json:"secretName"`
	SecretData  string `json:"secretData"`
	Version     string `json:"version"`
	Description string `json:"description"`
	State       bool   `json:"state"`
	CreatedBy   string `json:"createdBy"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	RequestId   string `json:"requestId"`
}

func CreateSetSecretValueRequest() *SetSecretValueRequest {
	req := SetSecretValueRequest{
		baseRequest: defaultBaseRequest("/api/secret/addSecretValue"),
	}
	return &req
}
func (c *SetSecretValueRequest) SetSecretName(secretName string) {
	c.secretName = secretName
	c.addQueryParam("secretName", secretName)
}
func (c *SetSecretValueRequest) GetSecretName() string {
	return c.secretName
}
func (c *SetSecretValueRequest) SetSecretData(secretData string) {
	c.secretData = secretData
	c.addBodyParam("secretData", secretData)
}
func (c *SetSecretValueRequest) GetSecretData() string {
	return c.secretData
}
