package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(bookCmd)

}

var bookCmd = &cobra.Command{
	Use:   "book {day} {time}",
	Short: "Create a booking",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		bs := getBookingService()

		date, err := extractDate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		start, end, err := extractTimes(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		available, err := bs.Available(date.Add(start), date.Add(end))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(available)

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if a := len(args); a > 2 || a == 1 {
			return fmt.Errorf("wrong number of arguments")
		}
		return nil
	},
}

func extractTimes(s string) (time.Duration, time.Duration, error) {
	parts := strings.Split(s, "-")
	if len(parts) == 2 {
		start, err := extractTime(parts[0])
		if err != nil {
			return time.Second, time.Second, err
		}
		end, err := extractTime(parts[1])
		if err != nil {
			return time.Second, time.Second, err
		}
		return start, end, nil
	}else{
		switch s {
		case "lunch":
			return time.Hour * 12, time.Hour * 13, nil
		default:
			return time.Second, time.Second, fmt.Errorf("failed to parse times from %s, s")
		}
	}
}

func extractTime(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	hour := ""
	minute := ""
	if len(parts) == 2 {
		hour = parts[0]
		minute = parts[1]
	} else if len(s) == 4 {
		hour = s[:2]
		minute = s[2:]
	} else if len(s) <= 2 {
		hour = s
		minute = "0"
	}
	h, err := strconv.Atoi(hour)
	if err != nil {
		return time.Nanosecond, err
	}
	m, err := strconv.Atoi(minute)
	if err != nil {
		return time.Nanosecond, err
	}
	return time.Hour * time.Duration(h) + time.Minute * time.Duration(m), nil
}

func extractDate(s string) (time.Time, error) {
	switch n := time.Now(); strings.ToLower(s) {
	case "today":
		return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location()), nil
	case "tomorrow":
		return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location()).Add(time.Hour * 24), nil
	default:
		weekday, err := parseWeekday(strings.ToLower(s))
		if err == nil {
			diff := daysToAdd(n.Weekday(), weekday)
			return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location()).Add(time.Hour * 24 * time.Duration(diff)), nil
		} else {
			switch len(s) {
			case 2:
				day, err := strconv.Atoi(s)
				if err != nil {
					return n, fmt.Errorf("couldn't parse day from %s", s)
				}
				t := time.Date(n.Year(), n.Month(), day, 0, 0, 0, 0, n.Location())
				if day < n.Day() {
					t = incMonth(t)
				}
				return t, nil

			case 4:
				day, err := strconv.Atoi(s[2:])
				if err != nil {
					return n, fmt.Errorf("couldn't parse day from %s", s[2:])
				}
				month, err := strconv.Atoi(s[:2])
				if err != nil {
					return n, fmt.Errorf("couldn't parse day from %s", s[:2])
				}
				year := n.Year()
				if month < int(n.Month()) || (month == int(n.Month()) && day < n.Day()) { //todo check time
					year++
				}
				return time.Date(year, time.Month(month), day, n.Hour(), n.Minute(), n.Second(), n.Nanosecond(), n.Location()), nil
			case 6:
				day, err := strconv.Atoi(s[4:])
				if err != nil {
					return n, fmt.Errorf("couldn't parse day from %s", s[4:])
				}
				month, err := strconv.Atoi(s[2:4])
				if err != nil {
					return n, fmt.Errorf("couldn't parse day from %s", s[2:4])
				}
				year, err := strconv.Atoi(s[:4])
				year += (n.Year() / 100) * 100
				if year >= n.Year() && month >= int(n.Month()) && day >= n.Day(){ //todo check time
					year += 100
				}
				return time.Date(year, time.Month(month), day, n.Hour(), n.Minute(), n.Second(), n.Nanosecond(), n.Location()), nil
			case 8: // will only work until year 9999
				day, err := strconv.Atoi(s[6:])
				if err != nil {
					return n, fmt.Errorf("couldn't parse day from %s", s[6:])
				}
				month, err := strconv.Atoi(s[4:6])
				if err != nil {
					return n, fmt.Errorf("couldn't parse day from %s", s[4:6])
				}
				year, err := strconv.Atoi(s[:4])
				return time.Date(year, time.Month(month), day, n.Hour(), n.Minute(), n.Second(), n.Nanosecond(), n.Location()), nil
			default:
				return n, fmt.Errorf("could not parse date from %s", s)
			}
		}
	}
}

func incMonth(t time.Time) time.Time {
	if t.Month() == 12 {
		return time.Date(t.Year() + 1, 1, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	} else {
		return time.Date(t.Year(), t.Month() + 1, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}
}

func daysToAdd (from, to time.Weekday) int {
	d := len(daysOfWeek)
	daysToAdd := (int(to) - int(from) + d) % d
	if daysToAdd == 0 {
		daysToAdd += d
	}
	return daysToAdd
}

var daysOfWeek = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

func parseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday format '%s'", v)
}