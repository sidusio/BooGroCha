package commands

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sidus.io/boogrocha/internal/booking"
	"strconv"
	"strings"
)

var deleteAll bool

func DeleteCmd(getBS func() booking.BookingService) *cobra.Command {
	DeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a booking",
		Long:  `Used to delete a booking`,
		Run: func(cmd *cobra.Command, args []string) {
			runDelete(cmd, args, getBS)
		},
	}
	DeleteCmd.Flags().BoolVarP(&deleteAll, "all", "", false, "Unbooks all current bookings")
	DeleteCmd.AddCommand(deleteIdCmd)
	return DeleteCmd
}

func runDelete(cmd *cobra.Command, args []string, getBS func() booking.BookingService) {
	bs := getBS()
	bookings, err := bs.MyBookings()
	if err != nil {
		panic(err)
	}

	fmt.Printf("    %-7s %-13s %-15s %s\n", "DATE", "TIME", "ROOM", "TEXT")
	for i, booking := range bookings {
		date := booking.Start.Format("02/01")
		time := fmt.Sprintf("%s-%s", booking.Start.Format("15:04"), booking.End.Format("15:04"))
		text := fmt.Sprintf("\"%s\"", booking.Text)
		fmt.Printf("[%d] %-7s %-13s %-15s %s\n",
			i+1,
			date,
			time,
			booking.Room,
			text,
		)
	}

	fmt.Println("==> Booking to delete")
	fmt.Print("==> ")

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.Replace(input, "\n", "", -1)

	n, err := strconv.Atoi(input)
	n--
	if err != nil {
		fmt.Printf("invalid booking")
		os.Exit(1)
	}

	if (n) < len(bookings) && (n) >= 0 {
		fmt.Printf("Deleting booking %d...\n", n+1)
		err := bs.UnBook(bookings[n])
		if err != nil {
			fmt.Println("couldn't delete booking")
			os.Exit(1)
		}
		fmt.Println("Booking deleted successfully!")

	} else {
		print("no such booking")
	}

}

var deleteIdCmd = &cobra.Command{
	Use:   "delete id {id}",
	Short: "Delete ",
	Long:  "",
	Run:   nil,
}
