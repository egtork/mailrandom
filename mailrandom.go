package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

type mailConfig struct {
	Server   string
	Port     int
	To       string
	From     string
	Username string
}

type option struct {
	label            string
	weight           float64
	cumulativeWeight float64
}

const (
	passwordEnvVar = "MAIL_PASS"
)

var (
	category       = flag.String("c", "selection", "category, e.g. \"exercise\", \"study topic\", etc.")
	optionsInput   = flag.String("o", "", "comma-separated list of options")
	weightsInput   = flag.String("w", "", "comma-separated list of weights")
	inputFile      = flag.String("i", "", "weighted option CSV file with one option and one weight per line; overrides -o and -w flags")
	mailConfigFile = flag.String("m", "", "mail configuration JSON file; use -mail-help flag for more information")
	authPassword   = flag.String("p", "", "mail server password; overrides MAIL_PASS environment variable")
	mailHelp       = flag.Bool("mail-help", false, "information on mail configuration")
)

// parseMailConfig parses a mail configuration JSON file.
func parseMailConfig(filename string) *mailConfig {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fatal(err)
	}
	m := new(mailConfig)
	err = json.Unmarshal(b, m)
	if err != nil {
		fatal(err)
	}
	fmt.Println(m)
	return m
}

// parseOptionsFile parses a CSV file containing a list of weighted options. Each line should
// contain an option and a weight separated by a comma.
func parseOptionsFile(f *os.File) []*option {
	var options []*option
	var weightSum float64

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fatal(err)
		}
		weight, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			fatal(err)
		}
		weightSum += weight
		x := &option{
			label:            record[0],
			weight:           weight,
			cumulativeWeight: weightSum,
		}
		options = append(options, x)
	}
	if len(options) == 0 {
		log.Fatal("No options defined in input file")
	}
	return options
}

// parseCommandLineOptions parses a list of comma separated options and (optionally) a list
// of comma separated weights. If weights are not provided, the options are equally weighted.
func parseCommandLineOptions() []*option {
	labels := strings.Split(*optionsInput, ",")
	weights := strings.Split(*weightsInput, ",")
	weightsSpecified := len(weights) > 0
	if weightsSpecified && len(labels) != len(weights) {
		panic("number of options should equal number of weights")
	}
	var o *option
	var sum float64
	options := make([]*option, len(labels))
	for k := 0; k < len(labels); k++ {
		if weightsSpecified {
			w, err := strconv.ParseFloat(weights[k], 64)
			sum += w
			if err != nil {
				panic(err)
			}
			o = &option{label: labels[k], weight: w, cumulativeWeight: sum}
		} else {
			o = &option{label: labels[k], weight: 1, cumulativeWeight: float64(k)}
		}
		options[k] = o
	}
	return options
}

// selectOption selections an option from a slice of weighted options.
func selectOption(options []*option) string {
	sum := options[len(options)-1].cumulativeWeight
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := sum * r.Float64()

	for k := 0; k <= len(options); k++ {
		if x < options[k].cumulativeWeight {
			return options[k].label
		}
	}
	return ""
}

// sendMail sends a selected option as an email using SMTP configuration options in cfg.
// If cfg.Username is set, authentication will be attempted using password supplied as
// a command line option or set as the MAIL_PASS environment variable.
func sendMail(cfg *mailConfig, selection string) error {
	var auth smtp.Auth
	if cfg.Username != "" {
		if *authPassword == "" {
			pass, ok := os.LookupEnv(passwordEnvVar)
			if !ok {
				fatal(errors.New("set SMTP authentication password using -password flag or " + passwordEnvVar + " environment variable"))
			}
			authPassword = &pass
		}
		auth = smtp.PlainAuth("", cfg.Username, *authPassword, cfg.Server)
	}
	msg := fmt.Sprintf("From: %s\r\n", cfg.From) +
		fmt.Sprintf("To: %s\r\n", cfg.To) +
		fmt.Sprintf("Subject: Today's %s is: %s\r\n\r\n", *category, selection) +
		"\r\n"

	server := cfg.Server + ":" + strconv.Itoa(cfg.Port)
	fmt.Println(server)
	err := smtp.SendMail(server,
		auth,
		cfg.From, []string{cfg.To}, []byte(msg))
	return err
}

// displayMailHelp displays the expected format of the mail configuration JSON file.
func displayMailHelp() {
	fmt.Println("The format of the mail configuration JSON file is as follows:")
	fmt.Println(`  {
    "Server":   "smtp.example-server.com",
    "Port":     587,
    "Username": "myself@example-server.com",
    "To":       "myself@example-server.com",
    "From":     "myself@example-server.com"
  }`)
	fmt.Println("Username can be omitted if authentication is not required.")
}

func main() {
	flag.Parse()
	if *mailHelp {
		displayMailHelp()
		os.Exit(0)
	}
	var options []*option
	if *inputFile == "" {
		if *optionsInput == "" {
			flag.Usage()
			os.Exit(0)
		} else {
			// Get weighted options from command line
			options = parseCommandLineOptions()
		}
	} else {
		// Get weighted options from CSV file
		f, err := os.Open(*inputFile)
		if err != nil {
			log.Fatal(err)
		}
		options = parseOptionsFile(f)
	}

	if len(options) == 0 {
		fmt.Println("No options specified")
		os.Exit(1)
	}
	selection := selectOption(options)

	if *mailConfigFile == "" {
		// Print selection to console
		fmt.Println("Today's " + *category + " is: " + selection)
	} else {
		// Send selection by mail
		cfg := parseMailConfig(*mailConfigFile)
		err := sendMail(cfg, selection)
		if err != nil {
			log.Println(err)
		}
	}
}

func fatal(s error) {
	fmt.Fprintln(os.Stderr, s.Error())
	os.Exit(1)
}
