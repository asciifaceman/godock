package framework

// Method wraps all HTTP methods valid for a route
type Method string

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
const (
	GET     Method = "GET"
	HEAD    Method = "HEAD"
	POST    Method = "POST" // create new object receive ID
	PUT     Method = "PUT"  // update object by ID
	DELETE  Method = "DELETE"
	CONNECT Method = "CONNECT"
	OPTIONS Method = "OPTIONS"
	TRACE   Method = "TRACE"
	PATCH   Method = "PATCH"
	ANY     Method = "ANY"
)

// String converts HTTP method to string
func (h Method) String() string {
	return string(h)
}