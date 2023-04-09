package kms

// interface
type LongforRequest interface {
	addQueryParam(key, value string)
	getQueryParams() map[string]string

	addBodyParam(key, value string)
	getBodyParam() map[string]string

	getRoute() string
	setRoute(route string)
}

// template class
type baseRequest struct {
	route       string
	queryParams map[string]string
	bodyParams  map[string]string
}

func (request *baseRequest) addQueryParam(key, value string) {
	request.queryParams[key] = value
}
func (request *baseRequest) getQueryParams() map[string]string {
	return request.queryParams
}

func (request *baseRequest) addBodyParam(key, value string) {
	request.bodyParams[key] = value
}
func (request *baseRequest) getBodyParam() map[string]string {
	return request.bodyParams
}

func (request *baseRequest) setRoute(route string) {
	request.route = route
}
func (request *baseRequest) getRoute() string {
	return request.route
}

func defaultBaseRequest(route string) (request *baseRequest) {
	return &baseRequest{
		route:       route,
		queryParams: make(map[string]string),
		bodyParams:  make(map[string]string),
	}
}
