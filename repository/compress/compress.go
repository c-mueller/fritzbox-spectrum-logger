package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Compress(data []byte) ([]byte, error) {
	var outputBuffer bytes.Buffer
	compressionWriter := gzip.NewWriter(&outputBuffer)
	_, err := compressionWriter.Write(data)
	if err != nil {
		return nil, err
	}
	compressionWriter.Close()

	return outputBuffer.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	inputBuffer := bytes.NewReader(data)
	compressionReader, err := gzip.NewReader(inputBuffer)
	if err != nil {
		return nil, err
	}

	defer compressionReader.Close()

	return ioutil.ReadAll(compressionReader)
}
