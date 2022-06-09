package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zachary-cauchi/golang-angular-sample-app/internal/message"
)

func GetMessageListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, message.Get())
}

func AddMessageHandler(c *gin.Context) {
	messageItem, statusCode, err := convertHTTPBodyToMessage(c.Request.Body)

	if err != nil {
		c.JSON(statusCode, err)

		return
	}

	c.JSON(statusCode, gin.H{"id": message.Add(messageItem.Text)})
}

func DeleteMessageHandler(c *gin.Context) {
	messageId := c.Param("id")

	if err := message.Delete(messageId); err != nil {
		c.JSON(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToMessage(httpBody io.ReadCloser) (message.Message, int, error) {
	body, err := ioutil.ReadAll(httpBody)

	if err != nil {
		return message.Message{}, http.StatusInternalServerError, err
	}

	defer httpBody.Close()

	return convertJSONBodyToMessage(body)
}

func convertJSONBodyToMessage(jsonBody []byte) (message.Message, int, error) {
	var messageItem message.Message

	err := json.Unmarshal(jsonBody, &messageItem)

	if err != nil {
		return message.Message{}, http.StatusBadRequest, err
	}

	return messageItem, http.StatusOK, nil
}
