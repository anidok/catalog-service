package opensearch

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
)

type mockTransport struct {
	resp *http.Response
	err  error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func newMockClient(respBody map[string]interface{}, statusCode int) *ClientImpl {
	body, _ := json.Marshal(respBody)
	resp := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}
	httpClient := &http.Client{
		Transport: &mockTransport{resp: resp},
	}
	osClient, _ := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://mock:9200"},
		Transport: httpClient.Transport,
	})
	return &ClientImpl{Client: osClient}
}

func unmarshalJSON(jsonStr string) map[string]interface{} {
	var result map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &result)
	return result
}
