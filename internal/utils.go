package internal

import "github.com/gofrs/uuid"

const (
	uuidV5Seed = "fd2e3c13-06f0-4ee9-9a0d-6b6ac5cf0eaa"
)

var (
	uuidV5NS = uuid.Must(uuid.FromString(uuidV5Seed))
)

func HashURL(url string) string {
	return uuid.NewV5(uuidV5NS, url).String()
}
