package services

import (
	"context"
	"errors"
	"fmt"
	"go-instaloader/config"
	"go-instaloader/google_auth"
	"go-instaloader/models"
	"go-instaloader/utils/rlog"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
)

type sheetService struct {
	srv *sheets.Service
}

type SheetRow struct {
	No       int    `json:"no"`
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Url      string `json:"url"`
	Status   string `json:"status"`
}

func newSheetService() *sheetService {
	client, err := google_auth.GetHttpClient()
	if err != nil {
		rlog.Error("Unable to get http client:", err.Error())
		return nil
	}

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		rlog.Error(fmt.Sprintf("unable to retrieve Sheets client: %v", err))
		return nil
	}
	return &sheetService{srv: srv}
}

func (s *sheetService) GetTalents(ctx context.Context, fetchRange string) ([]*models.Talent, error) {
	fetchRange = fmt.Sprintf("%s!%s", config.Instance.SheetName, fetchRange)

	resp, err := s.srv.Spreadsheets.Values.Get(config.Instance.SpreadSheetId, fetchRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		return nil, errors.New("no data found")
	}

	var talents []*models.Talent
	for _, row := range resp.Values {
		sheetRow := parseRow(row)

		// username is empty
		if sheetRow.Username == "" {
			// update the status column to fail
			rlog.Error(fmt.Sprintf("record %s is not complete", row[0].(string)))

			if sheetRow.Uuid != "" {
				if err = s.UpdateTalentStatus(ctx, models.StatusFail, sheetRow.Uuid, "record is not complete!"); err != nil {
					rlog.Error(fmt.Sprintf("unable to update status: %v", err))
				}
			}
			continue
		}

		talent := &models.Talent{
			Uuid:     sheetRow.Uuid,
			Username: sheetRow.Username,
			Url:      sheetRow.Url,
			Status:   models.StatusOnProcess,
		}

		if err = s.UpdateTalentStatus(ctx, models.StatusOnProcess, sheetRow.Uuid, ""); err != nil {
			rlog.Error(fmt.Sprintf("unable to update status: %v", err))
		}

		talents = append(talents, talent)
	}
	return talents, nil
}

func (s *sheetService) UpdateTalentStatus(ctx context.Context, status int, uuid, remark string) error {
	spreadSheetId := config.Instance.SpreadSheetId

	resp, err := s.srv.Spreadsheets.Values.Get(spreadSheetId, config.Instance.MaxFetchRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	conditionColumn := 1
	conditionValue := uuid
	var updateRow int
	for rowIndex, row := range resp.Values {
		if len(row) > conditionColumn && row[conditionColumn] == conditionValue {
			//rlog.Info(fmt.Sprintf("matching row %v", row))
			updateRow = rowIndex + 2
		}
	}

	if updateRow <= 0 {
		return errors.New("no data found")
	}

	rangeToUpdate := fmt.Sprintf("%s!%s%d", config.Instance.SheetName, config.Instance.StatusColumn, updateRow)
	//rlog.Info("updating ", rangeToUpdate)
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{status, remark},
		},
	}
	_, err = s.srv.Spreadsheets.Values.Update(spreadSheetId, rangeToUpdate, valueRange).ValueInputOption("USER_ENTERED").Do()
	return err
}

func parseRow(row []interface{}) *SheetRow {
	if len(row) <= 3 {
		return &SheetRow{}
	}
	var sheetRow SheetRow
	// no
	no, ok1 := row[0].(int)
	if ok1 {
		sheetRow.No = no
	}

	// uuid
	uuid, ok2 := row[1].(string)
	if ok2 {
		sheetRow.Uuid = uuid
	}

	// username
	username, ok3 := row[2].(string)
	if ok3 {
		sheetRow.Username = username
	}

	// url
	url, ok4 := row[3].(string)
	if ok4 {
		sheetRow.Url = url
	}

	// status
	status, ok5 := row[4].(string)
	if ok5 {
		sheetRow.Status = status
	}
	return &sheetRow
}
