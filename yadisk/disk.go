package yadisk

import (
	"context"
	"net/http"
)

// SystemFolders is the absolute addresses of Disk system folders.
type SystemFolders struct {
	Applications string `json:"applications"`
	Downloads    string `json:"downloads"`
}

// Disk represents the data about free and used space on the Disk.
// https://tech.yandex.com/disk/api/reference/response-objects-docpage/#disk
type Disk struct {
	// The cumulative size of the files in the Trash, in bytes.
	TrashSize uint `json:"trash_size"`

	// The total Disk space available to the user, in bytes.
	TotalSpace uint `json:"total_space"`

	// The cumulative size of the files already stored on the Disk, in bytes.
	UsedSpace uint `json:"used_space"`

	// Absolute addresses of Disk system folders.
	// Folder names depend on the user's interface language when the
	// personal Disk is created. For example, the "Downloads" folder
	// is created for an English-speaking user, the "Загрузки" folder
	// is created for a Russian-speaking user, and so on.
	//
	// The following folders are currently supported:
	//
	// * applications — folder for application files
	// * downloads — folder for files downloaded from
	// the internet (not from the user's device)
	SystemFolders SystemFolders `json:"system_folders"`
}

// DiskService handles communication with the data about a user's disk
// methods of the Yandex.Disk API.
//
// Yandex.Disk API docs: https://tech.yandex.com/disk/api/reference/capacity-docpage/
type DiskService service

// Get returns general information about
// a user's Disk: the available space, addresses of system folders, and so on.
// https://tech.yandex.com/disk/api/reference/capacity-docpage/
func (s *DiskService) Get(ctx context.Context) (*Disk, *http.Response, error) {
	url := "disk"
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	disk := new(Disk)
	resp, err := s.client.Do(ctx, req, disk)
	if err != nil {
		return nil, resp, err
	}

	return disk, resp, nil
}
