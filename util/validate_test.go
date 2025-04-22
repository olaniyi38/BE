package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidEmail(t *testing.T) {
	emails := []string{
		"john.doe@example.com",
		"sales+2024@company.co.uk",
		"invalid.email@",
		"noatsign.com",
	}

	validEmails := []string{}
	invalidEmails := []string{}

	for _, email := range emails {
		isValid := IsValidEmail(email)
		if isValid {
			validEmails = append(validEmails, email)
		} else {
			invalidEmails = append(invalidEmails, email)
		}

	}

	require.Equal(t, validEmails, []string{"john.doe@example.com", "sales+2024@company.co.uk"})
	require.Equal(t, invalidEmails, []string{"invalid.email@", "noatsign.com"})
}
