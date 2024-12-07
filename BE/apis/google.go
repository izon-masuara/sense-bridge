package apis

import (
	"encoding/json"
	"os/exec"
)

type Response struct {
	RequestId         string   `json:"requestId"`
	Results           []Result `json:"results"`
	TotalBilledTime   string   `json:"totalBilledTime"`
	UsingLegacyModels bool     `json:"usingLegacyModels"`
}

type Result struct {
	Alternatives  []Alternative `json:"alternatives"`
	LanguageCode  string        `json:"languageCode"`
	ResultEndTime string        `json:"resultEndTime"`
}

type Alternative struct {
	Confidence float64 `json:"confidence"`
	Transcript string  `json:"transcript"`
}

func SpeechToText(path string) (string, error) {
	cmd := exec.Command("gcloud", "ml", "speech", "recognize", path, "--language-code=id-ID")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	jsonData := string(output)

	var response Response

	err = json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		return "", err
	}

	return response.Results[0].Alternatives[0].Transcript, nil
}
