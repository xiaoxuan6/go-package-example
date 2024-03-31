package time

import (
    "strconv"
    "time"
)

type LocalTime struct {
    Name string
}

func NewTime() *LocalTime {
    return &LocalTime{
        Name: "Asia/Shanghai",
    }
}

func (l *LocalTime) getLoc() *time.Location {
    loc, err := time.LoadLocation(l.Name)
    if err != nil {
        return nil
    }

    return loc
}

func (l *LocalTime) LoadLocation(name string) *LocalTime {
    l.Name = name
    return l
}

func (l LocalTime) Year() string {
    return time.Now().In(l.getLoc()).Format("2006")
}

func (l LocalTime) YearInt() int {
    return l.toInt(l.Year())
}

func (l LocalTime) Month() string {
    return time.Now().In(l.getLoc()).Format("01")
}

func (l LocalTime) MonthInt() int {
    return l.toInt(l.Month())
}

func (l LocalTime) Day() string {
    return time.Now().In(l.getLoc()).Format("02")
}

func (l LocalTime) DayInt() int {
    return l.toInt(l.Day())
}

func (l LocalTime) Hour() string {
    return time.Now().In(l.getLoc()).Format("15")
}

func (l LocalTime) HourInt() int {
    return l.toInt(l.Hour())
}

func (l LocalTime) Minutes() string {
    return time.Now().In(l.getLoc()).Format("04")
}

func (l LocalTime) MinutesInt() int {
    return l.toInt(l.Minutes())
}

func (l LocalTime) Seconds() string {
    return time.Now().In(l.getLoc()).Format("05")
}

func (l LocalTime) SecondsInt() int {
    return l.toInt(l.Seconds())
}

func (l LocalTime) Date() string {
    return time.Now().In(l.getLoc()).Format("2006-01-02")
}

func (l LocalTime) Time() string {
    return time.Now().In(l.getLoc()).Format("2006-01-02 15:04:05")
}

func (l LocalTime) DateTime() time.Time {
    return time.Date(
        l.YearInt(),
        l.MonthEn(),
        l.DayInt(),
        l.HourInt(),
        l.MinutesInt(),
        l.SecondsInt(),
        0,
        l.getLoc(),
    )
}

func (l LocalTime) MonthEn() time.Month {
    var month time.Month
    switch l.MonthInt() {
    case 1:
        month = time.January
    case 2:
        month = time.February
    case 3:
        month = time.March
    case 4:
        month = time.April
    case 5:
        month = time.May
    case 6:
        month = time.June
    case 7:
        month = time.July
    case 8:
        month = time.August
    case 9:
        month = time.September
    case 10:
        month = time.October
    case 11:
        month = time.November
    case 12:
        month = time.December
    }
    return month
}

func (l LocalTime) toInt(target string) int {
    t, _ := strconv.Atoi(target)
    return t
}
