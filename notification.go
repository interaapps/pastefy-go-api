package pastefy

type Notifcation struct {
	UpdatedAt   string `json:"updated_at"`
	UserId      string `json:"user_id"`
	AlreadyRead bool   `json:"already_read"`
	CreatedAt   string `json:"created_at"`
	Received    bool   `json:"received"`
	Id          int    `json:"id"`
	Message     string `json:"message"`
	Url         string `json:"url"`
}

func (apiClient PastefyApiClient) GetUserNotifications() ([]Notifcation, error) {
	var notifications []Notifcation
	_, err := apiClient.RequestMap("GET", "/user/notification", nil, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (apiClient PastefyApiClient) MarkAllNotificationsAsRead() error {
	_, err := apiClient.Request("GET", "/user/readall", nil)
	return err
}
