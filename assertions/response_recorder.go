package assertions

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/assertions"
	"net/http"
	"net/http/httptest"
)

func ShouldHaveEmptyBody(actual interface{}, expected ...interface{}) string {
	resp := actual.(*httptest.ResponseRecorder)
	if resp.Body.Len() == 0 {
		return ""
	} else {
		return fmt.Sprintf("The body is not empty! %s", resp.Body.String())
	}
}

func ShouldHaveResponseCode(actual interface{}, expected ...interface{}) string {
	resp := actual.(*httptest.ResponseRecorder)
	if resp.Code == expected[0] {
		return ""
	} else {
		return fmt.Sprintf("Unexpected response code! Expected: %v Actual: %v", expected[0], resp.Code)
	}
}

func ShouldContainHeader(actual interface{}, expected ...interface{}) string {
	resp := actual.(*httptest.ResponseRecorder)
	expectedHeader := expected[0].(http.Header)
	for key, _ := range expectedHeader {
		actualValue := resp.Header().Get(key)
		if actualValue == "" {
			return fmt.Sprintf("Missing header '%v' !", key)
		}
		if actualValue != expectedHeader.Get(key) {
			return fmt.Sprintf("Header '%v' is present but has an unexpected value! %v \n %v", key, actualValue, expectedHeader.Get(key))
		}
	}
	return ""
}

func ShouldHaveJSONBody(actual interface{}, expected ...interface{}) string {
	resp := actual.(*httptest.ResponseRecorder)

	switch expected[0].(type) {
	default:
		return "Not supported"
	case []map[string]interface{}:
		expectedArray := expected[0].([]map[string]interface{})
		actualArray := make([]map[string]interface{}, 0)
		json.Unmarshal(resp.Body.Bytes(), &actualArray)
		if err := assertions.ShouldHaveLength(actualArray, len(expectedArray)); err != "" {
			return err
		}
		for _, actualElement := range actualArray {
			if err := assertions.ShouldContain(expectedArray, actualElement); err != "" {
				return err
			}
		}
		return ""
	case map[string]interface{}:
		expectedMap := expected[0].(map[string]interface{})
		actualMap := make(map[string]interface{})
		json.Unmarshal(resp.Body.Bytes(), &actualMap)
		return assertions.ShouldResemble(expectedMap, actualMap)
	}
}
