package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type cup int

func (c cup) String() string {
	s := fmt.Sprintf("%d cup", c)
	if c > 1 {
		s += "s"
	}
	return s
}

func startOven(degree int) {
	log.Helper()
	log.Fatal("Starting oven", "temperature", degree)
}

func main() {

	lgr := log.New(os.Stdout)

	lgr.SetFormatter(log.JSONFormatter)
	lgr.SetLevel(log.DebugLevel)

	lgr.Debug("SOME SHIT")

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "Baking 🍪 ",
	})

	// Override the default error level style.
	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR!!").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))
	// Add a custom style for key `err`
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)
	logger = log.New(os.Stderr)
	logger.SetStyles(styles)
	logger.Error("Whoops!", "err", "kitchen on fire")

	logger.Info("Starting oven!", "degree", 375)
	time.Sleep(1 * time.Second)
	logger.Info("Finished baking")

	log.SetTimeFormat(time.Kitchen)
	log.SetLevel(log.DebugLevel)

	var (
		butter    = cup(1)
		chocolate = cup(2)
		flour     = cup(3)
		sugar     = cup(5)
		temp      = 375
		bakeTime  = 10
	)

	startOven(temp)
	time.Sleep(time.Second)
	log.Debug("Mixing ingredients", "ingredients",
		strings.Join([]string{
			butter.String() + " of butter",
			chocolate.String() + " of chocolate",
			flour.String() + " of flour",
			sugar.String() + " of sugar",
		}, "\n"),
	)
	time.Sleep(time.Second)
	if sugar > 2 {
		log.Warn("That's a lot of sugar", "amount", sugar)
	}
	log.Info("Baking cookies", "time", fmt.Sprintf("%d minutes", bakeTime))
	time.Sleep(2 * time.Second)
	log.Info("Increasing temperature", "amount", 300)
	temp += 300
	time.Sleep(time.Second)
	if temp > 500 {
		log.Error("Oven is too hot", "temperature", temp)
		log.Fatal("The kitchen is on fire 🔥")
	}
}
