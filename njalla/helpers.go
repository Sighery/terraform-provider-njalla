package njalla

import (
	"fmt"
	"strings"
)

// parseImportID will parse a given resource ID when importing with the
// following format: `domain:id`, where `domain` and `id` will be any string.
func parseImportID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		msg := fmt.Errorf(
			"unexpected format of ID (%s), expected domain:id", id,
		)
		return "", "", msg
	}

	return parts[0], parts[1], nil
}
