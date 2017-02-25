package unit

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/chibisov/go-yadisk/yadisk"
)

func TestDisk_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/disk/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(
			w,
			`
            {
                "trash_size": 4631577437,
                "total_space": 319975063552,
                "used_space": 26157681270,
                "system_folders":
                {
                    "applications": "disk:/Applications",
                    "downloads": "disk:/Downloads/"
                }
            }
            `,
		)
	})

	disk, response, err := client.Disk.Get(context.Background())

	if err != nil {
		t.Errorf("Disk.Get returned error %v, %+v", err, response)
	}

	// Check returned Disk instance
	diskWant := &yadisk.Disk{
		TrashSize:  4631577437,
		TotalSpace: 319975063552,
		UsedSpace:  26157681270,
		SystemFolders: yadisk.SystemFolders{
			Applications: "disk:/Applications",
			Downloads:    "disk:/Downloads/",
		},
	}
	if !reflect.DeepEqual(disk, diskWant) {
		t.Errorf("Disk.Get returned %+v, want %+v", disk, diskWant)
	}
}

func TestDisk_Get_with_http_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/disk/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "conflict", 409)
	})

	disk, _, err := client.Disk.Get(context.Background())

	if _, ok := err.(*yadisk.APIError); !ok {
		t.Errorf("Disk.Get should return APIError if HTTP error occured")
	}

	// Check returned Disk instance
	if disk != nil {
		t.Errorf("Disk.Get should return disk as nil if HTTP error occured")
	}
}
