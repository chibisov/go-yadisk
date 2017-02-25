package integration

import (
	"context"
	"testing"

	"github.com/chibisov/go-yadisk/yadisk"
)

func TestDisk_Get(t *testing.T) {
	client := yadisk.NewClient(config(t).AccessToken)
	disk, resp, err := client.Disk.Get(context.Background())
	if err != nil {
		t.Errorf("Disk.Get returned error %v, %+v", err, resp)
	}

	if disk.TrashSize < 0 {
		t.Error("Disk TrashSize must be more than or equal to zero")
	}
	if disk.TotalSpace <= 0 {
		t.Error("Disk TotalSpace must be more than zero")
	}
	if disk.UsedSpace <= 0 {
		t.Error("Disk UsedSpace must be more than zero")
	}
	if got, want := disk.SystemFolders.Applications, "disk:/Приложения"; got != want {
		t.Errorf("disk SystemFolders.Applications = %v, want %v", got, want)
	}
	if got, want := disk.SystemFolders.Downloads, "disk:/Загрузки/"; got != want {
		t.Errorf("disk SystemFolders.Downloads = %v, want %v", got, want)
	}
}
