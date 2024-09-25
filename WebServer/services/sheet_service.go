package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-instaloader/config"
	"go-instaloader/google_auth"
	"go-instaloader/models"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
	"strconv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type sheetService struct {
	srv *sheets.Service
}

type SheetRow struct {
	No string `json:"no"`
	// Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Url      string `json:"url"`
	Status   string `json:"status"`
}

var ConfigColumnRows *ConfigRow

type ConfigRow struct {
	UsernameCol  int `json:"usernameCol"`
	StatusCol    int `json:"statusCol"`
	RemarkCol    int `json:"remarkCol"`
	TalentUrlCol int `json:"talentUrlCol"`
	TalentCount  int `json:"talentCount"`
}

func newSheetService() *sheetService {
	client, err := google_auth.GetHttpClient()
	if err != nil {
		rlog.Error("unable to get http client:", err.Error())
		return nil
	}

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		rlog.Errorf("unable to retrieve Sheets client: %v", err)
		return nil
	}
	return &sheetService{srv: srv}
}

func (s *sheetService) GetConfigColumnRows() error {
	dataRange := fmt.Sprintf("%s!%s", config.Instance.ConfigSheetName, config.Instance.ConfigCellRange)

	resp, err := s.srv.Spreadsheets.Values.Get(config.Instance.SpreadSheetId, dataRange).Do()
	if err != nil {
		rlog.Error(err)
		return fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		rlog.Error("no data found")
		return errors.New("no data found")
	}

	for _, rows := range resp.Values {
		usernameColumn, ok1 := rows[0].(string)
		if !ok1 {
			return errors.New("failed to parse username column")
		}

		statusColumn, ok2 := rows[1].(string)
		if !ok2 {
			return errors.New("failed to parse status column")
		}

		remarkColumn, ok3 := rows[2].(string)
		if !ok3 {
			return errors.New("failed to parse remark column")
		}

		talentUrlColumn, ok4 := rows[3].(string)
		if !ok4 {
			return errors.New("failed to parse talent url column")
		}

		talentCountStr, ok5 := rows[4].(string)
		if !ok5 {
			return errors.New("failed to parse talent count column")
		}

		talentCount, _ := strconv.Atoi(talentCountStr)

		ConfigColumnRows = &ConfigRow{
			UsernameCol:  charToNumber(usernameColumn),
			StatusCol:    charToNumber(statusColumn),
			RemarkCol:    charToNumber(remarkColumn),
			TalentUrlCol: charToNumber(talentUrlColumn),
			TalentCount:  talentCount,
		}

		// ConfigColumnRows.UsernameCol = charToNumber(usernameColumn)
		// ConfigColumnRows.StatusCol = charToNumber(statusColumn)
		// ConfigColumnRows.RemarkCol = charToNumber(remarkColumn)
		// ConfigColumnRows.TalentUrlCol = charToNumber(talentUrlColumn)
		// ConfigColumnRows.TalentCount = talentCount
	}
	return nil
}

func (s *sheetService) GetTalents(ctx context.Context, fetchRange string) ([]*models.Talent, error) {
	if err := s.GetConfigColumnRows(); err != nil {
		return nil, err
	}

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
			rlog.Errorf("record %s is not complete", row[0].(string))

			if err = s.UpdateTalentStatus(ctx, models.StatusFail, sheetRow.Username, "record is not complete!"); err != nil {
				rlog.Errorf("unable to update status: %v", err)
			}

			continue
		}

		talent := &models.Talent{
			SheetId:  sheetRow.No,
			Username: sheetRow.Username,
			Url:      sheetRow.Url,
			Status:   models.StatusOnProcess,
		}

		if err = s.UpdateTalentStatus(ctx, models.StatusOnProcess, sheetRow.Username, ""); err != nil {
			rlog.Errorf("unable to update status: %v", err)
		}

		if len(talent.SheetId) > 0 {
			byt, err := json.Marshal(&talent)
			if err != nil {
				rlog.Error(err)
			} else {
				fwRedis.RedisStore().LPush(context.Background(), models.RedisJobQueueKey, string(byt))
			}
		}

		talents = append(talents, talent)
	}
	return talents, nil
}

func (s *sheetService) UpdateTalentStatus(ctx context.Context, status int, username, remark string) error {
	if err := s.GetConfigColumnRows(); err != nil {
		return err
	}
	spreadSheetId := config.Instance.SpreadSheetId

	resp, err := s.srv.Spreadsheets.Values.Get(spreadSheetId, config.Instance.MaxFetchRange).Do()
	if err != nil {
		rlog.Errorf("Unable to retrieve data from sheet: %v", err)
		return err
	}

	conditionColumn := ConfigColumnRows.UsernameCol - 1
	conditionValue := username
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

	rangeToUpdate := fmt.Sprintf("%s!%s%d", config.Instance.SheetName, numberToChar(ConfigColumnRows.StatusCol), updateRow)
	//rlog.Info("updating ", rangeToUpdate)
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{status, remark},
		},
	}
	_, err = s.srv.Spreadsheets.Values.Update(spreadSheetId, rangeToUpdate, valueRange).ValueInputOption("USER_ENTERED").Do()
	return err
}

func charToNumber(c string) int {
	return int(c[0] - 'A' + 1)
}

func numberToChar(n int) string {
	return string(rune('A' + n - 1))
}

func parseRow(row []interface{}) *SheetRow {
	if len(row) <= 3 {
		return &SheetRow{}
	}
	var sheetRow SheetRow
	// no
	no, ok1 := row[0].(string)
	if ok1 {
		sheetRow.No = no
	}

	// username
	username, ok3 := row[(ConfigColumnRows.UsernameCol - 1)].(string)
	if ok3 {
		sheetRow.Username = username
	}

	// url
	url, ok4 := row[ConfigColumnRows.TalentUrlCol-1].(string)
	if ok4 {
		sheetRow.Url = url
	}

	// status
	status, ok5 := row[(ConfigColumnRows.StatusCol - 1)].(string)
	if ok5 {
		sheetRow.Status = status
	}
	return &sheetRow
}
