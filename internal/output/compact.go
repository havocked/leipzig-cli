package output

import (
	"fmt"
	"io"

	"github.com/havocked/leipzig-cli/internal/model"
)

func Compact(w io.Writer, events []model.Event) error {
	if len(events) == 0 {
		fmt.Fprintln(w, "No events found.")
		return nil
	}
	for _, e := range events {
		fmt.Fprintln(w, e.String())
	}
	return nil
}
