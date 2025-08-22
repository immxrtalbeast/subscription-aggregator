package domain

import "fmt"

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
