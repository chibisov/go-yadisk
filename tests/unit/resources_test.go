package unit

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/chibisov/go-yadisk/yadisk"
)

func TestResources_Get_file(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/resources/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(
			w,
			`
            {
                "sha256": "d69e72661e26d9f1d44ab12e59c6cebfde48c125299db7768c913cfb5e42dffd",
                "name": "Горы.jpg",
                "created": "2017-02-26T09:04:44+00:00",
                "revision": 1488099884247936,
                "resource_id": "471259284:a1b980b1e83355ec15e8d0cc7fdaa9a34020b252de9c38c9ec4279af50823a13",
                "modified": "2017-02-26T09:24:44+00:00",
  "preview": "https://downloader.disk.yandex.ru/preview/2df91f854d7b2aa36e69d185c1ce9f19f28ce1051cf3c1efcfa08302d78c797f/inf/P1Uu3IBvn2J_059WhlC9XydbkJ2sC_KEk2e_VtbYAOfuuRRO_ACnqYduWLfxvvz5VgHzs7xBXDvk4Qqe7B86gQ%3D%3D?uid=471259284&filename=%D0%93%D0%BE%D1%80%D1%8B.jpg&disposition=inline&hash=              limit=0&content_type=image%2Fjpeg&tknv=v2&size=S&crop=0",
                "media_type": "image",
                "path": "disk:/Горы.jpg",
                "md5": "1392851f0668017168ee4b5a59d66e7b",
                "type": "file",
                "mime_type": "image/jpeg",
                "size": 1762478
            }
            `,
		)
	})

	resource, response, err := client.Resources.Get(context.Background(), "/", nil)

	if err != nil {
		t.Errorf("Resources.Get returned error %v, %+v", err, response)
	}

	// Checking resource properties
	if resource.PublicKey != nil {
		t.Errorf("Returned resource PublicKey is %v, want nil", resource.PublicKey)
	}
	if resource.Embedded != nil {
		t.Errorf("Returned resource Embedded is %v, want nil", resource.Embedded)
	}
	if got, want := resource.Name, "Горы.jpg"; got != want {
		t.Errorf("Returned resource Name is %v, want %v", got, want)
	}
	wantCreated := time.Date(
		2017,          // year
		time.February, // month
		26,            // day
		9,             // hour
		4,             // min
		44,            // sec
		0,             // nsec
		time.FixedZone("+0000", 0), // loc
	)
	if !resource.Created.Equal(wantCreated) {
		t.Errorf("Returned resource Created is %+v, want %+v", resource.Created, wantCreated)
	}
	if resource.CustomProperties != nil {
		t.Errorf("Returned resource CustomProperties is %v, want nil", resource.CustomProperties)
	}
	if resource.PublicURL != nil {
		t.Errorf("Returned resource PublicURL is %v, want nil", resource.PublicURL)
	}
	if resource.OriginPath != nil {
		t.Errorf("Returned resource OriginPath is %v, want nil", resource.OriginPath)
	}
	wantModified := time.Date(
		2017,          // year
		time.February, // month
		26,            // day
		9,             // hour
		24,            // min
		44,            // sec
		0,             // nsec
		time.FixedZone("+0000", 0), // loc
	)
	if !resource.Modified.Equal(wantModified) {
		t.Errorf("Returned resource Modified is %+v, want %+v", resource.Modified, wantModified)
	}
	if got, want := resource.Path, "disk:/Горы.jpg"; got != want {
		t.Errorf("Returned resource Path is %v, want %v", got, want)
	}
	if got, want := resource.MD5, "1392851f0668017168ee4b5a59d66e7b"; got != want {
		t.Errorf("Returned resource MD5 is %v, want %v", got, want)
	}
	if got, want := resource.SHA256, "d69e72661e26d9f1d44ab12e59c6cebfde48c125299db7768c913cfb5e42dffd"; got != want {
		t.Errorf("Returned resource SHA256 is %v, want %v", got, want)
	}
	if got, want := resource.Revision, uint(1488099884247936); got != want {
		t.Errorf("Returned resource Revision is %v, want %v", got, want)
	}
	if got, want := resource.ResourceID, "471259284:a1b980b1e83355ec15e8d0cc7fdaa9a34020b252de9c38c9ec4279af50823a13"; got != want {
		t.Errorf("Returned resource ResourceID is %v, want %v", got, want)
	}
	if got, want := resource.Type, "file"; got != want {
		t.Errorf("Returned resource Type is %v, want %v", got, want)
	}
	if got, want := resource.MediaType, "image"; got != want {
		t.Errorf("Returned resource MediaType is %v, want %v", got, want)
	}
	if got, want := resource.MimeType, "image/jpeg"; got != want {
		t.Errorf("Returned resource MimeType is %v, want %v", got, want)
	}
	if got, want := resource.Size, uint(1762478); got != want {
		t.Errorf("Returned resource Size is %v, want %v", got, want)
	}
}

func _TestResources_Get_with_http_error(t *testing.T) {
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
