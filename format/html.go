package format

import (
	"chess-profile/types"
	"fmt"
	"sort"
	"strings"
)

func Format(profile types.Profile) string {
	html := strings.Replace(template, "%USERNAME%", profile.User, -1)
	html = strings.Replace(html, "%WHITE_OPENINGS%", formatOpenings(profile.OpeningsWhite), -1)
	html = strings.Replace(html, "%BLACK_OPENINGS%", formatOpenings(profile.OpeningsBlack), -1)
	html = strings.Replace(html, "%MATE_WIN_RATIO%", fmt.Sprintf("%.2f", profile.MateWinRatio), -1)
	html = strings.Replace(html, "%RESIGN_LOSS_RATIO%", fmt.Sprintf("%.2f", profile.ResignLossRatio), -1)
	html = strings.Replace(html, "%DURATION_PERCENTILES%", formatDurationPercentiles(profile.DurationPercentiles), -1)
	return html
}

func formatOpenings(openings []string) string {
	return strings.Join(openings, "<br>")
}

func formatDurationPercentiles(percentiles map[int]int) string {
	lines := make([]string, 0, len(percentiles))
	for pc, val := range percentiles {
		lines = append(lines, fmt.Sprintf("%d%%: %d", pc, val))
	}
	sort.Strings(lines)
	return strings.Join(lines, "<br>")
}

var template = `
<!DOCTYPE html>
<html>

<head>
  <title>%USERNAME%</title>
</head>

<body>
  <h2>%USERNAME%</h2>
  <strong>Favourite openings (White):</strong>
  <br>
  %WHITE_OPENINGS%
  <br>
  <strong>Favourite openings (Black):</strong>
  <br>
  %BLACK_OPENINGS%
  <br>
  <strong>Mate/Win:</strong> %MATE_WIN_RATIO% <br>
  <strong>Resign/Loss:</strong> %RESIGN_LOSS_RATIO% <br>
  <strong>Duration percentiles:</strong> <br> %DURATION_PERCENTILES% <br>
</body>

</html>
`
