package config

import "errors"

var ErrMissingCredentials = errors.New("missing API key or access token")