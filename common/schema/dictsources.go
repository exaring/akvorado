package schema

import (
	"io"
	"net/http"
	"os"

	"akvorado/common/helpers"
	"akvorado/common/s3"
)

type DictSource interface {
	Get(key string) (io.ReadCloser, error)
}

// DictSourceConfiguration represents the configuration of a cache backend.
type DictSourceConfiguration interface {
	New(c *s3.Component) (DictSource, error)
}

// FileDictSourceConfiguration is the configuration for a dict source reading
// from the local filesystem. There is no configuration.
type FileDictSourceConfiguration struct{}

// New creates a new file dict source from a file dict source configuration.
func (FileDictSourceConfiguration) New(c *s3.Component) (DictSource, error) {
	return FileDictSource{}, nil
}

// FileDictSource is a dict source reading from the local filesystem.
type FileDictSource struct{}

// Get returns a file reader for the specified filename.
func (FileDictSource) Get(key string) (io.ReadCloser, error) {
	return os.Open(key)
}

// DefaultFileDictSourceConfiguration returns the default configuration for a
// filesystem based dict source.
func DefaultFileDictSourceConfiguration() DictSourceConfiguration {
	return FileDictSourceConfiguration{}
}

// HttpDictSourceConfiguration is the configuration for a dict source reading
// from the specified HTTP endpoint.
type HttpDictSourceConfiguration struct {
	BaseURL string
}

// New creates a new HTTP dict source from a HTTP dict source configuration.
func (HttpDictSourceConfiguration) New(c *s3.Component) (DictSource, error) {
	return HttpDictSource{}, nil
}

// HttpDictSource is a dict source reading from the specified HTTP endpoint.
type HttpDictSource struct{}

// Get returns a file reader for the specified URL.
func (HttpDictSource) Get(key string) (io.ReadCloser, error) {
	resp, err := http.Get(key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// DefaultHttpDictSourceConfiguration returns the default configuration for a
// HTTP dict source configuration.
func DefaultHttpDictSourceConfiguration() DictSourceConfiguration {
	return HttpDictSourceConfiguration{}
}

// S3DictSourceConfiguration is the configuration for a dict source reading
// from the specified S3 config.
type S3DictSourceConfiguration struct {
	S3Config string
}

// New creates a new S3 dict source from an S3 dict source configuration.
func (sc S3DictSourceConfiguration) New(c *s3.Component) (DictSource, error) {
	return S3DictSource{config: sc, c: c}, nil
}

// S3DictSource is a dict source reading from the specified S3 bucket.
type S3DictSource struct {
	config S3DictSourceConfiguration
	c      *s3.Component
}

// Get returns a file reader for the specified S3 object.
func (s S3DictSource) Get(key string) (io.ReadCloser, error) {
	read, err := s.c.GetObject(s.config.S3Config, key)
	if err != nil {
		return nil, err
	}

	return read, nil
}

// DefaultS3DictSourceConfiguration returns the default configuration for a
// S3 dict source.
func DefaultS3DictSourceConfiguration() DictSourceConfiguration {
	return S3DictSourceConfiguration{}
}

// S3MockDictSourceConfiguration is the configuration for a dict source reading
// from a mocked S3 endpoint.
type S3MockDictSourceConfiguration struct {
}

// New creates a new S3 dict source from an S3 dict source configuration.
func (sc S3MockDictSourceConfiguration) New(c *s3.Component) (DictSource, error) {
	panic("not implemented")
}

// DefaultS3MockDictSourceConfiguration returns the default configuration for a
// mocked S3 dict source.
func DefaultS3MockDictSourceConfiguration() DictSourceConfiguration {
	return S3MockDictSourceConfiguration{}
}

var dictSourceConfigurationMap = map[string](func() DictSourceConfiguration){
	"file":   DefaultFileDictSourceConfiguration,
	"http":   DefaultHttpDictSourceConfiguration,
	"s3":     DefaultS3DictSourceConfiguration,
	"s3mock": DefaultS3MockDictSourceConfiguration,
}

func init() {
	helpers.RegisterMapstructureUnmarshallerHook(
		helpers.ParametrizedConfigurationUnmarshallerHook(CustomDict{}, dictSourceConfigurationMap))
}

// MarshalYAML undoes ConfigurationUnmarshallerHook().
func (cc CustomDict) MarshalYAML() (interface{}, error) {
	return helpers.ParametrizedConfigurationMarshalYAML(cc, dictSourceConfigurationMap)
}
