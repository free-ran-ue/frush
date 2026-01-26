package subscriber

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/util"
)

func getConsoleLoginToken(consoleAccountTemplatePath string) (string, error) {
	content, err := util.FileRead(consoleAccountTemplatePath)
	if err != nil {
		return "", err
	}

	headers := make(map[string]string)
	headers[constant.HTTP_HEADER_CONTENT_TYPE] = constant.HTTP_HEADER_CONTENT_TYPE_JSON

	response, err := util.SendHttpRequest(
		fmt.Sprintf("http://%s:%d%s", constant.CONSOLE_IP, constant.CONSOLE_PORT, constant.CONSOLE_LOGIN_PATH),
		http.MethodPost,
		headers,
		content,
	)
	if err != nil {
		return "", err
	}

	var bodyJSON map[string]interface{}
	if err := json.Unmarshal(response.Body, &bodyJSON); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %v", err)
	}

	tokenValue, exists := bodyJSON[constant.CONSOLE_ACCESS_TOKEN]
	if !exists {
		return "", fmt.Errorf("access token not found in response")
	}
	return tokenValue.(string), nil
}

func subscriberMain(token, subscriberTemplatePath, action string) error {
	content, err := util.FileRead(subscriberTemplatePath)
	if err != nil {
		return err
	}

	contentJSON := make(map[string]interface{})
	if err := json.Unmarshal(content, &contentJSON); err != nil {
		return fmt.Errorf("failed to parse subscriber template JSON: %v", err)
	}

	ueId, ok := contentJSON["ueId"]
	if !ok || ueId == nil {
		return fmt.Errorf("ueId not found in subscriber template")
	}
	imsi := ueId.(string)
	if len(imsi) > 5 && imsi[:5] == "imsi-" {
		imsi = imsi[5:]
	}

	plmnID, ok := contentJSON["plmnID"].(string)
	if !ok || plmnID == "" {
		return fmt.Errorf("plmnID not found or invalid")
	}

	headers := map[string]string{
		constant.HTTP_HEADER_CONTENT_TYPE: constant.HTTP_HEADER_CONTENT_TYPE_JSON,
		constant.CONSOLE_TOKEN:            token,
	}

	response, err := util.SendHttpRequest(
		fmt.Sprintf("http://%s:%d%s", constant.CONSOLE_IP, constant.CONSOLE_PORT, fmt.Sprintf(constant.CONSOLE_ADD_SUBSCRIBER_PATH, imsi, plmnID)),
		action,
		headers,
		content,
	)
	if err != nil {
		return err
	}

	switch action {
	case http.MethodPost:
		if response.StatusCode == http.StatusConflict {
			return fmt.Errorf("subscriber already exists in webconsole")
		}
		if response.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to add subscriber: %d", response.StatusCode)
		}
	case http.MethodDelete:
		if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("subscriber not found in webconsole")
		}
		if response.StatusCode != http.StatusNoContent {
			return fmt.Errorf("failed to delete subscriber: %d", response.StatusCode)
		}
	}

	return nil
}

func AddSubscriber(consoleAccountTemplatePath, subscriberTemplatePath string) error {
	token, err := getConsoleLoginToken(consoleAccountTemplatePath)
	if err != nil {
		return err
	}

	return subscriberMain(token, subscriberTemplatePath, http.MethodPost)
}

func DeleteSubscriber(consoleAccountTemplatePath, subscriberTemplatePath string) error {
	token, err := getConsoleLoginToken(consoleAccountTemplatePath)
	if err != nil {
		return err
	}

	return subscriberMain(token, subscriberTemplatePath, http.MethodDelete)
}