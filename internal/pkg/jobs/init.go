package jobs

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/manifest"
)

type aentBootstrap struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// NewInitJob creates a new Job instance.
func NewInitJob(ctx *context.Context, m *manifest.Manifest) (Job, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/theaentmachine/colony-registry/master/bootstrap.json")
	if err != nil {
		return nil, errors.Wrap("init job", err)
	}
	defer resp.Body.Close()
	var aent aentBootstrap
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if err := json.Unmarshal(buf.Bytes(), &aent); err != nil {
		return nil, errors.Wrap("init job", err)
	}
	return NewAddJob(aent.Image, ctx, m)
}
