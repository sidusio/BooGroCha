package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"sidus.io/boogrocha/internal/booking"
)

var listJSON bool

func ListCmd(getBS func() *booking.Directory) *cobra.Command {
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

func runList(cmd *cobra.Command, args []string, getBS func() *booking.Directory) {
	bookings, err := getBS().MyBookings()
	if err != nil {
		panic(err)
	}

	if !listJSON {
		fmt.Printf("%-7s %-13s %-15s %s\n", "DATE", "TIME", "ROOM", "TEXT")
		for _, booking := range bookings {
			date := booking.Start.Format("02/01")
			time := fmt.Sprintf("%s-%s", booking.Start.Format("15:04"), booking.End.Format("15:04"))
			text := fmt.Sprintf("\"%s\"", booking.Text)
			fmt.Printf("%-7s %-13s %-15s %s\n",
				date,
				time,
				booking.Room,
				text,
			)
		}
	} else {
		b, _ := json.Marshal(bookings)
		fmt.Println(string(b))
	}
}
