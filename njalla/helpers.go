package njalla

import (
	"fmt"
	"strconv"
	"strings"
)

// parseImportID will parse a given resource ID when importing with the
// following format: `domain:id`, where `domain` will be any string, and `id`
// will be numeric.
func parseImportID(id string) (string, int, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		msg := fmt.Errorf(
			"unexpected format of ID (%s), expected domain:id", id,
		)
		return "", 0, msg
	}

	parsedID, err := strconv.Atoi(parts[1])
	if err != nil {
		msg := fmt.Errorf("expected id to be numeric, got: %s", parts[1])
		return "", 0, msg
	}

	return parts[0], parsedID, nil
}
