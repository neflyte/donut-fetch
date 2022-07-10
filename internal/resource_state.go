package internal

import (
	"fmt"
	"time"
)

type ResourceState struct {
	LastModified time.Time `json:"last-modified,omitempty"`
	ETag         string    `json:"etag,omitempty"`
}

func (rs ResourceState) String() string {
	return fmt.Sprintf("LastModified=%s,ETag=%s", rs.LastModified.String(), rs.ETag)
}

func (rs ResourceState) IsEmpty() bool {
	return rs.LastModified == time.Time{} && rs.ETag == ""
}

func (rs ResourceState) IsETagStale(sourceETag string) bool {
	return sourceETag != rs.ETag
}

func (rs ResourceState) IsLastModifiedPast(sourceLastModified time.Time) bool {
	return sourceLastModified.After(rs.LastModified)
}
