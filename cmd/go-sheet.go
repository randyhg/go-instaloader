package cmd

import (
	"context"
	"errors"
	"fmt"
	"go-instaloader/config"
	"go-instaloader/models"
	"go-instaloader/utils/rlog"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
)

type SheetRow struct {
	No       int    `json:"no"`
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Url      string `json:"url"`
	Status   string `json:"status"`
}

func GetTalents(client *http.Client, ctx context.Context) ([]*models.Talent, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(config.Instance.SpreadSheetId, config.Instance.MaxFetchRange).Do()
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
				if err = UpdateTalentStatus(client, ctx, models.StatusFail, sheetRow.Uuid, "record is not complete!"); err != nil {
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

		if err = UpdateTalentStatus(client, ctx, models.StatusOnProcess, sheetRow.Uuid, ""); err != nil {
			rlog.Error(fmt.Sprintf("unable to update status: %v", err))
		}

		talents = append(talents, talent)
	}
	return talents, nil
}

func parseRow(row []interface{}) *SheetRow {
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

func UpdateTalentStatus(client *http.Client, ctx context.Context, status int, uuid, remark string) error {
	spreadSheetId := config.Instance.SpreadSheetId
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadSheetId, config.Instance.MaxFetchRange).Do()
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
	_, err = srv.Spreadsheets.Values.Update(spreadSheetId, rangeToUpdate, valueRange).ValueInputOption("USER_ENTERED").Do()
	return err
}

func UpdateTalentStoryImg(client *http.Client, ctx context.Context, imgUrl, uuid string) error {
	spreadSheetId := config.Instance.SpreadSheetId
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadSheetId, config.Instance.MaxFetchRange).Do()
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

	rangeToUpdate := fmt.Sprintf("%s!%s%d", config.Instance.SheetName, "G", updateRow)
	//rlog.Info("updating ", rangeToUpdate)
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{imgUrl},
		},
	}
	_, err = srv.Spreadsheets.Values.Update(spreadSheetId, rangeToUpdate, valueRange).ValueInputOption("USER_ENTERED").Do()
	return err
}
