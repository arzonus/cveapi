package interfaces

import "testing"
import "github.com/arzonus/cveapi/usecases"

func TestFindCVEs(t *testing.T) {
	ucves := []usecases.CVE{
		{
			Id:      "CVE-2016-0746",
			Summary: "123",
			CPEs: []string{
				"cpe:/a:1024cms:1024_cms:0.7",
				"cpe:/a:1024cms:1024_cms:0.8",
			},
		},
		{
			Id:      "CVE-2016-0747",
			Summary: "123",
			CPEs: []string{
				"cpe:/a:1024cms:1024_cms:0.7",
				"cpe:/a:1024cms:1024_cms:0.8",
			},
		},
	}

	cpe := "cpe:/a:1024cms:1024_cms:0.7"

	cves := []string{
		"CVE-2016-0746",
		"CVE-2016-0747",
	}

	rcves := findCVEs(cpe, ucves)
	if arrayStringEq(cves, rcves) != true {
		t.Fatalf("Expected ", cves, " got ", rcves)
	}
}

func arrayStringEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
