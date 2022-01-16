package pastefy

import "errors"

type Folder struct {
	UserId   string   `json:"user_id"`
	Name     string   `json:"name"`
	Exists   bool     `json:"exists"`
	Id       string   `json:"id"`
	Pastes   []Paste  `json:"pastes,omitempty"`
	Children []Folder `json:"children,omitempty"`
}

func (apiClient PastefyApiClient) GetFolder(id string) (Folder, error) {
	folder := Folder{}
	_, err := apiClient.RequestMap("GET", "/folder/"+id, nil, &folder)
	if err != nil {
		return Folder{}, err
	}
	return folder, nil
}

func (apiClient PastefyApiClient) CreateFolder(folder Folder) (*Folder, error) {
	ret := struct {
		folder Folder `json:"folder"`
	}{}
	_, err := apiClient.RequestMap("POST", "/folder", folder, &ret)
	if err != nil {
		return nil, errors.New("error while fetching folder")
	}
	return &ret.folder, nil
}

func (apiClient PastefyApiClient) UpdateFolder(id string, folder Folder) (*Folder, error) {
	ret := struct {
		folder Folder `json:"folder"`
	}{}
	_, err := apiClient.RequestMap("PUT", "/folder/"+id, folder, &ret)

	if err != nil {
		return nil, errors.New("error while fetching folder")
	}
	return &ret.folder, nil
}

func (apiClient PastefyApiClient) SaveFolder(folder Folder) (*Folder, error) {
	if folder.Id != "" {
		return apiClient.UpdateFolder(folder.Id, folder)
	} else {
		return apiClient.CreateFolder(folder)
	}
}

func (apiClient PastefyApiClient) DeleteFolder(id string) error {
	_, err := apiClient.Request("DELETE", "/folder/"+id, nil)
	return err
}

func (apiClient PastefyApiClient) GetUserFolders() ([]Folder, error) {
	var folders []Folder
	_, err := apiClient.RequestMap("GET", "/user/folders", nil, &folders)
	if err != nil {
		return nil, err
	}
	return folders, nil
}
