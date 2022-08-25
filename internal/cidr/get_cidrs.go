package cidr

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetCIDRList(uris []string) (string, error) {
	var cidrList []string
	// TODO: Fix this right but for now let's just ignore it
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for i := range uris {
		res, err := client.Get(uris[i])
		if err != nil {
			return "", err
		}

		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("didn't receive ok from server: %d", res.StatusCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}

		s := strings.Fields(string(resBody))
		cidrList = append(cidrList, strings.Join(s, ","))

	}

	return strings.Join(cidrList, ","), nil
}
