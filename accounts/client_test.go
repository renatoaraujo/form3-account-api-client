package accounts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAccountIDFormat(t *testing.T) {
	tests := []struct {
		name      string
		accountID string
		wantErr   bool
	}{
		{
			name:      "Failed due an invalid account id format",
			accountID: "this is an invalid uuid format",
			wantErr:   true,
		},
		{
			name:      "Successfully validate the account id format",
			accountID: "00000000-0000-0000-0000-000000000000",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAccountIDFormat(tt.accountID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
