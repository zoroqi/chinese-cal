package main

import (
	"flag"
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"math"
	"strconv"
	"strings"
	"time"
)

func main() {
	three := flag.Bool("3", false, "")
	y := flag.Int("year", 0, "year")
	m := flag.Int("month", 0, "month")
	flag.Parse()

	print := func(t time.Time) {
		lines := layout(month(t))
		for _, s := range lines {
			fmt.Println(s)
		}
		fmt.Println()
	}

	now := time.Now()
	if *y != 0 {
		if *y > 0 && *y < 10000 {
			now = now.AddDate(*y-now.Year(), 0, 0)
		}
	}
	if *m != 0 {
		if *m > 0 && *m < 13 {
			now = now.AddDate(0, *m-int(now.Month()), 0)
			// 月份不一致, 说明是计算日期偏移到下个月了, 需要调整到指定月份.
			// 2021-01-31 + 1 month = 2021-02-31(2021-03-03), 而我期望的是 2021-02-28
			for int(now.Month()) != *m {
				now = now.AddDate(0, 0, -1)
			}
		}
	}
	if *three {
		print(firstDay(now).AddDate(0, -1, 0))
	}
	print(now)
	if *three {
		print(firstDay(now).AddDate(0, 1, 0))
	}
}

func firstDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 1, t.Minute(), t.Second(), 0, t.Location())
}

var week = []string{"日", "一", "二", "三", "四", "五", "六"}

const blockWidth = 8
const width = blockWidth * 7

var space = strings.Repeat(" ", width)

var today = time.Now()

func layout(m []day) []string {
	// 因为需要显示中文, 每一个日期包含三个字, 加上左右两个空格, 一行半角共 56 个字符.
	lines := []string{}
	// 26 2 2 4 26
	lines = append(lines, block(fmt.Sprintf("%02d  %d", m[0].Month(), m[0].Year()), width, false))
	lines = append(lines, block(strings.Join(week, space[:blockWidth-2]), width, false))
	csb := strings.Builder{}
	nsb := strings.Builder{}

	if m[0].Weekday() != 0 {
		nsb.WriteString(space[:m[0].Weekday()*blockWidth])
		csb.WriteString(space[:m[0].Weekday()*blockWidth])
	}

	for _, d := range m {
		if d.Weekday() == 0 {
			if csb.String() != "" {
				lines = append(lines, nsb.String())
				lines = append(lines, csb.String())
			}
			csb.Reset()
			nsb.Reset()
		}
		nsb.WriteString(d.nString())
		csb.WriteString(d.cString())
	}
	lines = append(lines, nsb.String())
	lines = append(lines, csb.String())

	return lines
}

// s 长度不足 length 的时候进行空格填充, 尽量保证 s 是居中的.
func block(s string, length int, color bool) string {
	l := 0
	for _, r := range s {
		if r <= 255 {
			l += 1
		} else {
			l += 2
		}
	}

	left := max(0, int(math.Ceil(float64(length-l)/2.0)))
	right := max(0, int(math.Floor(float64(length-l)/2.0)))

	if color {
		return fmt.Sprintf("%c[7;40;37m%s%s%s%c[0m", 0x1B, space[:left], s, space[:right], 0x1B)
	} else {
		return fmt.Sprintf("%s%s%s", space[:left], s, space[:right])
	}
}

type day struct {
	time.Time
	lunar *calendar.Lunar
	today bool
}

func (d day) nString() string {
	return block(strconv.Itoa(d.Day()), blockWidth, d.today)
}

func (d day) cString() string {
	s := ""
	if d.lunar.GetJieQi() != "" {
		s = d.lunar.GetJieQi()
	} else {
		s = d.lunar.GetMonthInChinese() + d.lunar.GetDayInChinese()
	}
	return block(s, blockWidth, d.today)
}

func isToday(t time.Time) bool {
	return t.Year() == today.Year() && t.Month() == today.Month() && t.Day() == today.Day()
}

func month(t time.Time) []day {
	first := firstDay(t)
	lunar := calendar.NewLunarFromDate(first)
	current := first
	days := []day{}

	cm := current.Month()
	for cm == current.Month() {
		days = append(days, day{
			Time:  current,
			today: isToday(current),
			lunar: lunar,
		})
		current.YearDay()
		current = current.AddDate(0, 0, 1)
		lunar = lunar.Next(1)
	}

	return days
}
