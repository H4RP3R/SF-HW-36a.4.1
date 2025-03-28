package rss

import "time"

var rssTimeFormats = []string{
	"Mon, 2 Jan 2006 15:04:05 +0000", // Custom for single digit date.
	time.RFC1123,                     // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,                    // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.RFC822,                      // "02 Jan 06 15:04 MST"
	time.RFC822Z,                     // "02 Jan 06 15:04 -0700"
}

func ConvertToUTC(rssTime string) (t time.Time, err error) {
	for _, format := range rssTimeFormats {
		t, err = time.Parse(format, rssTime)
		if err == nil {
			return t.UTC(), nil
		}
	}
	return
}
