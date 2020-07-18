package support

// MIT Licensed - see LICENSE
// Copyright (C) 2014 Philip Schlump

import (
	"github.com/pschlump/uuid"
)

func GenUUID() string {
	newUUID, _ := uuid.NewV4()
	return newUUID.String()
}
