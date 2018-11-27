package fritz

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const supportDataURL = "/cgi-bin/firmwarecfg"

func (s *Session) DownloadSupportData() ([]byte, error) {
	queryURL, err := s.getUrl(supportDataURL)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	w.WriteField("sid", s.sessionInfo.SID)
	w.WriteField("SupportData", "")

	req, err := http.NewRequest("POST", queryURL.String(), &buffer)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "*/*")

	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	return data, err
}
