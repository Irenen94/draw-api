package kms

func (client *Client) GetSecret(request *GetSecretValueRequest) (*GetSecretValueResponse, error) {
	var response = GetSecretValueResponse{}
	err := client.DoAction(request, &response)
	return &response, err
}

type GetSecretValueRequest struct {
	*baseRequest
	secretName string `json:"secretName"`
	version    string `json:"version,omitempty"`
}
type GetSecretValueResponse struct {
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

func CreateGetSecretValueRequest() (request *GetSecretValueRequest) {
	req := GetSecretValueRequest{
		baseRequest: defaultBaseRequest("/api/secret/getSecretValue"),
	}
	return &req
}
func (c *GetSecretValueRequest) SetSecretName(secretName string) {
	c.secretName = secretName
	c.addQueryParam("secretName", secretName)
}
func (c *GetSecretValueRequest) GetSecretName() string {
	return c.secretName
}
func (c *GetSecretValueRequest) SetVersion(version string) {
	c.version = version
	c.addQueryParam("version", version)
}
func (c *GetSecretValueRequest) GetVersion() string {
	return c.version
}
