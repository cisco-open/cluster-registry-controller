// Code generated by vfsgen; DO NOT EDIT.

package cluster_registry

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Chart statically implements the virtual filesystem provided to vfsgen.
var Chart = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Time{},
		},
		"/.helmignore": &vfsgen۰CompressedFileInfo{
			name:             ".helmignore",
			modTime:          time.Time{},
			uncompressedSize: 349,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x4c\x8e\x41\x6e\xe3\x30\x0c\x45\xf7\x3c\xc5\x1f\x78\x33\x63\x0c\xe4\x43\x24\xb3\x98\x55\x0b\xa4\xc8\xb6\x90\x6d\x46\x62\x22\x8b\x82\x44\x27\x6d\x17\x3d\x7b\x91\x04\x41\xbb\x79\x20\x3f\xc8\x8f\xd7\xe1\xd9\x9b\x71\xcd\x0d\xa6\x90\x90\xb5\x32\x2e\x91\x33\xc6\x55\xd2\x2c\x39\xa0\xf8\xe9\xe4\x03\x37\x47\x1d\x5e\xa2\x34\xb4\xb5\x14\xad\xd6\xd0\x22\xa7\x84\x90\x74\xc4\xe2\x6d\x8a\x92\xc3\x5f\x54\x4e\xde\xe4\xcc\x28\xde\xe2\x8f\xdc\xe7\x99\x3a\x64\x0e\xde\x44\x33\x7e\x97\xca\x07\x79\xe3\x19\x17\xb1\x88\x5f\x7f\x1c\x9e\x72\x7a\x87\xe6\xdb\xe7\x55\x09\x85\x2b\x92\x64\x76\xe4\xb6\xbb\xd7\x9d\x69\x65\xea\xb0\xd1\x65\xd1\x8c\xfd\x66\x87\x59\x6a\x23\x17\xc4\x86\x1b\xef\xfa\xe4\xc6\x8f\x3a\xdc\xf8\x08\x62\x18\xae\x78\xac\xed\x9c\x87\xef\xa2\xd1\x4f\xa7\xb5\xe0\x20\x89\x1b\xf5\xae\x5d\x0a\xf5\x6e\xf4\x27\xea\x9d\x2d\xd7\x59\xab\x04\xea\x3f\xa9\xc3\xde\x57\xd1\xb5\xe1\xff\xf6\x5f\x23\x57\xaa\x1e\x79\x32\x72\x32\xb3\x1f\xee\xe7\x55\x8f\xe4\xce\x6d\xd2\x99\x07\xfa\x0a\x00\x00\xff\xff\x16\xec\x32\x27\x5d\x01\x00\x00"),
		},
		"/Chart.yaml": &vfsgen۰CompressedFileInfo{
			name:             "Chart.yaml",
			modTime:          time.Time{},
			uncompressedSize: 257,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x5c\x8f\x3d\x6e\x84\x40\x0c\x85\x7b\x9f\xc2\x17\x58\xd8\xed\xa2\xa9\xa2\x50\xa6\xa7\x37\x83\x01\x2b\xcc\x8f\xec\x19\x24\x72\xfa\x08\x10\x52\x92\xf6\x7d\xef\xd3\xb3\x29\x4b\xcf\x6a\x92\xa2\xc3\xed\x05\x91\x02\x3b\xf4\x6b\xb5\xc2\xfa\x50\x9e\xc5\x8a\xee\x30\xb2\x79\x95\x5c\xce\x5a\x77\x51\xbc\x29\x4e\x49\xf1\xf3\xcd\x60\x49\x87\xbc\x94\x92\xcd\xb5\xed\x40\xf1\x9b\xc4\xaf\xa9\x8e\x8d\x4f\x01\x20\x90\xc4\x42\x12\x59\xcd\xc1\x03\x39\x90\xac\x0e\x25\x4e\xe9\xfd\x7f\x17\xf1\x3a\xe4\xe3\xcc\xb1\x3b\x00\x58\xaa\xea\xf9\x74\xef\x8d\x59\xca\x52\x87\x43\xf9\x3d\xd7\x0e\xe4\xbf\x76\xd2\xd1\x00\xb6\xfb\xb9\x67\xf3\x6c\x5e\x40\x39\xf7\x7f\x93\x9f\x00\x00\x00\xff\xff\x3b\x93\x34\x08\x01\x01\x00\x00"),
		},
		"/crds": &vfsgen۰DirInfo{
			name:    "crds",
			modTime: time.Time{},
		},
		"/crds/crds.yaml": &vfsgen۰CompressedFileInfo{
			name:             "crds.yaml",
			modTime:          time.Time{},
			uncompressedSize: 15352,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xec\x1a\x6d\x6f\xdb\xb8\xf9\xbb\x7f\xc5\x83\x6c\x40\xdb\x21\x72\x5a\xdc\x97\xcd\xc0\xa1\x0b\x92\x1e\x96\xb5\xd7\x06\x49\xae\x1f\x56\x74\x03\x2d\x3d\xb6\xb9\x50\xa4\x4a\x52\x6e\xd5\xc3\xfd\xf7\xe1\x21\x29\x59\x92\x25\x9b\x76\x5f\x70\xdb\xe2\x2f\x89\xf9\xf2\xbc\xbf\x9b\x93\x24\x49\x26\xac\xe0\x6f\x51\x1b\xae\xe4\x0c\x58\xc1\xf1\x93\x45\x49\xdf\xcc\xf4\xfe\xcf\x66\xca\xd5\xd9\xfa\xd9\x1c\x2d\x7b\x36\xb9\xe7\x32\x9b\xc1\x45\x69\xac\xca\x6f\xd0\xa8\x52\xa7\x78\x89\x0b\x2e\xb9\xe5\x4a\x4e\x72\xb4\x2c\x63\x96\xcd\x26\x00\x4c\x4a\x65\x19\x2d\x1b\xfa\x0a\x90\x2a\x69\xb5\x12\x02\x75\xb2\x44\x39\xbd\x2f\xe7\x38\x2f\xb9\xc8\x50\x3b\x0c\x35\xfe\xf5\xd3\xe9\x0f\xd3\xa7\x13\x80\x54\xa3\xbb\x7e\xc7\x73\x34\x96\xe5\xc5\x0c\x64\x29\xc4\x04\x40\xb2\x1c\x67\x90\x8a\xd2\x58\xd4\x66\x1a\xfe\xd1\xb8\xe4\xc6\xea\xca\xd1\x9c\x72\x93\xaa\x69\xaa\xf2\x89\x29\x30\x75\xf4\x64\x99\x23\x92\x89\x6b\xcd\xa5\x45\x7d\xa1\x44\x99\x7b\xe2\x12\xf8\xfb\xed\x9b\xd7\xd7\xcc\xae\x66\x30\x35\x96\xd9\xd2\xb8\x3f\xe8\x08\xf7\xe8\x6e\xdd\xb2\x5b\xb0\x55\x81\x33\x30\x56\x73\xb9\x1c\xb9\x4d\x47\x5a\x97\xef\xea\xaf\x11\x57\x53\x25\x3d\xa5\xe6\xdd\xf3\xc7\x7f\x75\x90\x7e\xfc\xf1\xe4\x22\x70\x7b\x5b\xc9\x14\xb3\x93\x27\xef\xc3\xf1\x36\x85\x6e\x2b\x16\x4d\x10\x78\xeb\xfe\xdb\xd6\x4a\xa1\xb9\xd2\xdc\x56\x33\x78\x16\x0b\xb0\xd0\x6a\xcd\x33\xd4\x2d\x88\xd7\xed\xa5\x23\x40\x66\xa4\x50\x3e\x2f\x6d\x97\xd0\xcb\xfe\xf2\x11\xa0\x85\x4a\x99\xe0\xb6\x9a\x92\xd9\x74\xa0\xdf\x6c\x16\x8e\x80\x9b\xa3\x31\x6c\xb9\x6d\x38\xf0\x73\x6b\xe3\x08\xb8\xb1\x56\x31\x80\xbf\x92\xe9\x21\xd8\x97\x5a\x95\x45\xe3\x5e\x23\x5e\xe5\x81\x07\xcf\x0e\x51\xc1\x5f\x70\x2b\x82\x1b\xfb\xb2\xbd\xfa\x8a\x1b\xeb\x91\x8b\x52\x33\xb1\xf1\x5e\xb7\x68\xb8\x5c\x96\x82\xe9\x66\x79\x42\x54\xa2\x41\xbd\xc6\x5f\xe4\xbd\x54\x1f\xe5\x4f\x1c\x45\x66\x66\xb0\x60\xc2\x10\x1b\x26\x55\x44\xf5\x06\xa9\x29\xe7\x3a\x44\xa4\x40\x96\x97\xdc\x0c\x7e\xfd\x6d\x02\xb0\x66\x82\x67\x2e\x9e\xf8\x4d\x55\xa0\x3c\xbf\xbe\x7a\xfb\xc3\x6d\xba\xc2\x9c\xf9\x45\x80\x0c\x4d\xaa\x79\xe1\xce\xd5\xc0\x81\x1b\xb0\x2b\x04\x7f\x12\x16\x4a\xbb\xaf\x35\x07\x70\x7e\x7d\x15\x6e\x17\x5a\x15\xa8\x2d\xaf\x29\xa0\x4f\x2b\xb6\x36\x6b\x3d\x3c\x8f\x88\x10\x7f\x06\x32\x8a\xa6\xe8\x11\x06\x17\xc5\x0c\x8c\x47\xad\x16\x60\x57\xdc\x80\x46\x27\x1d\xe9\xe3\x6b\x0b\x2c\xd0\x11\x26\x41\xcd\xff\x8d\xa9\x9d\xc2\x2d\x49\x50\x1b\x30\x2b\x55\x8a\x8c\x42\xf0\x1a\xb5\x05\x8d\xa9\x5a\x4a\xfe\xb9\x81\x6c\xc0\x2a\x87\x52\x30\x8b\x41\x53\xf5\xc7\x05\x4b\xc9\x04\x89\xb0\xc4\x53\x60\x32\x83\x9c\x55\xa0\x91\x70\x40\x29\x5b\xd0\xdc\x11\x33\x85\x9f\x95\x46\xe0\x72\xa1\x66\xb0\xb2\xb6\x30\xb3\xb3\xb3\x25\xb7\x75\x36\x49\x55\x9e\x97\x92\xdb\xea\xcc\xe5\x04\xf2\x64\xa5\xcd\x59\x86\x6b\x14\x67\x86\x2f\x13\xa6\xd3\x15\xb7\x98\xda\x52\xe3\x19\x2b\x78\xe2\x08\x97\xce\xfc\xa7\x79\xf6\x87\x46\xd1\x8f\x5a\x94\xf6\xcc\xd8\x7f\x9c\x69\x8e\xca\x9d\x4c\x94\xb4\xcb\xc2\x35\x4f\xff\x46\xbc\xb4\x44\x52\xb9\x79\x71\x7b\x07\x35\x52\xa7\x82\xae\xcc\x9d\xb4\x37\xd7\xcc\x46\xf0\x24\x28\x2e\x17\xa8\xbd\xe2\x16\x5a\xe5\x0e\x22\xca\xac\x50\x5c\xda\x60\x49\x1c\x65\x57\xe8\xa6\x9c\xe7\xdc\x92\xa6\x3f\x94\x68\x2c\xe9\x67\x0a\x17\x2e\xa7\xc2\x1c\xa1\x2c\x32\x66\x31\x9b\xc2\x95\x84\x0b\x96\xa3\xb8\x60\x06\xbf\xb9\xd8\x49\xc2\x26\x21\x91\xee\x17\x7c\xbb\x14\xe8\x1e\xf4\xd2\x6a\x96\xeb\xfc\x3c\xa8\xa1\xe0\x81\xb7\x05\xa6\x1d\xcf\xc8\xd0\x70\x4d\xd6\x4b\x49\x9a\x6c\xbe\x1d\x7c\xc6\x7d\xd1\xf9\x63\x69\x57\x57\x24\xa2\xce\x6a\x0f\xef\x79\x38\x04\x2b\x25\x32\xe3\x44\xaa\x73\xe7\x6c\x60\x57\xcc\x86\xc3\x73\x34\xb0\x52\x1f\x81\x0d\x69\xd0\x95\x3c\x4c\xc2\x12\x2d\x95\x32\x19\xc9\x91\x09\xe7\x68\x2c\x4d\xd1\x98\x76\x10\x99\xf6\xae\x8e\x11\xef\x04\x86\xa9\x46\x7b\x83\x8b\xed\xad\x1e\x17\x2f\x3e\x94\x7c\xcd\x04\x4a\xeb\x22\x07\x69\x6f\xfa\x9a\xc2\x76\xc1\x52\xcc\xe8\x3f\xf8\xc8\xed\xca\x65\x1b\xb0\x6c\x69\x06\x00\xee\xa2\xa4\xc9\x31\x83\x3b\x23\x86\xd1\xbf\xec\x88\x39\x12\xc2\xa0\x41\xed\xd9\x0a\x12\xbf\xba\xdc\x69\x00\xbf\x5c\x5d\xfa\x68\x8b\x40\x85\x6a\x62\x2a\x63\x31\xdf\x10\x3c\x89\xa2\x93\xbc\x97\xec\xb4\x8d\x2a\xd9\x50\xb0\xd7\x35\x7c\x02\xdb\xe7\x1c\xbe\xbe\x68\xbb\x87\x9a\xbb\xd4\x79\x94\x7f\x6c\x0a\x8d\x9d\x02\xba\x68\x8e\xb9\xc2\x9e\x71\x19\x3c\x93\x2f\x16\xa8\xc9\xe4\x1a\x40\x81\x0f\x34\x94\x3a\xb7\xb4\xe8\xe2\xe2\x88\x1b\x70\x8b\xf9\x80\xdd\x0d\x49\xa1\xa1\x67\x43\xce\x86\x80\xb6\xff\x52\xfa\x66\x03\xb6\x34\x42\xc2\x3e\x0f\x10\xcc\xd8\xbf\x21\xd3\x76\x8e\xcc\x52\xab\x32\x6c\xca\x1d\x92\x5f\xf5\xef\xd4\x15\x06\x01\x03\x4b\x0b\x5e\x2a\x35\x03\x23\xde\xf1\x91\x99\x26\x13\x0c\x1e\xf1\x5c\xcf\x80\x8e\x24\x04\x77\x72\x84\x93\x11\x51\x77\x9a\x49\xc3\xeb\x6e\x2c\x92\xc5\xee\xa5\x21\x1e\x71\x2f\x8b\xe9\x8a\xc9\x25\x66\x3e\x71\x2a\x89\xc1\x96\x5c\x14\x95\xca\xae\x86\x14\xf6\xd5\x38\x0f\xb5\x74\x04\xbb\xa1\xbc\xf6\xb5\xc4\xaa\xcc\x99\x4c\x34\xb2\x8c\xcd\x05\xd6\x50\x80\xcb\x8c\xa7\xcc\xd5\x14\x19\x5a\xc6\x85\x19\xe1\x99\xcd\x55\x69\x37\xb2\x0a\x1c\x7b\x49\x4c\x8f\xe1\x43\x23\x33\xdd\xf2\x73\x84\x8d\x1b\x77\xd0\x73\xf1\x78\xae\x39\x2e\x9e\x84\xcb\x9b\xaa\xb7\x56\xd8\x23\xe3\xc8\x1b\xe1\xe1\xcb\x89\xde\x0e\x7e\x23\x44\x87\xf8\x17\xcc\x2b\x20\x0e\xb1\xbb\xa1\x76\x0a\x6f\xa4\x0b\x84\x77\xba\xc4\xd3\x11\xa2\x7f\xa2\xde\xe2\x14\x42\xc7\x71\x14\xd5\x6e\x7b\x3f\xcd\x77\x55\xd1\x38\x04\x5d\x69\xe8\x0d\x1d\xc7\x86\xee\xc3\x89\x18\x4a\x3a\x75\xea\x69\xcd\x0c\xba\x1b\xcd\xbc\xe2\xa0\xf4\xca\xb4\x66\x55\x67\xa7\xdd\xb1\xcf\x26\x91\x54\x53\x86\xbd\xd6\xea\x53\x15\x9a\xa0\x2d\xad\x8f\xe4\x81\x1d\x62\x18\xa3\x8f\x50\x09\xb4\xdf\x1e\x91\x40\x96\xa1\x1e\x16\xc1\x5c\x29\x81\xac\x1b\xf2\xea\x79\xc4\xec\x80\x4a\xd0\x8f\x2e\x66\x93\x83\x0d\x64\x39\xc4\xfb\x0e\xfe\xa3\x4b\xb1\x6d\x39\xd0\xe7\xb3\x92\xf8\xdd\xd0\x8d\x1a\xed\x48\x34\x1f\x45\x54\xcf\xb3\xa2\x2f\xb8\x62\x2b\xfa\xf4\x50\xa4\x18\x3d\xbc\xde\x9e\x20\xec\x38\x3f\x20\x82\xde\xd2\x66\xda\xfa\x8c\x89\x62\xc5\x9e\x6d\xd6\xc2\x40\xd4\xcf\x8e\x5a\xdb\xd4\x78\x50\x4d\x39\x03\xab\x4b\x0c\x03\x16\xa5\x49\xa2\x7e\x65\x13\xb0\xa9\xbf\x29\xac\xef\x30\x3a\x23\xa2\x93\x93\xce\x0c\xc8\x7d\x6d\xd5\x9b\xf0\xee\xfd\xc4\x43\xc5\xac\xf1\x50\x5a\xfc\xef\x1d\x52\x37\x93\x8a\x4a\xa6\xba\x14\x18\x3b\xad\x3e\x76\x08\x57\xf3\x7b\x5b\xc9\xf4\xa6\x14\xd8\x9b\xc6\xf5\xb7\xb7\xc6\x72\x5b\xf4\xf6\xe6\x73\xfd\xfd\xdf\xc7\xa0\xae\xcf\xd6\xc8\xc4\xae\x99\xe0\x10\xf5\x40\xe4\x3f\xcc\xee\x1e\x66\x77\x0f\xb3\xbb\xaf\x36\xbb\x1b\x2b\x56\x5c\x34\x0b\x5e\xf2\xb2\xa7\xd5\x7d\x45\x8e\x8f\x84\x87\xd6\x38\xf7\x03\x58\xf6\x5e\x1a\xc9\xb2\x51\x15\xe0\x40\xc1\xe1\xe2\x67\x64\x85\xb9\x7b\xc8\x90\x33\x9b\xae\x86\xeb\xa4\x1d\x25\xd4\xfe\xe1\x9d\x8b\x71\xfd\x14\x38\xfc\xd9\x89\x27\x1e\x5b\x8b\x9f\xf3\x38\xc4\x81\xca\xd6\xaf\xb7\xb1\x58\xa2\x46\x90\x91\xfd\xce\x20\x07\x2f\x3e\x51\xd0\x30\x71\x1c\x44\x08\x70\x68\x10\xdd\xd2\x0f\x18\x14\x98\x5a\xa5\xeb\x16\x2f\x47\x69\x81\x9b\x08\x98\x40\x61\xb2\xbe\xed\x26\xd8\xcd\x94\xcc\x87\xfd\x53\x60\x70\x8f\x95\xcf\x10\x4c\x46\x81\x24\x3d\xb0\x06\xa0\x46\x97\x7f\xfc\xc8\x14\x2b\x07\x28\xa4\x94\x08\x68\xc5\x41\x5a\x05\xc2\x10\x77\xb0\x27\x4f\xa2\xac\x19\x44\xcd\x51\xb8\x05\x47\xbf\x1b\x1e\x04\x11\x45\x42\xa6\x02\xa1\x10\x1c\x5d\x84\x8f\xbc\x73\x90\x41\xb6\xa5\x7c\x14\xbb\x8d\x8a\x36\xf9\xcd\x2b\xfa\x91\xf1\x0a\x23\xdb\x5d\xf1\x22\x9a\x61\xab\x9c\x25\xb9\x1f\x14\xea\x82\xe1\x2d\x55\x6c\x0d\x2a\x03\x4c\x23\x5c\xc9\xd3\x68\x98\xaf\x95\xbd\x92\xa7\xf0\xe2\x13\xa7\x64\x49\x76\x73\xa9\xd0\xbc\x56\xd6\xad\x7c\x33\xc1\x7a\xf2\x8f\x12\xab\xbf\xea\xaa\x0f\xe9\x9b\x50\x92\x47\xbb\x0e\x31\xd3\x68\xf6\xaf\xfc\xe8\xa7\x51\x15\x37\x54\x19\x28\x5d\xcb\xc5\x55\x93\x0e\x66\xbc\x59\x3a\x92\xf2\xd2\xb8\x82\x43\x2a\x99\x60\x5e\xd8\x6a\x3a\x80\x2b\x1a\x66\x50\x8f\xd2\x1d\xed\xb4\xc9\x6b\xa1\x8d\x86\x3a\x47\x08\xa4\xdd\x51\x8d\xe5\x21\xf8\x2a\x59\xb0\x14\x33\xc8\x4a\x27\x54\x16\x0d\xd1\x58\xcd\x2c\x2e\x79\x0a\x39\xea\x25\x42\x41\x91\x3a\x56\x1b\xd1\x41\x3a\x44\x2d\x66\xa9\xbe\x9e\xc1\x3f\x1f\x3f\x7e\x77\x9e\xfc\x83\x25\x9f\x9f\x26\x7f\x79\xff\x2e\x69\xfe\xff\xd7\xf4\xfd\x9f\x9e\x3c\x6f\xed\x3d\x79\xfe\xc7\x78\x67\x3b\xd4\xa4\x77\x8f\x61\x86\x06\xc3\x63\xa3\xc2\xed\x4f\x42\x61\x23\xea\x5c\x6d\x5d\xd1\xe9\x39\x2a\xe3\xc6\xf3\x16\x09\x34\x06\x1c\x65\x49\x94\xf6\xbb\x15\x46\x51\x69\xed\x00\xb3\x70\x8e\x19\x51\x5d\xc9\xea\xcd\x62\xff\xb1\x24\xa0\xa6\xae\x72\x89\x3a\xfa\x7c\xa4\x05\x7f\x4a\xee\xcb\x39\x6a\x89\x16\x4d\xc2\xa5\x4d\x94\x4e\xfc\xd5\xd6\xcc\xe9\xcb\x4c\x79\xbf\x11\x27\x5e\x66\xdf\xcb\xc0\x5c\x11\xf2\xa5\x85\x77\xaf\x62\xf4\x85\x4d\x53\xee\xb9\x1e\xd9\xaf\x7d\x28\x51\x57\xa0\xd6\xa8\xf7\x06\xd4\x90\xe3\x9b\x6e\x9d\x02\xb4\x9b\xa0\x94\xc2\x6d\xb8\x02\xf8\x95\xa3\x3e\x4c\x15\xba\x05\xf1\x64\x5f\x86\x42\x38\x7f\x7d\x49\x6d\xf0\xb9\xf4\x29\xa0\x4f\xb7\x83\x48\x59\x45\x88\x20\xeb\xbd\x49\xf5\xdc\x0d\xdf\xc6\x00\x49\x15\x07\xe7\xc0\x3e\xe6\xa0\x2e\xa0\xa3\xaa\xfe\xf5\xa0\x2a\x6e\x9c\x84\xbb\x5c\xc4\x47\xf3\xdc\x4f\x31\xbc\xba\x36\x2b\x2d\x91\x7f\xb3\x5e\xa5\x27\xf8\x6e\x9b\xd2\x6a\x41\xa2\x92\x59\x44\x9b\xd2\x6d\x41\xa2\xa0\x3e\xb4\x29\x0f\x6d\xca\x43\x9b\xf2\xd0\xa6\xfc\x3f\xb5\x29\x0f\x7d\xc4\x91\xbc\xb5\x8a\x9c\xdf\xc9\x88\x72\xbb\x7e\x08\x35\x98\xcb\xaf\x39\x2b\xc8\xc3\x7f\xa5\x14\xe9\x8c\xfd\x37\x28\x18\xd7\x51\x5e\x7e\xee\x7e\xe5\x13\xd8\xb9\xcd\xa5\x73\x9c\x36\x22\xc2\xc1\x0d\x60\xf3\xb6\x74\x12\x17\x8e\x25\xa0\xf0\xa5\x40\x5d\x3d\xb6\x2a\x9f\x53\xf8\xb8\x52\xc6\x67\xe4\x05\x47\x91\x45\x00\xe5\x06\x4e\xee\xb1\x3a\x39\xdd\x8a\x4b\x27\x57\xf2\xc4\x97\x08\x7d\xaf\x8f\x00\xdb\x54\x1c\x4a\x8a\x0a\x4e\xdc\xed\x93\x2f\x2b\xa7\xa2\xad\xf3\x2b\x36\x16\xcd\x4b\xd5\x2f\x6d\x2e\x22\xcd\x33\x86\x26\xcf\xd8\xcb\xdd\x35\x51\x6c\x89\xb5\xeb\xe1\xf1\xc1\xae\xb5\xf7\x21\xf2\x51\xd2\xd8\xa9\xc8\x7d\x31\x33\x69\xa9\x70\x72\x24\x96\xdd\x4a\xc9\xcb\x9d\x3f\xbf\xec\x57\x04\x35\x92\x9a\x67\xbb\x34\xb5\xd7\xc0\x62\xd5\x5d\x30\x6d\xf0\x6d\xcc\x38\x63\xfc\x8d\xd5\xc0\xf8\x6e\xf5\xd5\x4c\x68\xfc\xe5\xdf\x11\xc0\xd6\xf1\x8c\x7e\x0d\x5b\x8c\x71\x5f\x53\xc9\xf4\x76\xc7\x9b\xcc\x38\xd1\xef\x24\xe6\xe0\xf7\x86\xc3\xef\xeb\xfb\xbf\x36\x77\x36\x37\x8f\x5a\x76\xe0\xdc\x7e\x7b\xfa\x3f\xf7\xa4\xea\x3f\x01\x00\x00\xff\xff\x3d\x1f\x6f\x46\xf8\x3b\x00\x00"),
		},
		"/examples": &vfsgen۰DirInfo{
			name:    "examples",
			modTime: time.Time{},
		},
		"/examples/test.yaml": &vfsgen۰CompressedFileInfo{
			name:             "test.yaml",
			modTime:          time.Time{},
			uncompressedSize: 153,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x3c\xcc\x4d\x8e\xc2\x30\x0c\x40\xe1\x7d\x4e\x91\x0b\xb8\x23\x0f\x69\xeb\x66\x0b\x1b\x2e\xc0\xde\x75\x2c\x88\xfa\xab\x38\x20\x71\x7b\x84\x40\x6c\x3f\x3d\x3d\x00\x70\xbc\xe7\x8b\x16\xcb\xdb\x1a\xbd\xcc\x77\xab\x5a\x8a\x5e\xb3\xd5\xf2\x6c\x26\xb2\x46\xb2\xc9\xd6\xc8\xb6\xfc\x3d\x90\xe7\xfd\xc6\xe8\xa6\xbc\xa6\xe8\x8f\x9f\xd8\x2d\x5a\x39\x71\xe5\xe8\xbc\x5f\x79\xd1\xdf\x06\xd0\xd9\xae\xf2\xf6\xaf\x9c\x4f\xd1\xa7\x16\x03\x92\xf6\xd0\x86\x5e\x21\x74\x34\x02\xa7\x51\x60\x90\xff\xa1\xa3\x03\x0d\x48\xc1\xbd\x02\x00\x00\xff\xff\x73\x74\xf2\x63\x99\x00\x00\x00"),
		},
		"/templates": &vfsgen۰DirInfo{
			name:    "templates",
			modTime: time.Time{},
		},
		"/templates/_helpers.tpl": &vfsgen۰CompressedFileInfo{
			name:             "_helpers.tpl",
			modTime:          time.Time{},
			uncompressedSize: 2035,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x54\x41\x6b\x1b\x3d\x10\xbd\xfb\x57\x0c\xcb\x17\xf8\x9a\xb2\xf2\xa1\xd0\x83\x21\x87\x92\xf6\x50\x5a\x52\x68\x20\x3d\x16\x59\x3b\x8a\x87\x6a\xb5\xaa\x46\x72\x63\x12\xff\xf7\x22\x69\xbd\x59\x27\x8e\xb3\x29\xb9\x89\xd5\x9b\x37\x33\xef\x3d\xed\xed\xed\xfc\x14\xd6\xd4\x2e\x80\x31\x80\x26\x83\x61\xe3\xf0\xac\x8d\x1c\xa4\x5a\xe1\x02\x4e\xe7\xdb\xed\x2c\xa1\x66\x9f\x6e\x9c\xb4\x0d\x84\x15\x82\x95\x2d\x42\xa7\xf3\x59\xad\xa4\x0f\x62\xd6\xe3\x6a\x68\x50\x93\x45\xa8\x94\x89\x1c\xd0\xd7\x1e\xaf\x89\x83\xdf\x88\x54\x54\x41\x7d\x0f\x93\xd1\x04\x10\xe7\xb9\xfe\x22\x31\x8a\x2b\x69\x22\x72\x46\x7e\x5b\xa3\xf7\xd4\x20\xdc\x41\xf0\xd1\x2a\x78\xff\x2e\x1f\xa9\xbd\x8c\x5a\xd3\x0d\x54\x75\x05\x3d\x17\xda\x26\x1d\xcb\x98\xe7\x1e\x65\x40\x90\x43\x07\x1d\x8d\xd9\xc0\xef\x28\x0d\x69\xc2\x06\xa4\x73\x79\x01\x31\xfb\x81\x85\x3b\xe3\x43\xea\x90\x96\x61\x58\xa2\x92\x91\x11\xb8\x6b\x11\xbe\xc4\x25\x7a\x8b\x01\xb9\xac\xad\x09\x4d\xc3\x20\x3d\x82\xa1\x96\x02\x36\x10\x3a\x08\x2b\x62\xf8\x7f\xb9\xc9\x92\x7c\xbc\xb8\x4c\x58\xb2\xd7\xc0\x0e\xd5\x1b\x31\xfb\xac\xc1\xa3\x41\xc9\xbd\x76\xaa\xb3\x41\x92\xe5\xa2\x5e\xf9\x46\x01\xfe\x90\x31\xb0\x44\x88\x9c\xe6\x64\x90\x79\xf8\x7e\xda\xe7\x15\x4e\xe0\x7d\x95\x49\x0f\xa2\xee\x2e\x07\x61\x7b\xc8\x93\xf7\x53\x84\x37\x3c\xf0\xfc\x97\x97\x58\x9c\x4d\x77\xf6\x7e\xc6\x41\x8e\x42\x22\xbe\x17\xad\x4a\xed\x6e\xce\xbd\x8f\x2f\x1c\xce\x79\xb2\x41\x43\x75\xc2\xf5\x09\x57\x0f\xb8\x4a\xd3\xe9\x39\x3b\x7c\xdc\x4b\xdf\xc8\xd6\xf4\x66\xd6\xe8\x99\x3a\x9b\x2c\xcd\xd6\xf6\x39\x29\x28\x23\x97\x68\xa6\xd8\x9b\xe1\xf7\xde\x3e\xdc\x69\x2c\x77\x39\x5f\xf5\x6d\xef\xc0\xa3\x33\x52\x21\x54\x6f\x2b\xa8\x7e\x56\x2f\x79\x54\x47\x67\xaa\x93\x73\xbe\x33\x06\xfd\xa3\xf4\x01\x59\x65\x62\x73\x3c\xa8\x02\xb6\xdb\x11\xc9\x03\x41\xa7\x35\x9e\xd8\xf4\xf5\x1a\x66\xc7\x38\x2b\x25\x9d\x5b\xc0\xb1\xb6\x87\x15\x12\x7d\xad\xf8\x35\xfc\x5d\x04\x75\xf3\x74\x39\x9d\x6e\x44\xb5\x42\xd3\x0a\x5e\xcd\x73\x44\x8e\x33\xec\x62\xf4\xc4\x08\xad\xb4\xf2\x1a\x9b\x7a\xb9\xc9\x34\xc3\x4b\xb9\x44\xbf\x26\x85\x87\x8b\xc8\x72\x90\x56\xe1\x7e\xc9\xee\xf5\x3e\xc6\xf7\xef\xa1\xc0\x4b\x56\x3f\x38\xf7\x74\x5c\x0f\x92\xa8\xae\x75\x9d\x45\x1b\x16\x70\x44\xa5\x03\x85\x4e\xfa\x50\x77\xfa\x19\x99\x46\xea\xfe\x4b\x46\x18\x0d\xaa\xd0\xf9\xaf\x7d\x56\xea\xd7\x35\xfc\xa5\x1e\x8c\x56\xf8\x1b\x00\x00\xff\xff\x0e\x0d\x26\xe8\xf3\x07\x00\x00"),
		},
		"/templates/deployment.yaml": &vfsgen۰CompressedFileInfo{
			name:             "deployment.yaml",
			modTime:          time.Time{},
			uncompressedSize: 2671,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xbc\x96\x5d\x6f\xf2\x36\x14\xc7\xef\xf9\x14\x16\xf7\xa4\x7a\xa6\x69\x9a\x72\x97\x42\x5a\x55\xa2\x80\x00\x55\xea\x15\x3a\x38\x27\x60\xcd\xb1\x2d\xfb\x24\x1b\xa3\x7c\xf7\xc9\x24\xa4\x09\x10\x56\xd6\xe9\xc9\x15\x3a\xc7\xfe\xff\xce\x8b\xed\x03\x18\xf1\x86\xd6\x09\xad\x42\x06\xc6\xb8\x87\xe2\x47\xef\x0f\xa1\x92\x90\x8d\xd0\x48\xbd\xcb\x50\x51\x2f\x43\x82\x04\x08\xc2\x1e\x63\x0a\x32\x0c\xd9\x7e\xcf\x84\xe2\x32\x4f\x90\xf5\xb9\xcc\x1d\xa1\x1d\x58\xdc\x08\x47\x76\x37\xe0\x5a\x91\xd5\x52\xa2\x0d\xd2\x5c\x4a\xbf\xa3\xcf\x02\x76\x38\x54\xdb\x9d\x01\x5e\x6a\x04\x73\x94\x08\x0e\x83\xc9\xc9\x5c\xae\x92\xb0\x46\xe9\x3c\x8e\xb1\xfd\x7e\xf0\x35\x56\xb9\xc9\x93\x3e\x98\x12\x2a\x41\x45\xec\x57\xaf\xe7\x0c\x72\xaf\x65\xd1\x48\xc1\xc1\x95\xe8\x37\x90\x39\xba\xe0\x64\x2c\xc1\x0e\x25\x72\xd2\xb6\x44\x67\x40\x7c\x3b\x6e\xc4\x72\x47\x34\x27\xa5\xf1\x95\xa8\x7e\x2b\x61\x84\x99\x91\x40\x58\xc1\x1a\x45\xf6\x1f\x28\xa5\x09\x48\x68\x55\xc3\x2b\x7c\x5a\x07\x2f\x1c\x09\x1d\x58\x2c\x84\xef\x60\xa9\x5a\x7e\xa5\x67\x0d\xea\x6f\x10\x5c\xea\x3c\x09\x84\x7e\xb0\x58\xb4\x72\xef\xde\xee\x39\xa8\x92\x73\xd3\x9f\x82\xb6\xf5\x6e\xa3\x93\xe8\x33\xc6\xf3\xa5\xa4\xdf\x21\x93\xad\xb4\x7f\xbf\x49\x90\xad\x3a\x7f\xb3\xef\x35\xea\xd4\xfb\xe3\x6f\xe4\xb9\x15\xb4\x1b\x6a\x45\xf8\x17\x85\x57\xe3\xfd\x4c\x6e\xd1\x5e\x7e\x3d\x11\x87\xb6\x10\x1c\x23\xce\x75\xae\x68\xf2\xbd\xcb\xe1\x3f\xef\x07\xa1\xd0\x36\x4a\x31\xa8\x6e\x5d\x06\x0a\x36\x68\x6b\xfb\x8d\x94\x3a\xd2\x72\x9d\x39\xfd\xf8\xa5\xd9\x1d\xc6\x44\x06\x1b\x0c\x59\xbf\x79\x5c\xbc\xc9\x5f\x18\xed\x04\x69\xbb\x63\x87\x43\x78\xe1\x26\xd8\xb0\x0f\x96\x60\x0a\xb9\x24\x16\x0c\xb7\x60\x29\x88\x8c\xa9\x5e\x19\x76\x38\xf4\xcf\x29\xb3\x5c\xca\x99\x96\x82\xef\xda\xa7\xf3\xa8\x67\x6a\x67\x3b\x3e\xa3\x2d\xb9\x76\xbe\x75\x99\x90\xac\xe0\xae\xe5\x6b\x14\x76\xa6\x2d\xb5\x40\x55\x0f\x03\x2f\xd9\x86\x1c\x41\x56\x93\xe6\x5a\x86\x6c\x39\x9c\x35\x7c\x52\x14\xa8\xd0\xb9\x99\xd5\x6b\x6c\x07\xb2\x25\x32\xcf\x78\xd6\x0d\xc6\x0c\xd0\x36\x64\x0f\xd7\xc3\x33\xc7\xa8\x2e\x7d\x16\x21\x11\x3f\x85\xe3\x74\x6e\x39\xba\x7f\x3d\x43\xf5\xca\x1b\xa7\x07\x55\x71\xbd\x37\xaf\xf1\x72\xfe\x32\x5c\xac\xa2\xd1\x68\x7e\x16\x59\xe1\xe5\x43\xd6\x0f\xbb\x5b\xd3\xbf\xaa\x39\x8e\xa3\x51\x3c\x5f\xc5\xe3\x78\xb8\x7c\x99\x4e\x56\xf1\x24\x7a\x1c\xc7\xa3\x0e\xf9\x86\x7a\xf3\x11\x41\x48\xd0\xc6\xfe\xcd\x16\x5a\x05\xa8\x60\x2d\x31\xf9\x32\x72\x12\xbd\xc6\xdf\xe1\x79\xd5\xbb\x60\x8b\x59\x34\xbc\x45\xbc\x36\x58\x3b\xd4\xa7\xcf\xab\xa7\xe9\xfc\x35\x5a\xde\x97\x80\xde\x04\xa9\xb6\x19\xd0\x4d\xe5\xb7\x78\xfe\x38\x5d\xbc\x2c\xdf\xef\x16\x2f\xd0\xae\xfd\x43\xb3\xeb\xd4\xbf\x59\x87\x27\xab\xb3\xf3\x8b\xc1\x58\x2a\x50\x26\x73\x4c\x2f\x3d\x95\x6f\x76\xbc\x3b\xa7\x59\x1c\xd4\x7f\x57\x7a\x1d\x43\x50\xe9\x04\x17\xd5\xa8\xff\xbc\x02\x4d\x6b\x78\xc7\x58\xbc\x18\x8a\x17\x3c\x48\x53\xa1\xca\xa2\x9c\xfe\x27\x54\x96\xff\x97\x43\x5a\xa2\x3d\x9f\xec\x0d\xe3\x7f\xa4\xfd\x13\x00\x00\xff\xff\xc9\x95\x7d\xe6\x6f\x0a\x00\x00"),
		},
		"/templates/poddistruptionbudget.yaml": &vfsgen۰CompressedFileInfo{
			name:             "poddistruptionbudget.yaml",
			modTime:          time.Time{},
			uncompressedSize: 434,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x94\x91\xc1\x4a\x03\x31\x10\x40\xef\xfb\x15\x43\xef\x8d\x2c\x88\x87\xbd\x29\x1e\x8b\x88\x87\xde\x67\x93\x69\x1d\x9c\x9d\x84\x64\x52\x28\x6b\xff\x5d\xd2\x50\x10\xf4\xa0\xd7\x90\xf7\xde\x24\xb3\xae\x5b\xe0\x03\xb8\x3d\x4a\xa5\xe2\x52\x0c\xcf\x5c\x72\x4d\xc6\x51\x9f\x6a\x38\x92\x39\x52\x9c\x85\x02\x5c\x2e\x03\x26\xde\x53\x2e\x1c\x75\x82\x14\x85\xfd\xf9\xee\x34\xce\x64\x38\x0e\x1f\xac\x61\x82\xd7\x9f\xfc\xb0\x90\x61\x40\xc3\x69\x00\x50\x5c\x68\x82\x75\x05\x56\x2f\x35\x10\x6c\xbc\xd4\x62\x94\xb7\x99\x8e\x5c\x2c\x9f\xb7\x3e\xaa\xe5\x28\x42\xd9\x1d\xaa\x48\x23\x36\xe0\x5a\xbd\xe3\x25\xa1\xef\x0e\xf7\x46\x42\x58\xc8\xbd\xdc\x8e\xfb\x2d\xc1\x99\xa4\xb4\x1c\xc0\xf5\x7d\x7f\x69\x75\xa8\x95\x3e\x41\x59\x03\xa9\xc1\x7d\xf3\x95\x44\xbe\xb9\x16\xd6\xc7\x13\xb2\xb4\xcf\x98\x60\x1c\x00\x0a\x09\x79\x8b\xb9\x97\x16\x34\xff\xbe\xfb\x96\xfe\x47\xfc\x66\xda\xfd\x32\xc4\x43\x1b\xa2\x99\x48\xaf\x3b\xf8\x0a\x00\x00\xff\xff\xa2\xa6\x23\xc5\xb2\x01\x00\x00"),
		},
		"/templates/rbac-leader-election.yaml": &vfsgen۰CompressedFileInfo{
			name:             "rbac-leader-election.yaml",
			modTime:          time.Time{},
			uncompressedSize: 1093,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xdc\x93\xbb\xae\x14\x31\x0c\x86\xfb\x3c\x85\xb5\x7d\xe6\x08\x89\x02\x4d\x07\x0d\x1d\xc5\x22\xd1\x7b\x93\x7f\xf7\x84\xe3\x4d\x46\x8e\x33\x08\x86\x79\x77\x34\x97\x45\x08\x04\x5a\x2e\x15\xd5\x58\x1e\xf9\xfb\xac\x5f\x0e\x0f\xe9\x1d\xb4\xa6\x92\x7b\xd2\x13\x87\x8e\x9b\x3d\x16\x4d\x9f\xd8\x52\xc9\xdd\xd3\x8b\xda\xa5\xf2\x30\x3e\x73\x4f\x29\xc7\x9e\x8e\x45\xe0\xae\x30\x8e\x6c\xdc\x3b\xa2\xcc\x57\xf4\x34\x4d\x94\x72\x90\x16\x41\x87\x20\xad\x1a\xd4\x2b\x2e\xa9\x9a\x7e\xf4\xa1\x64\xd3\x22\x02\xed\xce\x4d\x64\x99\x38\x50\x47\xf3\xec\x05\x1c\xa1\x1e\x82\xb0\xd8\x76\x5c\x1d\x38\x6c\xcc\xee\x08\x01\x57\x74\x6f\x6e\x6d\x9a\x67\x47\x24\x7c\x82\xd4\x45\x4f\x34\x4d\xfe\x3e\xf7\x36\xb4\x98\x3f\x53\x4e\x39\x22\x1b\x3d\x5f\x78\xda\x04\xb5\x77\x9e\x78\x48\xaf\xb5\xb4\x61\x25\x7b\x3a\x1c\x1c\x91\xa2\x96\xa6\x01\x7b\x2f\x94\x7c\x4e\x97\x2b\x0f\xd5\x11\x8d\xd0\xd3\xde\xbf\xc0\xd6\xaf\xa4\xba\x15\x1f\xd8\xc2\xe3\x36\xa2\x60\xc3\x5a\xb6\x21\xde\xca\xe1\xeb\xff\x08\x81\xe1\x77\xf5\x0f\xd5\xd8\xda\x4f\xb6\xf8\xc1\x73\x17\x1c\x23\xb2\x7d\x47\xdc\x97\xf7\xde\xbb\x3f\xb8\x94\x57\x29\xc7\x94\x2f\xff\xdd\xc1\x14\xc1\x11\xe7\x05\x77\x8b\xf5\x17\x91\x38\xa2\x6f\xde\xce\x3f\x0e\xa0\xb6\xd3\x7b\x04\x5b\xcf\x77\xb3\xbc\x85\x8e\x29\xe0\x65\x08\xa5\x65\xfb\x4b\xdf\x9d\x01\x7f\x09\x00\x00\xff\xff\x4f\x3b\x02\xe0\x45\x04\x00\x00"),
		},
		"/templates/rbac.yaml": &vfsgen۰CompressedFileInfo{
			name:             "rbac.yaml",
			modTime:          time.Time{},
			uncompressedSize: 2975,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xcc\x55\xc1\x6e\xdb\x30\x0c\xbd\xfb\x2b\x08\x1f\x07\x58\xc5\x80\x1d\x06\xdf\xb6\x1d\x76\x19\x76\x48\x81\x02\xc3\xd0\x03\x23\xb1\x8e\x56\x45\x32\x28\x2a\xc5\x96\xf5\xdf\x07\xcb\x49\x9a\xd5\x49\xe7\xac\x99\xbb\x93\x69\x41\x7a\x8f\x7c\x94\x1e\xb1\xb5\x57\xc4\xd1\x06\x5f\xc3\xea\x75\x71\x6b\xbd\xa9\xe1\x92\x78\x65\x35\xbd\xd3\x3a\x24\x2f\xc5\x92\x04\x0d\x0a\xd6\x05\x80\xc7\x25\xd5\xb0\x5e\x83\xf5\xda\x25\x43\x50\x6a\x97\xa2\x10\x57\x4c\x8d\x8d\xc2\xdf\x2b\x1d\xbc\x70\x70\x8e\x58\xdd\x24\xe7\xba\x13\x25\x28\xb8\xbf\xdf\x1c\x8f\x2d\xea\x1e\x43\xcd\xc8\x11\x46\x52\x9f\xb7\xcb\xfd\x2e\x87\x73\x72\xb1\xa3\x03\x58\xaf\xab\x71\x5c\xfd\xa1\x8e\xe9\x27\x78\xeb\x0d\x79\x81\x37\x3d\x5e\x87\x71\x67\x65\x01\xea\x0a\x5d\xa2\xa8\xe2\x6f\x05\x2a\xf4\x3e\x08\x8a\x0d\x3e\xf6\x07\xf6\x16\x1e\xb2\x90\xf0\x05\x97\xee\x18\x3e\x79\xd3\xfd\x54\x55\x55\xec\x6b\xca\x73\xd4\x0a\x93\x2c\x02\xdb\x1f\x19\x51\xdd\xbe\x8d\xca\x86\x8b\x9d\xda\x1f\xfa\xa2\x66\xc1\xd1\x01\xa9\x9f\x28\xb9\xc2\xa6\x61\x6a\x50\xc8\x14\xdb\xd0\x06\x3f\x4b\x8e\xba\xf3\xfa\x01\xf7\x92\x1c\x69\x09\x9c\xab\xa9\x60\x89\xa2\x17\x9f\xf6\x44\x86\x01\x4d\xce\x52\xdb\xa8\x83\xd2\x61\x79\x71\x90\xb4\x86\x52\x38\x51\x59\x70\x72\x14\x6b\xf8\x7a\x7d\xfe\xf2\xff\xfa\xa6\x9d\xfb\x0e\x9d\x45\xa2\xa2\x02\x6c\xed\x47\x0e\xa9\xed\xf4\xda\xe6\x72\x18\xb1\xbc\x2e\x00\x98\x62\x48\xac\xb3\xbc\xe5\xab\xbc\xb4\x22\x9e\x6f\x1a\xd9\x90\xe4\xaf\xb3\xb1\x0f\xee\xba\xce\xe6\x48\x33\xa1\x50\x0e\x53\x6b\xb6\xa1\x21\x47\x9b\xb0\xcd\x5b\x1f\x25\xf4\x88\x33\x6f\xa4\x15\x79\x89\x2f\x40\xbc\x33\x8b\xd8\xff\x06\xb3\x89\x22\x69\xa6\x91\x29\x8d\xe0\x39\x08\x37\xb6\x8c\x67\x5d\xf8\xf7\xd6\x1b\xeb\x9b\xff\xf7\xde\x73\x70\x34\xa3\x9b\x0e\x6e\xf8\x5a\x4f\xb4\x28\xd8\x75\xe2\x09\x99\x8a\x98\xe6\xdf\x48\x4b\x7e\x2b\x07\xa7\xd1\x24\x33\xa8\x18\x34\x76\xfc\x70\x1c\x24\xc3\x84\x86\x78\x24\xf3\xbf\x1f\x20\x7d\x3a\x13\x0e\x8f\x01\xe1\x14\x83\xa3\x3c\x52\x77\xf9\x12\xb3\xe1\x4f\x0a\x4c\x34\x17\x26\xf3\xdc\x13\x9a\x39\x27\xc1\x93\x9c\xf1\xd8\xf3\xda\xb3\xaa\x31\x3e\x73\x92\xa1\x0d\x9f\xcc\x68\x9b\x7a\x9e\x1b\xfc\x0a\x00\x00\xff\xff\x84\x8f\xc9\x65\x9f\x0b\x00\x00"),
		},
		"/templates/service.yaml": &vfsgen۰CompressedFileInfo{
			name:             "service.yaml",
			modTime:          time.Time{},
			uncompressedSize: 470,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x94\x90\xcd\x6a\xeb\x30\x10\x85\xf7\x7e\x8a\x21\x7b\x0b\x2e\xdc\x95\xb7\xdd\x96\x12\xda\x92\xbd\x2a\x9f\xa6\xa2\x63\x49\xcc\x8c\x03\xc1\xcd\xbb\x17\x45\x0a\xdd\xb4\xd0\x7a\xe7\x39\x3f\x1f\x47\xbe\xc4\x03\x44\x63\x4e\x13\x9d\xfe\x0d\xef\x31\xcd\x13\x3d\x41\x4e\x31\x60\x58\x60\x7e\xf6\xe6\xa7\x81\x28\xf9\x05\x13\x6d\x1b\xc5\x14\x78\x9d\x41\xbb\xc0\xab\x1a\x64\x14\x1c\xa3\x9a\x9c\xc7\x90\x93\x49\x66\x86\xb8\xd7\x95\xb9\x26\x76\xe4\xe8\x72\xe9\x71\x2d\x3e\xb4\x0e\xf7\x08\x86\x57\xb8\x87\xdb\xb9\xb9\xd8\xbf\x80\xb5\xe2\x88\xb6\x6d\xfc\x1d\xab\x85\x2a\xe9\x83\x52\x4c\x33\x92\xd1\xff\xda\xa7\x05\xa1\x76\xd9\xb9\x74\xec\xc1\xf3\x0a\x75\xda\xf6\xb9\x2a\x34\x70\xc9\x62\x9d\x3b\x5e\x7f\xbe\xf5\x57\xa1\xf9\xeb\x67\x5e\x8e\xb0\xfd\xd5\xbc\xc0\x24\x06\xed\x4a\x91\x6c\x39\x64\x9e\xe8\xf9\x6e\xdf\x6f\xed\xfd\xde\xcc\xca\xf8\x65\x56\x30\x82\x65\xf9\xe3\xe2\x5b\xec\xfe\xa7\xe5\x9f\x01\x00\x00\xff\xff\xb2\xbe\xa1\xef\xd6\x01\x00\x00"),
		},
		"/values.yaml": &vfsgen۰CompressedFileInfo{
			name:             "values.yaml",
			modTime:          time.Time{},
			uncompressedSize: 800,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x64\x92\xc1\x6f\xd5\x30\x0c\xc6\xef\xf9\x2b\xac\xb7\xf3\x1b\xdd\xc6\xd0\x94\xdb\x63\x9b\x10\x12\x93\x26\x0d\x0e\x08\x71\xc8\x52\xb7\x18\xd2\xb8\xd8\x4e\x1f\x65\xda\xff\x8e\x9a\x32\x98\xf6\x6e\xcd\xaf\xdf\x97\xf8\xb3\x7d\x04\x57\xd8\x85\x92\x0c\xa6\x90\x0a\x2a\x74\x2c\x10\x53\x51\x43\xd9\x0a\xf6\xa4\x26\xf3\x36\x72\x36\xe1\x94\x50\xdc\x11\x7c\xfc\x46\x0a\xa4\x10\xe0\xf3\xee\xe6\xc3\xb6\x63\x19\x82\x19\xb6\xd0\x51\xc2\x63\xb7\xdc\x18\x53\x10\x84\x29\x08\x85\xfb\x84\x0a\xc6\x70\x8f\x30\x06\x55\x6c\x81\xb2\x31\xcc\x5c\x04\x0c\x87\x31\x05\x43\x3d\x76\x4e\x70\x4c\x14\x83\x7a\x38\x71\x8e\xd4\x88\xbd\x03\x10\x9c\x48\x89\xb3\x87\xcd\xc6\xb9\x91\xdb\x5d\xce\x6c\xc1\x88\xb3\x7a\xf8\xf2\xb5\xb2\x3b\x8c\x45\xc8\xe6\x4b\xce\x86\xbf\xac\xfa\x4a\xde\xe9\x27\x45\xf1\xf0\xe6\xfc\xfc\xec\xf5\x13\x7a\x27\x5c\xc6\x27\xa6\x87\xbe\x90\x12\xef\x6f\x85\x26\x4a\xd8\xe3\xb5\xc6\x90\xea\x63\x1e\xba\x90\x14\x1d\x0d\xa1\xc7\xb5\xb0\x91\x95\x8c\x65\xf6\xb0\x0f\x73\xfe\xdd\xc8\xab\x97\x6d\x73\x00\x16\x7a\x0f\x43\x58\xb0\x03\x18\x4b\x4a\xb7\x9c\x28\xce\x1e\x76\x69\x1f\x66\x75\x2e\x73\x8b\x77\x98\x30\x1a\x8b\x87\x87\x47\x17\xba\x8e\x32\xd9\x5c\x0f\xc6\x09\xe5\x59\x5c\x41\xe5\x22\x11\x75\xad\xe1\x67\x41\xb5\xfa\x0d\x30\xe0\x50\xab\xd9\x9c\x34\xcd\x0d\x6d\x2a\x8b\x63\x59\xc1\xb0\x9c\x13\x0d\x74\xa0\x3e\x7d\xa9\x3e\xab\x6a\xa7\x28\x13\xc5\x9a\xd5\xe6\x11\x3d\x5c\xae\xe9\xde\xdf\x2e\x41\x58\xcc\xc3\x45\x73\xd1\xfc\x13\xee\x62\xe4\x92\xd7\x26\x3e\x9f\xd1\xc3\x63\x9d\xd1\x15\xa9\x94\x71\x61\x6f\x4b\xdb\x63\xd5\x61\x5e\xb6\xa3\x7d\xea\xad\xfb\xbf\x65\xcb\xdf\x84\xa1\x45\xb9\x5e\x3a\xb3\x0c\xa0\x56\xf8\xc2\xb1\xa0\x1c\x06\xf4\xb0\x39\x58\xd9\xd5\xbe\xc5\xbf\xfe\x9a\x9f\xfb\xf5\x9a\x75\x65\x3d\x7c\x57\xce\x15\x4c\x28\xf7\xcb\x38\x67\x0f\x8d\x03\xd8\xb3\xfc\x40\x51\x0f\xa7\xee\x4f\x00\x00\x00\xff\xff\x18\x75\xb4\xec\x20\x03\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/.helmignore"].(os.FileInfo),
		fs["/Chart.yaml"].(os.FileInfo),
		fs["/crds"].(os.FileInfo),
		fs["/examples"].(os.FileInfo),
		fs["/templates"].(os.FileInfo),
		fs["/values.yaml"].(os.FileInfo),
	}
	fs["/crds"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/crds/crds.yaml"].(os.FileInfo),
	}
	fs["/examples"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/examples/test.yaml"].(os.FileInfo),
	}
	fs["/templates"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/templates/_helpers.tpl"].(os.FileInfo),
		fs["/templates/deployment.yaml"].(os.FileInfo),
		fs["/templates/poddistruptionbudget.yaml"].(os.FileInfo),
		fs["/templates/rbac-leader-election.yaml"].(os.FileInfo),
		fs["/templates/rbac.yaml"].(os.FileInfo),
		fs["/templates/service.yaml"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
