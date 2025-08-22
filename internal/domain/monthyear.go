package domain

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type MonthYear struct {
	Year  int
	Month int
}

func NewMonthYear(year, month int) (MonthYear, error) {
	if year < 2000 || year > 2100 {
		return MonthYear{}, fmt.Errorf("incorrect year: %d", year)
	}
	if month < 1 || month > 12 {
		return MonthYear{}, fmt.Errorf("incorrect month: %d", month)
	}
	return MonthYear{Year: year, Month: month}, nil
}

func (my MonthYear) String() string {
	return fmt.Sprintf("%02d-%d", my.Month, my.Year)
}

func (my MonthYear) IsBefore(other MonthYear) bool {
	if my.Year != other.Year {
		return my.Year < other.Year
	}
	return my.Month < other.Month
}

func ParseMonthYear(s string) (MonthYear, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return MonthYear{}, fmt.Errorf("incorrect format. Expecting MM-YYYY")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return MonthYear{}, fmt.Errorf("month should be an int")
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return MonthYear{}, fmt.Errorf("year should be an int")
	}

	return NewMonthYear(year, month)
}

///GORM

func (my MonthYear) Value() (driver.Value, error) {
	return my.String(), nil
}

func (my *MonthYear) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("MonthYear: expected string, getting %T", value)
	}

	temp, err := ParseMonthYear(s)
	if err != nil {
		return err
	}
	my.Year = temp.Year
	my.Month = temp.Month

	return nil
}
func (MonthYear) GormDataType() string {
	return "char(7)" // MM-YYYY = 7 символов
}
