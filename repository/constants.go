package repository

import "errors"

const (
	SpectrumListBucketName = "Spectra"
)

var BucketNotFoundError = errors.New("repository: Bucket not Found")
