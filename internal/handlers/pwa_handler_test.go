package handlers

import (
	"testing"
)

func TestValidateDocumentRequest(t *testing.T) {
	testCases := map[string]struct {
		documentType string
		documentID   string
		documentDate string
		expectError  bool
	}{
		"emptyDocumentType":    {documentType: "", documentID: "123", documentDate: "2023-11-01", expectError: true},
		"monthlyMissingDate":   {documentType: "monthly", documentID: "", documentDate: "", expectError: true},
		"monthlyValidDate":     {documentType: "monthly", documentID: "", documentDate: "2023-11-01", expectError: false},
		"otherTypeMissingID":   {documentType: "invoice", documentID: "", documentDate: "2023-11-01", expectError: true},
		"otherTypeWithValidID": {documentType: "invoice", documentID: "123", documentDate: "2023-11-01", expectError: false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateDocumentRequest(tc.documentType, tc.documentID, tc.documentDate)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}
