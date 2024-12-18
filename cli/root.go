package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportTimestamp: true,
	TimeFormat:      time.TimeOnly,
	Prefix:          "🎄 aoc",
})

var VERSION string

func setLogLevel(cmd *cobra.Command) {
	quiet, err := cmd.Flags().GetBool("quiet")

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	if quiet {
		logger.SetLevel(log.WarnLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
}

var rootCmd = &cobra.Command{
	Use:   "aoc",
	Short: "\n🎄🎄🎄 Advent of Code 🎄🎄🎄\n\nAoC command-line tool",
	Run: func(cmd *cobra.Command, args []string) {

		version, err := cmd.Flags().GetBool("version")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if version {
			fmt.Println(VERSION)
			os.Exit(0)
		}

		setLogLevel(cmd)

		// If no command specified, run in interactive mode
		Interactive()
	},
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Create a template folder for a specific day",
	Example: "aoc init --day 1",
	Run: func(cmd *cobra.Command, args []string) {
		setLogLevel(cmd)

		year, err := validateYearFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		validateDayFlag(cmd)

		day, err := validateDayFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		InitialiseDay(fmt.Sprint(year), fmt.Sprint(day))
	},
}

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download puzzle inputs for specific year/day",
	Example: "aoc download --day 1",
	Run: func(cmd *cobra.Command, args []string) {
		setLogLevel(cmd)

		year, err := validateYearFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		day, err := validateDayFlag(cmd)

		if day == -1 {
			fmt.Println(err)
			os.Exit(1)
		}

		DownloadInput(fmt.Sprint(year), fmt.Sprint(day))
	},
}

var solveCmd = &cobra.Command{
	Use:     "solve",
	Short:   "Run the solution for a specific day",
	Example: "aoc solve --day 1 --part 1",
	Run: func(cmd *cobra.Command, args []string) {
		setLogLevel(cmd)

		year, err := validateYearFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		day, err := validateDayFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		part, err := validatePartFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		example, err := cmd.Flags().GetBool("example")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		SolveDay(fmt.Sprint(year), fmt.Sprint(day), fmt.Sprint(part), example)
	},
}

var submitCmd = &cobra.Command{
	Use:     "submit",
	Short:   "Submit an answer for a specific day",
	Example: "aoc submit --day 1 --part 1",
	Run: func(cmd *cobra.Command, args []string) {
		setLogLevel(cmd)

		year, err := validateYearFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		day, err := validateDayFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		part, err := validatePartFlag(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		SubmitAnswer(fmt.Sprint(year), fmt.Sprint(day), fmt.Sprint(part))
	},
}

func validateYearFlag(cmd *cobra.Command) (int, error) {
	year, err := cmd.Flags().GetInt("year")
	if err != nil {
		return 0, err
	}

	if year < 2015 {
		return 0, fmt.Errorf("error: The 'year' flag must be greater than 2015")
	}

	return year, nil
}

func validateDayFlag(cmd *cobra.Command) (int, error) {
	day, err := cmd.Flags().GetInt("day")
	if err != nil {
		return -1, err
	}

	if day < 1 || day > 25 {
		if day == 0 {
			return 0, fmt.Errorf("error: The 'day' flag must be greater than 0")
		}
		return -1, fmt.Errorf("error: The 'day' flag must be between 1 and 25")
	}

	return day, nil
}

func validatePartFlag(cmd *cobra.Command) (int, error) {
	part, err := cmd.Flags().GetInt("part")
	if err != nil {
		return 0, err
	}

	if part < 1 || part > 2 {
		return 0, fmt.Errorf("error: The 'part' flag must be either 1 or 2")
	}

	return part, nil
}

func init() {

	currentYear, currentMonth, currentDay := time.Now().Date()

	defaultDay := 0
	defaultYear := currentYear

	if currentMonth == 12 {
		defaultDay = currentDay
	} else {
		defaultYear = defaultYear - 1
	}

	rootCmd.Flags().BoolP("version", "v", false, "show version")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet mode")

	initCmd.Flags().IntP("year", "y", defaultYear, "puzzle year")
	initCmd.Flags().IntP("day", "d", defaultDay, "puzzle day")
	initCmd.MarkFlagRequired("day")

	solveCmd.Flags().IntP("year", "y", defaultYear, "puzzle year")
	solveCmd.Flags().IntP("day", "d", defaultDay, "puzzle day")
	solveCmd.Flags().IntP("part", "p", 1, "puzzle part")
	solveCmd.Flags().BoolP("example", "e", false, "use example input")
	solveCmd.MarkFlagRequired("day")

	submitCmd.Flags().IntP("year", "y", defaultYear, "puzzle year")
	submitCmd.Flags().IntP("day", "d", defaultDay, "puzzle day")
	submitCmd.Flags().IntP("part", "p", 1, "puzzle part")
	submitCmd.MarkFlagRequired("day")

	downloadCmd.Flags().IntP("year", "y", defaultYear, "puzzle year")
	downloadCmd.Flags().IntP("day", "d", 0, "puzzle day")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(solveCmd)
	rootCmd.AddCommand(submitCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
