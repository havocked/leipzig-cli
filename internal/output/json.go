package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/havocked/leipzig-cli/internal/model"
)

func JSON(w io.Writer, events []model.Event) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(events); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
