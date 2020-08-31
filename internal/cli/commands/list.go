package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sidus.io/boogrocha/internal/booking"
)

var listJSON bool

func ListCmd(getBS func() booking.BookingService) *cobra.Command {
	ListCmd := &cobra.Command{
		Use:   "list",
		Short: "List your current bookings",
		Long:  `Used to list all upcoming and current bookings`,
		Run: func(cmd *cobra.Command, args []string) {
			runList(cmd, args, getBS)
		},
	}

	ListCmd.Flags().BoolVarP(&listJSON, "json", "", false, "Formats the output in a machine readable format with id")

	return ListCmd
}

func runList(cmd *cobra.Command, args []string, getBS func() booking.BookingService) {
	bookings, err := getBS().MyBookings()
	if err != nil {
		fmt.Printf("Failed to get bookings: %v \n", err)
		os.Exit(1)
	}

	if !listJSON {
		fmt.Printf("%-7s %-13s %-15s %s\n", "DATE", "TIME", "ROOM", "TEXT")
		for _, booking := range bookings {
			date := formatDateWithWeekday(booking)
			time := formatTime(booking)
			text := fmt.Sprintf("\"%s\"", booking.Text)
			fmt.Printf("%-7s %-13s %-15s %s\n",
				date,
				time,
				booking.Room.Id,
				text,
			)
		}
	} else {
		b, _ := json.Marshal(bookings)
		fmt.Println(string(b))
	}
}

func formatDateWithWeekday(booking booking.Booking) (date string) {
	return booking.Start.Format("Mon 02/01")
}

func formatTime(booking booking.Booking) (time string) {
	return fmt.Sprintf("%s-%s", booking.Start.Format("15:04"), booking.End.Format("15:04"))
}
