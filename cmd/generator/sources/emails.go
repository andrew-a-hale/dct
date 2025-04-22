package sources

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

func GenerateEmails() []string {
	emails := []string{}
	emailSanitiser := strings.NewReplacer(" ", "", "@", "")

	for range 200 {
		firstNameIdx := rand.IntN(len(FirstNames))
		firstName := FirstNames[firstNameIdx].Name
		lastNameIdx := rand.IntN(len(LastNames))
		lastName := LastNames[lastNameIdx]
		companiesIdx := rand.IntN(len(Companies))
		company := Companies[companiesIdx]
		emails = append(emails, fmt.Sprintf("%s.%s@%s.com", firstName, lastName, emailSanitiser.Replace(company)))
	}

	return emails
}
