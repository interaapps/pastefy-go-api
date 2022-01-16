package pastefy

import "errors"

type MultiPastePart struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}

type Paste struct {
	Id        string `json:"id"`
	Created   string `json:"created"`
	Encrypted bool   `json:"encrypted"`
	UserId    string `json:"user_id,omitempty"`
	Exists    bool   `json:"exists"`
	Title     string `json:"title"`
	Type      string `json:"type,omitempty"`
	FolderId  string `json:"folder,omitempty"`
	Content   string `json:"content"`
	RawUrl    string `json:"raw_url,omitempty"`
}

func (apiClient PastefyApiClient) GetPaste(id string) (Paste, error) {
	paste := Paste{}
	_, err := apiClient.RequestMap("GET", "/paste/"+id, nil, &paste)
	if err != nil {
		return Paste{}, err
	}
	return paste, nil
}
func (apiClient PastefyApiClient) CreatePaste(paste Paste) (*Paste, error) {
	ret := struct {
		Paste Paste `json:"paste"`
	}{}
	_, err := apiClient.RequestMap("POST", "/paste", paste, &ret)
	if err != nil {
		return nil, errors.New("error while fetching paste")
	}
	return &ret.Paste, nil
}

func (apiClient PastefyApiClient) UpdatePaste(id string, paste Paste) (*Paste, error) {
	ret := struct {
		Paste Paste `json:"paste"`
	}{}
	_, err := apiClient.RequestMap("PUT", "/paste/"+id, paste, &ret)

	if err != nil {
		return nil, errors.New("error while fetching paste")
	}
	return &ret.Paste, nil
}

func (apiClient PastefyApiClient) DeletePaste(id string) error {
	_, err := apiClient.Request("DELETE", "/paste/"+id, nil)
	return err
}

func (apiClient PastefyApiClient) AddFriendToPaste(paste string, friend string) error {
	_, err := apiClient.Request("DELETE", "/paste/"+paste+"/friend", struct {
		Friend string `json:"friend"`
	}{
		Friend: friend,
	})
	return err
}

func (apiClient PastefyApiClient) SavePaste(paste Paste) (*Paste, error) {
	if paste.Id != "" {
		return apiClient.UpdatePaste(paste.Id, paste)
	} else {
		return apiClient.CreatePaste(paste)
	}
}

func (apiClient PastefyApiClient) GetUserPastes() ([]Paste, error) {
	var pastes []Paste
	_, err := apiClient.RequestMap("GET", "/user/pastes", nil, &pastes)
	if err != nil {
		return nil, err
	}
	return pastes, nil
}

func (apiClient PastefyApiClient) GetUserSharedPastes() ([]Paste, error) {
	var pastes []Paste
	_, err := apiClient.RequestMap("GET", "/user/sharedpastes", nil, &pastes)
	if err != nil {
		return nil, err
	}
	return pastes, nil
}

func (paste Paste) Encrypt(password string) (Paste, error) {
	paste.Encrypted = true
	title, err := AesEncrypt(paste.Title, password)
	if err != nil {
		return Paste{}, err
	}
	content, err := AesEncrypt(paste.Content, password)
	if err != nil {
		return Paste{}, err
	}
	paste.Title = title
	paste.Content = content
	return paste, nil
}

func (paste Paste) Decrypt(password string) (Paste, error) {
	paste.Encrypted = false
	if paste.Title != "" {
		title, err := AesDecrypt(paste.Title, password)
		if err != nil {
			return Paste{}, err
		}
		paste.Title = title
	}
	if paste.Content != "" {
		content, err := AesDecrypt(paste.Content, password)
		if err != nil {
			return Paste{}, err
		}
		paste.Content = content
	}
	return paste, nil
}
