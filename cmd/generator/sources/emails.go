package sources

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var Emails [200]string

func init() {
	emailSanitiser := strings.NewReplacer(" ", "", "@", "")

	for i := range 200 {
		firstNameIdx := rand.IntN(len(FirstNames))
		firstName := FirstNames[firstNameIdx]
		lastNameIdx := rand.IntN(len(LastNames))
		lastName := LastNames[lastNameIdx]
		companiesIdx := rand.IntN(len(Companies))
		company := Companies[companiesIdx]
		Emails[i] = fmt.Sprintf("%s.%s@%s.com", firstName, lastName, emailSanitiser.Replace(company))
	}
}
