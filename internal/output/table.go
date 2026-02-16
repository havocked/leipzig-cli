package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/havocked/leipzig-cli/internal/model"
)

func Table(w io.Writer, events []model.Event) error {
	if len(events) == 0 {
		fmt.Fprintln(w, "No events found.")
		return nil
	}

	for _, e := range events {
		date := e.StartTime.Format("Mon 02 Jan")
		timeStr := e.StartTime.Format("15:04")
		if e.StartTime.Hour() == 0 && e.StartTime.Minute() == 0 {
			timeStr = "     "
		}
		cat := fmt.Sprintf("%-10s", e.Category)
		venue := e.Venue
		if len(venue) > 25 {
			venue = venue[:22] + "..."
		}
		name := e.Name
		if len(name) > 40 {
			name = name[:37] + "..."
		}
		price := e.Price
		if price == "" {
			price = ""
		}
		_ = strings.TrimSpace
		fmt.Fprintf(w, "%-11s %5s  %s  %-40s  %-25s  %s\n", date, timeStr, cat, name, venue, price)
	}
	return nil
}
