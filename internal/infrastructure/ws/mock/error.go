package mock

import "errors"

var ErrTopicNotFound = errors.New("topic not found")

var ErrTopicAlreadyExists = errors.New("topic already exist")

var ErrConnectionClosed = errors.New("connection closed")