package kms

func (client *Client) ExportKeySecret(request *ExportKeySecretRequest) (*ExportKeySecretResponse, error) {
	var response = ExportKeySecretResponse{}
	err := client.DoAction(request, &response)
	return &response, err
}

type ExportKeySecretRequest struct {
	*baseRequest
	KeyId        string `json:"keyId"`
	KeyVersionId string `json:"KeyVersionId,omitempty"` //主密钥版本
	PublicKey    string `json:"publicKey"`
}

type ExportKeySecretResponse struct {
	KeyId        string `json:"keyId"`
	KeySpec      string `json:"keySpec"`
	KeyUsage     string `json:"keyUsage"`
	State        string `json:"state"`
	KeyVersionId string `json:"KeyVersionId"` //主密钥版本
	SecretKey    string `json:"secretKey"`
	Secret       string `json:"secret"` //密钥信息，通过公钥加密
	RequestId    string `json:"requestId"`
}

func CreateExportKeySecretRequest() (request *ExportKeySecretRequest) {
	req := ExportKeySecretRequest{
		baseRequest: defaultBaseRequest("/api/key/exportKeySecret"),
	}
	return &req
}
func (c *ExportKeySecretRequest) SetKeyId(keyid string) {
	c.KeyId = keyid
	c.addQueryParam("keyId", keyid)
}
func (c *ExportKeySecretRequest) GetKeyId() string {
	return c.KeyId
}
func (c *ExportKeySecretRequest) SetKeyVersionId(keyVersionId string) {
	c.KeyVersionId = keyVersionId
	c.addQueryParam("keyVersionId", keyVersionId)
}
func (c *ExportKeySecretRequest) GetKeyVersionId() string {
	return c.KeyVersionId
}
func (c *ExportKeySecretRequest) SetPublicKey(publicKey string) {
	c.PublicKey = publicKey
	c.addQueryParam("publicKey", publicKey)
}
func (c *ExportKeySecretRequest) GetPublicKey() string {
	return c.PublicKey
}

type KeyInfo struct {
	KeyId        string `json:"keyId"`
	KeySpec      string `json:"keySpec"`
	KeyUsage     string `json:"keyUsage"`
	State        string `json:"state"`
	KeyVersionId string `json:"KeyVersionId"` //主密钥版本
	Secret       []byte `json:"secret"`       //密钥信息，通过公钥加密
}

func NewKeyInfo(keyId, keySpec, keyUsage, state, KeyVersionId string, secret []byte) KeyInfo {
	k := KeyInfo{
		KeyId:        keyId,
		KeySpec:      keySpec,
		KeyUsage:     keyUsage,
		State:        state,
		KeyVersionId: KeyVersionId,
		Secret:       make([]byte, 0),
	}
	k.Secret = append(k.Secret, secret...)
	return k
}
