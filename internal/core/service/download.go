// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package service

// type FileRepository interface {
// 	Get()
// }
//
// // Downloader streams a file to an endpoint.
// type Downloader interface {
// 	io.Writer
// }
//
// type downloader struct {
// 	// data is the data that is streamed in the request.
// 	data []byte
// 	// filePath is the file to download.
// 	filePath string
// 	// repo is the repository that serves the file, such as the file system
// 	// or a remote service like S3.
// 	repo FileRepository
// }
//
// func NewDownloader(filePath string) Downloader {
// 	return &downloader{
// 		data:     []byte{},
// 		filePath: filePath,
// 	}
// }
//
// func (d *downloader) Write(p []byte) (int, error) {
// 	return 0, nil
// }
