package config

type KmsMode int

const (
	LOCAL  KmsMode = 0
	REMOTE KmsMode = 1
)

func (c KmsMode) String() string {
	switch c {
	case LOCAL:
		return "LOCAL"
	case REMOTE:
		return "REMOTE"
	default:
		return "UNKNOWN"
	}
}

type LongforConfig struct {
	AccessKeyId       string
	AccessKeySecret   string
	BasePath          string
	GaiaApiKey        string
	timeout           int
	localCacheTimeout int
	kmsMode           KmsMode
}

func DefaultLongforConfig(accessKeyId, accessKeySecret, basePath string) *LongforConfig {
	return &LongforConfig{
		AccessKeyId:       accessKeyId,
		AccessKeySecret:   accessKeySecret,
		BasePath:          basePath,
		timeout:           3000,
		kmsMode:           REMOTE,
		localCacheTimeout: 60 * 60 * 4, //default 4 hour
	}
}
func (c *LongforConfig) GetKmsMode() KmsMode {
	return c.kmsMode
}
func (c *LongforConfig) SetKmsMode(mode KmsMode) {
	c.kmsMode = mode
}
func (c *LongforConfig) SetGaiaApiKey(gaiaApiKey string) {
	c.GaiaApiKey = gaiaApiKey
}
func (c *LongforConfig) SetLocalCacheTimeout(d int) {
	c.localCacheTimeout = d
}
func (c *LongforConfig) GetLocalCacheTimeout() int {
	return c.localCacheTimeout
}
