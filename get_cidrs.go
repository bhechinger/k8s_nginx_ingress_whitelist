package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getCIDRList(uris []string) (string, error) {
	var cidrList []string

	for i := range uris {
		res, err := http.Get(uris[i])
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
