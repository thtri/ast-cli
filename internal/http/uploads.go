package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	uploads "github.com/checkmarxDev/uploads/api/rest/v1"
	"github.com/pkg/errors"
)

type UploadsWrapper interface {
	Create(sourcesFile string) (*string, error)
}

type UploadsHttpWrapper struct {
	url string
}

func (u UploadsHttpWrapper) Create(sourcesFile string) (*string, error) {
	var body bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&body)

	file, err := os.Open(sourcesFile)
	if err != nil {
		return nil, errors.Errorf("Failed to open file %s: %s", sourcesFile, err.Error())
	}
	// Close the file later
	defer file.Close()

	// Initialize the file field
	var fileWriter io.Writer
	fileWriter, err = multiPartWriter.CreateFormFile("sources", sourcesFile)
	if err != nil {
		return nil, errors.Errorf("Failed creating FormFile - %s", err.Error())
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, errors.Errorf("Failed to copy file: %s", err.Error())
	}
	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	var req *http.Request
	req, err = http.NewRequest("POST", u.url, &body)
	if err != nil {
		return nil, errors.Errorf("Requesting error model failed - %s", err.Error())
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	var client = &http.Client{
		Timeout: time.Second * time.Duration(5),
	}
	var resp *http.Response
	fmt.Printf("Uploading to %s\n", u.url)
	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.Errorf("Invoking HTTP request failed - %s", err.Error())
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	switch resp.StatusCode {
	case http.StatusBadRequest:
		errorModel := uploads.ErrorModel{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&errorModel)
		if err != nil {
			return nil, errors.Errorf("Parsing error model failed - %s", err.Error())
		}
		return nil, errors.Errorf("%d - %s", errorModel.Code, errorModel.Message)

	case http.StatusCreated:
		model := uploads.UploadModel{}
		err = decoder.Decode(&model)
		if err != nil {
			return nil, errors.Errorf("Parsing upload model failed - %s", err.Error())
		}
		return &model.URL, nil

	default:
		return nil, errors.Errorf("Unknown response status code %d", resp.StatusCode)
	}
}

func NewUploadsHttpWrapper(url string) UploadsWrapper {
	return &UploadsHttpWrapper{
		url: url,
	}
}
