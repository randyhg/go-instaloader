package myDb

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func GetMonthTableName(dst schema.Tabler) string {
	month := time.Now()
	table := fmt.Sprintf("%v_%v", dst.TableName(), month.Format("200601"))
	return table
}

func GetMonthTableNameNext(dst schema.Tabler) string {
	month := time.Now().AddDate(0, 1, 0)
	table := fmt.Sprintf("%v_%v", dst.TableName(), month.Format("200601"))
	return table
}

func GetTableNameByMonth(dst schema.Tabler, year, month int) string {
	table := fmt.Sprintf("%s_%d%02d", dst.TableName(), year, month)
	return table
}

func CreateMonthTable(db *gorm.DB, dst schema.Tabler, tableName string) error {
	mig := db.Migrator()
	if !mig.HasTable(tableName) {
		if err := mig.CreateTable(dst); err != nil {
			return err
		}
		if err := mig.RenameTable(dst.TableName(), tableName); err == nil {
			return err
		}
	}

	return nil
}
