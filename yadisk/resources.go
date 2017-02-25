package yadisk

import (
	"context"
	"net/http"
	"time"
)

// Resource is a description or metainformation about a file or folder.
// https://tech.yandex.com/disk/api/reference/response-objects-docpage/#resource
type Resource struct {
	// Key of a published resource.
	// It is included in the response only if
	// the specified file or folder is published.
	PublicKey string `json:"public_key"`

	// Resources directly contained in the folder.
	// It is included in the response only
	// when requesting metainformation about a folder.
	Embedded ResourceList `json:"_embedded"`

	// Resource name
	Name string `json:"name"`

	// The date and time of the resource was created.
	// JSON data is in ISO 8601 format.
	Created time.Time `json:"created"`

	// Object with the user defined attributes
	CustomProperties map[string]string `json:"custom_properties"`

	// Link to a published resource.
	// It is included in the response only if
	// the specified file or folder is published.
	PublicURL string `json:"public_url"`

	// Path to the resource before it was moved to the Trash.
	// Included in the response only for a request
	// for metainformation about a resource in the Trash.
	OriginPath string `json:"origin_path"`

	// The date and time the resource was modified.
	// JSON data is in ISO 8601 format.
	Modified time.Time `json:"modified"`

	// Full path to the resource on Disk.
	// In metainformation for a published folder, paths are relative
	// to the folder itself. For published files, the value
	// of the key is always "/". For a resource located in the Trash,
	// this attribute may have the unique ID appended to it
	// (for example, trash:/foo_1408546879).
	// Using this ID, the resource can be differentiated from other
	// deleted resources with the same name.
	Path string `json:"path"`

	// MD5 hash of the file.
	MD5 string `json:"md5"`

	// Resource type:
	// * "dir" - folder
	// * "file" - file
	Type string `json:"type"`

	// The MIME type of the file.
	MimeType string `json:"mime_type"`

	// File size.
	Size uint `json:"size"`
}

// ResourceList is a list of resources contained in the folder.
// Contains Resource objects and list properties.
// https://tech.yandex.com/disk/api/reference/response-objects-docpage/#resourcelist
type ResourceList struct {
	// The field used for sorting the list.
	Sort string `json:"sort"`

	// The key of a published folder that contains resources from this list.
	// It is included in the response only when requesting
	// metainformation about a public folder.
	PublicKey string `json:"public_key"`

	// Array of resources contained in the folder.
	Items []Resource `json:"items"`

	// The path to the folder whose contents are described
	// in this ResourceList object.
	// For a public folder, the value of the attribute is always "/".
	Path string `json:"path"`

	// The maximum number of items in the items array; set in the request.
	Limit uint `json:"limit"`

	// How much to offset the beginning of the
	// list from the first resource in the folder.
	Offset uint `json:"offset"`

	// The total number of resources in the folder.
	Total uint `json:"total"`
}

// ResourcesService handles communication with the metainformation
// about files and folders. Metainformation includes the properties of
// files and folders, and the properties and contents of subfolders.
type ResourcesService service

// ResourcesOptions specifies the optional parameters to the
// ResourcesService.Get method
type ResourcesOptions struct {
	// The attribute used for sorting the list of resources in the folder.
	// The names of the following keys for the Resource object can
	// be used as the value:
	//
	// * "name" (resource name)
	// * "path" (path to the resource on Disk)
	// * "created"" (date the resource was created)
	// * "modified" (date the resource was modified)
	// * "size" (file size)
	//
	// To sort in reverse order, add a hyphen to the value of the parameter,
	// for example: sort="-name".
	Sort string `url:"sort"`

	// The number of resources in the folder that should be described
	// in the response (for example, for paginated output).
	// The default value is 20.
	Limit uint `url:"limit"`

	// The number of resources from the top of the list that
	// should be skipped in the response (for example, for paginated output).
	//
	// Let's say the /foo folder contains three files.
	// If we request metainformation about the folder with the offset=1
	// parameter and default sorting, the Yandex.Disk API returns
	// only the descriptions of the second and third files.
	Offset uint `url:"offset"`

	// List of JSON keys that should be included in the response.
	// Keys that are not included in this list will be discarded when
	// forming the response. If the parameter is omitted, the response is
	// returned in full, without discarding anything.
	//
	// Embedded keys should be separated by dots.
	// For example: ["name", "_embedded.items.path"].
	Fields []string `url:"fields"`

	// The required size of the reduced image (file preview),
	// which the API returns a reference to in the preview key.
	//
	// You can define the exact size of the preview, or the length of one
	// of the sides. The resulting image can be cropped to a square
	// using the PreviewCrop parameter.
	//
	// Possible values:
	//
	// * Predefined length of the longest side.
	//   The image is reduced to the specified length of the longest side,
	//   and the proportions of the source image are preserved.
	//   For example, for the size "S" and an image sized 120×200,
	//   a preview sized 90×150 will be generated, while for an image
	//   sized 300×100, the preview will be 150×50.
	//   Supported values:
	//     * "S" — 150 pixels
	//     * "M" — 300 pixels
	//     * "L" — 500 pixels
	//     * "XL" — 800 pixels
	//     * "XXL" — 1024 pixels
	//     * "XXXL" — 1280 pixels
	// * The exact width (for example, "120" or "120x") or the exact
	//   height (for example, "x145").
	//   The image is reduced to the specified width or height,
	//   and the proportions of the source image are preserved.
	//   If the PreviewCrop parameter is passed,
	//   a square with the set side length is cut out of the center
	//   of the reduced image.
	// * The exact size (in the format <width>x<height>, such as "120x240").
	//   The image is reduced to the smallest of the specified dimensions,
	//   and the proportions of the source image are preserved.
	//   If the PreviewCrop parameter is passed,
	//   a section is cut from the center of the source image with
	//   the maximum size in the set proportions of width
	//   to height (in the example, this is 1/2).
	//   Then the cropped section is scaled to the specified dimensions.
	PreviewSize string `url:"preview_size"`

	// This parameter cuts the preview to the size specified
	// in the PreviewSize parameter. When set to false (default setting),
	// the parameter is ignored.
	// When set to true, the preview is cropped as follows:
	// * If only the width or height is passed, the image is reduced
	//   to this size with the proportions preserved.
	//   Then a square with the specified length of side is cut out
	//   of the center of the reduced image.
	// * If the exact size is passed (for example, "120x240""),
	//   a section is cut from the center of the source image with the
	//   maximum size in the set proportions of width to height.
	//   Then the cropped section is scaled to the specified dimensions.
	PreviewCrop bool `url:"preview_crop"`
}

// Get retunes metainformation for the path. The path to the desired resource
// is relative to the Disk root directory.
// The path to a resource in the Trash should be relative to the Trash root directory.
//
// Yandex.Disk API docs: https://tech.yandex.com/disk/api/reference/meta-docpage/
func (s *ResourcesService) Get(
	ctx context.Context,
	path string,
	opt *ResourcesOptions,
) ([]*ResourceList, *http.Response, error) {
	return nil, nil, nil
}
