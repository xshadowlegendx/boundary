// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package credentialstores

import (
	"strconv"
	"strings"

	"github.com/hashicorp/boundary/api"
)

// Option is a func that sets optional attributes for a call. This does not need
// to be used directly, but instead option arguments are built from the
// functions in this package. WithX options set a value to that given in the
// argument; DefaultX options indicate that the value should be set to its
// default. When an API call is made options are processed in ther order they
// appear in the function call, so for a given argument X, a succession of WithX
// or DefaultX calls will result in the last call taking effect.
type Option func(*options)

type options struct {
	postMap                 map[string]interface{}
	queryMap                map[string]string
	withAutomaticVersioning bool
	withSkipCurlOutput      bool
	withFilter              string
	withRecursive           bool
}

func getDefaultOptions() options {
	return options{
		postMap:  make(map[string]interface{}),
		queryMap: make(map[string]string),
	}
}

func getOpts(opt ...Option) (options, []api.Option) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o != nil {
			o(&opts)
		}
	}
	var apiOpts []api.Option
	if opts.withSkipCurlOutput {
		apiOpts = append(apiOpts, api.WithSkipCurlOutput(true))
	}
	if opts.withFilter != "" {
		opts.queryMap["filter"] = opts.withFilter
	}
	if opts.withRecursive {
		opts.queryMap["recursive"] = strconv.FormatBool(opts.withRecursive)
	}
	return opts, apiOpts
}

// If set, and if the version is zero during an update, the API will perform a
// fetch to get the current version of the resource and populate it during the
// update call. This is convenient but opens up the possibility for subtle
// order-of-modification issues, so use carefully.
func WithAutomaticVersioning(enable bool) Option {
	return func(o *options) {
		o.withAutomaticVersioning = enable
	}
}

// WithSkipCurlOutput tells the API to not use the current call for cURL output.
// Useful for when we need to look up versions.
func WithSkipCurlOutput(skip bool) Option {
	return func(o *options) {
		o.withSkipCurlOutput = true
	}
}

// WithFilter tells the API to filter the items returned using the provided
// filter term.  The filter should be in a format supported by
// hashicorp/go-bexpr.
func WithFilter(filter string) Option {
	return func(o *options) {
		o.withFilter = strings.TrimSpace(filter)
	}
}

// WithRecursive tells the API to use recursion for listing operations on this
// resource
func WithRecursive(recurse bool) Option {
	return func(o *options) {
		o.withRecursive = true
	}
}

func WithVaultCredentialStoreAddress(inAddress string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["address"] = inAddress
		o.postMap["attributes"] = val
	}
}

func WithAttributes(inAttributes map[string]interface{}) Option {
	return func(o *options) {
		o.postMap["attributes"] = inAttributes
	}
}

func DefaultAttributes() Option {
	return func(o *options) {
		o.postMap["attributes"] = nil
	}
}

func WithVaultCredentialStoreCaCert(inCaCert string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["ca_cert"] = inCaCert
		o.postMap["attributes"] = val
	}
}

func DefaultVaultCredentialStoreCaCert() Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["ca_cert"] = nil
		o.postMap["attributes"] = val
	}
}

func WithVaultCredentialStoreClientCertificate(inClientCertificate string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["client_certificate"] = inClientCertificate
		o.postMap["attributes"] = val
	}
}

func DefaultVaultCredentialStoreClientCertificate() Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["client_certificate"] = nil
		o.postMap["attributes"] = val
	}
}

func WithVaultCredentialStoreClientCertificateKey(inClientCertificateKey string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["client_certificate_key"] = inClientCertificateKey
		o.postMap["attributes"] = val
	}
}

func DefaultVaultCredentialStoreClientCertificateKey() Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["client_certificate_key"] = nil
		o.postMap["attributes"] = val
	}
}

func WithDescription(inDescription string) Option {
	return func(o *options) {
		o.postMap["description"] = inDescription
	}
}

func DefaultDescription() Option {
	return func(o *options) {
		o.postMap["description"] = nil
	}
}

func WithName(inName string) Option {
	return func(o *options) {
		o.postMap["name"] = inName
	}
}

func DefaultName() Option {
	return func(o *options) {
		o.postMap["name"] = nil
	}
}

func WithVaultCredentialStoreNamespace(inNamespace string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["namespace"] = inNamespace
		o.postMap["attributes"] = val
	}
}

func DefaultVaultCredentialStoreNamespace() Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["namespace"] = nil
		o.postMap["attributes"] = val
	}
}

func WithVaultCredentialStoreTlsServerName(inTlsServerName string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["tls_server_name"] = inTlsServerName
		o.postMap["attributes"] = val
	}
}

func DefaultVaultCredentialStoreTlsServerName() Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["tls_server_name"] = nil
		o.postMap["attributes"] = val
	}
}

func WithVaultCredentialStoreTlsSkipVerify(inTlsSkipVerify bool) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["tls_skip_verify"] = inTlsSkipVerify
		o.postMap["attributes"] = val
	}
}

func DefaultVaultCredentialStoreTlsSkipVerify() Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["tls_skip_verify"] = nil
		o.postMap["attributes"] = val
	}
}

func WithVaultCredentialStoreToken(inToken string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["token"] = inToken
		o.postMap["attributes"] = val
	}
}

func WithVaultCredentialStoreWorkerFilter(inWorkerFilter string) Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["worker_filter"] = inWorkerFilter
		o.postMap["attributes"] = val
	}
}

func DefaultVaultCredentialStoreWorkerFilter() Option {
	return func(o *options) {
		raw, ok := o.postMap["attributes"]
		if !ok {
			raw = interface{}(map[string]interface{}{})
		}
		val := raw.(map[string]interface{})
		val["worker_filter"] = nil
		o.postMap["attributes"] = val
	}
}
