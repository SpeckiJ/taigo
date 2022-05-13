package taigo

import (
	"fmt"
	"net/http"
	"strconv"
)

// ExporterService is a handle to actions related to Export/Import
//
// https://taigaio.github.io/taiga-doc/dist/api.html#export-import-export-dump
type ExporterService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// CreateExport creates a new Export
// https://docs.taiga.io/api.html#export-import-export-dump
func (s *ExporterService) CreateExport(project *Project) (*ExportDetail, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(project.ID))
	type ExportRaw struct {
		Url      string `json:"url"`
		ExportId string `json:"export_id"`
	}
	var ex ExportRaw

	res, err := s.client.Request.Get(url, &ex)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusAccepted {
		return &ExportDetail{
			fmt.Sprintf("%s/exports/%d/%s-%s.json", s.client.MediaURL, project.ID, project.Slug, ex.ExportId)}, nil
	} else {
		return &ExportDetail{Url: ex.Url}, nil
	}
}
