package documents

import (
	"testing"
)

func TestPagesNeeded(t *testing.T) {
	testCases := map[string]struct {
		itemsLength   int
		expectedPages int
	}{
		"fitsOnFirstPage":         {itemsLength: 8, expectedPages: 1},
		"oneExtraItemForTwoPages": {itemsLength: 9, expectedPages: 2},
		"fitsOnTwoPages":          {itemsLength: 20, expectedPages: 2},
		"oneMoreThanTwoPages":     {itemsLength: 21, expectedPages: 3},
		"zeroItems":               {itemsLength: 0, expectedPages: 1},
		"threePages":              {itemsLength: 32, expectedPages: 3},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualPages := pagesNeeded(tc.itemsLength)
			if actualPages != tc.expectedPages {
				t.Errorf("For itemsLength %d, expected %d pages, got %d pages", tc.itemsLength, tc.expectedPages, actualPages)
			}
		})
	}
}
