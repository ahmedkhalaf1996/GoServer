package countrycodes

import (
	"testing"
)

func TestFindByName(t *testing.T) {
	matches := FindByName("United States Minor")

	if len(matches) != 1 {
		t.Fatalf("Extra matches found")
	}

	um, _ := GetByAlpha2("UM")
	if matches[0] != um {
		t.Fatalf("Match for united States Minor outlyin Islands faild")
	}
}

func TestGetByNumeric(t *testing.T) {
	code, _ := GetByNumeric(840)

	if code.Name != "United States" {
		t.Fatalf("GetByNumeric faild")
	}
}
