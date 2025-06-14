package alpaca

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AlpacaAccount struct {
	Id string `json:"id"`
}

func CreateAlpacaAccount(email string, firstName string, lastName string) (*AlpacaAccount, error) {

	url := "https://broker-api.sandbox.alpaca.markets/v1/accounts"

	payload := strings.NewReader(fmt.Sprintf("{\"contact\":{\"email_address\":\"%s\",\"phone_number\":\"+15556667788\",\"city\":\"San Mateo\",\"postal_code\":\"94401\",\"street_address\":[\"20 N San Mateo Dr\"],\"country\":\"USA\",\"state\":\"CA\"},\"identity\":{\"tax_id_type\":\"NOT_SPECIFIED\",\"given_name\":\"%s\",\"family_name\":\"%s\",\"date_of_birth\":\"1990-01-01\",\"tax_id\":\"666-55-4321\",\"funding_source\":[\"employment_income\"],\"country_of_tax_residence\":\"USA\",\"country_of_citizenship\":\"USA\",\"country_of_birth\":\"USA\"},\"disclosures\":{\"is_control_person\":false,\"is_affiliated_exchange_or_finra\":false,\"is_politically_exposed\":false,\"immediate_family_exposed\":false},\"agreements\":[{\"agreement\":\"account_agreement\",\"signed_at\":\"2019-09-11T18:09:33Z\",\"ip_address\":\"185.13.21.99\"},{\"agreement\":\"margin_agreement\",\"signed_at\":\"2019-09-11T18:09:33Z\",\"ip_address\":\"185.13.21.99\"},{\"agreement\":\"customer_agreement\",\"signed_at\":\"2019-09-11T18:09:33Z\",\"ip_address\":\"185.13.21.99\"}],\"account_type\":\"trading\"}", email, firstName, lastName))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Basic Q0tKSzM2UzNFRTZPTDRBWTU1SDE6THJkYXFpaUQyZzRGeHV4cG1rOE1yWmxibE5qYTJVOU5HR0lsQ0tGMg==")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var account_id AlpacaAccount

	if err := json.Unmarshal(body, &account_id); err != nil {
		return nil, err
	}

	return &account_id, nil

}
