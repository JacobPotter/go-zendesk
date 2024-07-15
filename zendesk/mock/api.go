package mock

import (
	"github.com/JacobPotter/go-zendesk/zendesk"
)

var _ zendesk.API = (*Client)(nil)
