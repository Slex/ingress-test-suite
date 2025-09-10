package runner

import (
	"fmt"
	"ingress-test-suite/pkg/messages"
	"ingress-test-suite/test_load"
	"net/http"
)

type HTTPTester struct{}

func (t *HTTPTester) Test(entry test_load.IngressTestEntry) TestResult {
	result := TestResult{Host: entry.Host, Path: entry.Path}

	url := fmt.Sprintf("http://%s:%d%s", entry.Host, entry.ExtPort, entry.Path)
	log.Printf(messages.RequestURL, url)

	resp, err := http.Get(url)

	if err != nil {
		result.Success = false
		result.ErrorMessage = fmt.Sprintf(messages.HttpRequestFailed, err)
		return result
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf(messages.FailedCloseResponseBody, cerr)
		}
	}()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode == entry.ExpectedStatus
	return result
}
