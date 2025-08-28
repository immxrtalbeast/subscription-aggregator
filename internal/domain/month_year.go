package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	if my.Year == 0 && my.Month == 0 {
		return nil, nil
	}
	return my.String(), nil
}

func (my *MonthYear) Scan(value interface{}) error {
	if value == nil {
		my.Year = 0
		my.Month = 0
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

func (my MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(my.String())
}

func (my *MonthYear) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseMonthYear(s)
	if err != nil {
		return err
	}

	my.Year = parsed.Year
	my.Month = parsed.Month
	return nil
}
func (my MonthYear) ToTime() time.Time {
	return time.Date(my.Year, time.Month(my.Month), 1, 0, 0, 0, 0, time.UTC)
}

func FromTime(t time.Time) MonthYear {
	return MonthYear{
		Year:  t.Year(),
		Month: int(t.Month()),
	}
}
