package accounts

const (
	basePath = "/v1/organisation/accounts"
)

type httpUtils interface {
	Delete(resourcePath string) error
}

type Client struct {
	http httpUtils
}

func NewClient(httpUtils httpUtils) Client {
	return Client{http: httpUtils}
}
