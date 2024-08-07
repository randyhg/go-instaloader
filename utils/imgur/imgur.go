package imgur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const imgurUploadURL = "https://api.imgur.com/3/upload"
const imgurClientId = "1b728fe145c1dba"

type ImgurResponse struct {
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

func UploadImageToImgur(imgUrl string) (string, error) {
	httpResp, httpErr := http.Get(imgUrl)
	if httpErr != nil {
		return "", fmt.Errorf("failed to download image: %v", httpErr)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: received non-200 response code %d", httpResp.StatusCode)
	}

	// Open the image file
	//file, err := os.Open(imgUrl)
	//if err != nil {
	//	return "", fmt.Errorf("failed to open image file: %v", err)
	//}
	//defer file.Close()

	// Create a buffer to store the image data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create the form file field
	part, err := writer.CreateFormFile("image", imgUrl)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the file content to the form file field
	_, err = io.Copy(part, httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", imgurUploadURL, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Client-ID "+imgurClientId)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var imgurResp ImgurResponse
	err = json.Unmarshal(body, &imgurResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// Check if the upload was successful
	if !imgurResp.Success {
		return "", fmt.Errorf("failed to upload image: status %d", imgurResp.Status)
	}

	// Return the image link
	return imgurResp.Data.Link, nil
}
