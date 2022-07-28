package repository

import (
	"context"
	"course/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const getUserUrl = "internal/users/"

type MicroserviceRepo struct {
	hostname string
	username string
	password string
	client   *http.Client
}

func NewMicroserviceRepo() *MicroserviceRepo {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	return &MicroserviceRepo{
		hostname: os.Getenv("API_HOST"),
		username: os.Getenv("API_USERNAME"),
		password: os.Getenv("API_PASSWORD"),
		client:   &client,
	}
}

func (mr MicroserviceRepo) IsUserExists(ctx context.Context, userID int) bool {
	url := fmt.Sprintf("%s%s%d", mr.hostname, getUserUrl, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	req.SetBasicAuth(mr.username, mr.password)
	resp, err := mr.client.Do(req)
	if err != nil {
		return false
	}
	var user domain.User
	json.NewDecoder(resp.Body).Decode(&user)
	return user.ID > 0
}
