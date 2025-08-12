package userclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserClient interface {
	FindById(it uint) (*UserDTO, error)
}

type userClient struct {
	baseURL string
}

func NewUserClient() UserClient {
	return &userClient{
		baseURL: "http://user-service:8080/api", // TODO: This should not be hardcoded
	}
}

func (c *userClient) FindById(id uint) (*UserDTO, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", c.baseURL, id))
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var obj UserDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}
