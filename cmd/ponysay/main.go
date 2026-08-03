package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ponysay-go/pkg/assets"
	"ponysay-go/pkg/balloon"
	"ponysay-go/pkg/color"
	"ponysay-go/pkg/pony"
	"ponysay-go/pkg/term"
	"ponysay-go/pkg/update"
)

var (
	version    = "v1.0.0"
	commitHash = "dev"
)

func getVersionString() string {
	if commitHash != "" && commitHash != "dev" {
		return fmt.Sprintf("ponysay-go %s", commitHash)
	}
	return fmt.Sprintf("ponysay-go %s", version)
}

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	execName := filepath.Base(os.Args[0])
	isThink := strings.HasSuffix(execName, "think")

	var ponyFiles stringSliceFlag
	var nonMLPPonies stringSliceFlag
	var anyPonies stringSliceFlag
	var quotePonies stringSliceFlag

	var balloonStyle string
	var wrapCol int

	var listMLP bool
	var listMLPAliases bool
	var listNonMLP bool
	var listNonMLPAliases bool
	var listAll bool
	var listAllAliases bool
	var oneListMLP bool
	var oneListNonMLP bool
	var oneListAll bool
	var listBalloons bool
	var listQuoters bool

	var compress bool
	var ponyOnly bool
	var infoLevel int // 0: off, 1: --info, 2: ++info
	var showVersion bool
	var showHelp bool
	var doUpdate bool

	var colorMode string // "256", "tty", "kms"
	var balloonColor string
	var linkColor string
	var msgColor string

	flagSet := flag.NewFlagSet("ponysay", flag.ExitOnError)

	flagSet.Var(&ponyFiles, "f", "Select a pony")
	flagSet.Var(&ponyFiles, "file", "Select a pony")
	flagSet.Var(&ponyFiles, "pony", "Select a pony")
	flagSet.Var(&ponyFiles, "files", "Select ponies")
	flagSet.Var(&ponyFiles, "ponies", "Select ponies")

	flagSet.Var(&quotePonies, "q", "Select a pony quote")
	flagSet.Var(&quotePonies, "quote", "Select a pony quote")
	flagSet.Var(&quotePonies, "quotes", "Select pony quotes")

	flagSet.StringVar(&balloonStyle, "b", "cowsay", "Select a balloon style")
	flagSet.StringVar(&balloonStyle, "bubble", "cowsay", "Select a balloon style")
	flagSet.StringVar(&balloonStyle, "balloon", "cowsay", "Select a balloon style")

	flagSet.IntVar(&wrapCol, "W", 0, "Wrap column width")
	flagSet.IntVar(&wrapCol, "wrap", 0, "Wrap column width")

	flagSet.BoolVar(&listMLP, "l", false, "List pony names")
	flagSet.BoolVar(&listMLP, "list", false, "List pony names")

	flagSet.BoolVar(&listMLPAliases, "L", false, "List pony names with alternatives")
	flagSet.BoolVar(&listMLPAliases, "symlist", false, "List pony names with alternatives")
	flagSet.BoolVar(&listMLPAliases, "altlist", false, "List pony names with alternatives")

	flagSet.BoolVar(&listAll, "A", false, "List all pony names")
	flagSet.BoolVar(&listAll, "all", false, "List all pony names")

	flagSet.BoolVar(&listAllAliases, "symall", false, "List all pony names with alternatives")
	flagSet.BoolVar(&listAllAliases, "altall", false, "List all pony names with alternatives")

	flagSet.BoolVar(&listBalloons, "B", false, "List balloon styles")
	flagSet.BoolVar(&listBalloons, "bubblelist", false, "List balloon styles")
	flagSet.BoolVar(&listBalloons, "balloonlist", false, "List balloon styles")

	flagSet.BoolVar(&listQuoters, "quoters", false, "List ponies with quotes")
	flagSet.BoolVar(&oneListMLP, "onelist", false, "List output in one line")

	flagSet.BoolVar(&compress, "c", false, "Compress messages")
	flagSet.BoolVar(&compress, "compress", false, "Compress messages")
	flagSet.BoolVar(&compress, "compact", false, "Compress messages")

	flagSet.BoolVar(&ponyOnly, "o", false, "Print only the pony")
	flagSet.BoolVar(&ponyOnly, "pony-only", false, "Print only the pony")
	flagSet.BoolVar(&ponyOnly, "ponyonly", false, "Print only the pony")

	var infoStandard bool
	flagSet.BoolVar(&infoStandard, "i", false, "Print metadata of pony")
	flagSet.BoolVar(&infoStandard, "info", false, "Print metadata of pony")

	var color256 bool
	flagSet.BoolVar(&color256, "X", false, "256 color mode")
	flagSet.BoolVar(&color256, "256-colours", false, "256 color mode")
	flagSet.BoolVar(&color256, "256colours", false, "256 color mode")

	var colorTTY bool
	flagSet.BoolVar(&colorTTY, "V", false, "TTY 16 color mode")
	flagSet.BoolVar(&colorTTY, "tty-colours", false, "TTY 16 color mode")
	flagSet.BoolVar(&colorTTY, "ttycolours", false, "TTY 16 color mode")

	var colorKMS bool
	flagSet.BoolVar(&colorKMS, "K", false, "KMS color mode")
	flagSet.BoolVar(&colorKMS, "kms-colours", false, "KMS color mode")

	flagSet.BoolVar(&showVersion, "v", false, "Print version")
	flagSet.BoolVar(&showVersion, "version", false, "Print version")

	flagSet.BoolVar(&showHelp, "h", false, "Print help message")
	flagSet.BoolVar(&showHelp, "help", false, "Print help message")

	flagSet.BoolVar(&doUpdate, "u", false, "Update to latest release from GitHub")
	flagSet.BoolVar(&doUpdate, "update", false, "Update to latest release from GitHub")

	flagSet.StringVar(&balloonColor, "colour-bubble", "", "Color of balloon border")
	flagSet.StringVar(&balloonColor, "colour-balloon", "", "Color of balloon border")
	flagSet.StringVar(&linkColor, "colour-link", "", "Color of link stem")
	flagSet.StringVar(&msgColor, "colour-msg", "", "Color of message text")
	flagSet.StringVar(&msgColor, "colour-message", "", "Color of message text")

	// Pre-process custom arguments (+f, +l, +L, +A, ++info, +c, --f, --q, update, etc.)
	args := os.Args[1:]
	var processedArgs []string
	var messageArgs []string
	quoteModeRequested := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "update" || arg == "-u" || arg == "--update" {
			doUpdate = true
		} else if arg == "-l" || arg == "--list" {
			listMLP = true
		} else if arg == "+l" || arg == "++list" {
			listNonMLP = true
		} else if arg == "-L" || arg == "--symlist" || arg == "--altlist" {
			listMLPAliases = true
		} else if arg == "+L" || arg == "++symlist" || arg == "++altlist" {
			listNonMLPAliases = true
		} else if arg == "-A" || arg == "--all" {
			listAll = true
		} else if arg == "+A" || arg == "++all" || arg == "++symall" || arg == "++altall" || arg == "--symall" || arg == "--altall" {
			listAllAliases = true
		} else if arg == "--onelist" {
			oneListMLP = true
		} else if arg == "++onelist" {
			oneListNonMLP = true
		} else if arg == "--Onelist" {
			oneListAll = true
		} else if arg == "-B" || arg == "--bubblelist" || arg == "--balloonlist" {
			listBalloons = true
		} else if arg == "--quoters" {
			listQuoters = true
		} else if arg == "-i" || arg == "--info" {
			infoLevel = 1
		} else if arg == "+i" || arg == "++info" {
			infoLevel = 2
		} else if strings.HasPrefix(arg, "+c") || strings.HasPrefix(arg, "--colour") {
			if strings.Contains(arg, "=") {
				balloonColor = arg[strings.Index(arg, "=")+1:]
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				balloonColor = args[i+1]
				i++
			}
		} else if arg == "-F" || arg == "+F" {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				anyPonies = append(anyPonies, args[i+1])
				i++
			}
		} else if arg == "--F" || arg == "++F" || arg == "--any-ponies" || arg == "++any-ponies" {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				anyPonies = append(anyPonies, args[i+1])
				i++
			}
		} else if strings.HasPrefix(arg, "-F") || strings.HasPrefix(arg, "+F") {
			val := arg[2:]
			if val != "" {
				anyPonies = append(anyPonies, val)
			}
		} else if arg == "+f" || arg == "++file" || arg == "++pony" {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				nonMLPPonies = append(nonMLPPonies, args[i+1])
				i++
			}
		} else if arg == "++f" || arg == "++files" || arg == "++ponies" {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				nonMLPPonies = append(nonMLPPonies, args[i+1])
				i++
			}
		} else if strings.HasPrefix(arg, "+f") {
			val := arg[2:]
			if val != "" {
				nonMLPPonies = append(nonMLPPonies, val)
			}
		} else if arg == "-f" || arg == "-file" || arg == "--file" || arg == "-pony" || arg == "--pony" {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				ponyFiles = append(ponyFiles, args[i+1])
				i++
			}
		} else if arg == "--f" || arg == "--files" || arg == "--ponies" {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				ponyFiles = append(ponyFiles, args[i+1])
				i++
			}
		} else if strings.HasPrefix(arg, "-f") {
			val := arg[2:]
			if val != "" {
				ponyFiles = append(ponyFiles, val)
			}
		} else if arg == "-q" || arg == "+q" || arg == "-quote" || arg == "--quote" {
			quoteModeRequested = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				quotePonies = append(quotePonies, args[i+1])
				i++
			}
		} else if arg == "--q" || arg == "--quotes" || arg == "++q" || arg == "++quotes" {
			quoteModeRequested = true
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				quotePonies = append(quotePonies, args[i+1])
				i++
			}
		} else if strings.HasPrefix(arg, "-q") || strings.HasPrefix(arg, "+q") {
			quoteModeRequested = true
			val := arg[2:]
			if val != "" {
				quotePonies = append(quotePonies, val)
			}
		} else if arg == "-b" || arg == "-bubble" || arg == "--bubble" || arg == "-balloon" || arg == "--balloon" || arg == "-W" || arg == "-wrap" || arg == "--wrap" {
			processedArgs = append(processedArgs, arg)
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !strings.HasPrefix(args[i+1], "+") {
				processedArgs = append(processedArgs, args[i+1])
				i++
			}
		} else if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "+") {
			processedArgs = append(processedArgs, arg)
		} else {
			messageArgs = append(messageArgs, arg)
		}
	}

	_ = flagSet.Parse(processedArgs)

	if doUpdate {
		if err := update.SelfUpdate(getVersionString(), update.DefaultRepo); err != nil {
			fmt.Fprintf(os.Stderr, "Update error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if color256 {
		colorMode = "256"
	} else if colorTTY {
		colorMode = "tty"
	} else if colorKMS {
		colorMode = "kms"
	}
	_ = colorMode

	if infoStandard && infoLevel == 0 {
		infoLevel = 1
	}

	if showHelp {
		printHelp(isThink)
		os.Exit(0)
	}

	if showVersion {
		fmt.Println(getVersionString())
		os.Exit(0)
	}

	am := assets.NewAssetManager()

	outputGroupedList := func(includeStandard, includeExtra, withAliases bool) {
		width := term.GetTerminalWidth()
		groups := am.GetPonyGroups(includeStandard, includeExtra, withAliases, true)
		for _, g := range groups {
			fmt.Println()
			fmt.Printf("\x1b[1mponies located in %s\x1b[0m\n", g.DirectoryPath)
			fmt.Print(assets.FormatColumnisedList(g.Ponies, width))
		}
	}

	outputList := func(items []string) {
		width := term.GetTerminalWidth()
		fmt.Print(assets.FormatColumnisedList(items, width))
	}

	if listQuoters {
		for _, item := range am.ListQuoters(true, false) {
			fmt.Println(item)
		}
		os.Exit(0)
	}

	if oneListAll || (oneListMLP && oneListNonMLP) {
		for _, item := range am.ListPoniesOneList(true, true) {
			fmt.Println(item)
		}
		os.Exit(0)
	}

	if oneListMLP {
		for _, item := range am.ListPoniesOneList(true, false) {
			fmt.Println(item)
		}
		os.Exit(0)
	}

	if oneListNonMLP {
		for _, item := range am.ListPoniesOneList(false, true) {
			fmt.Println(item)
		}
		os.Exit(0)
	}

	if listAllAliases || (listMLPAliases && listNonMLPAliases) {
		outputGroupedList(true, true, true)
		os.Exit(0)
	}

	if listAll || (listMLP && listNonMLP) {
		outputGroupedList(true, true, false)
		os.Exit(0)
	}

	if listMLPAliases {
		outputGroupedList(true, false, true)
		os.Exit(0)
	}

	if listMLP {
		outputGroupedList(true, false, false)
		os.Exit(0)
	}

	if listNonMLPAliases {
		outputGroupedList(false, true, true)
		os.Exit(0)
	}

	if listNonMLP {
		outputGroupedList(false, true, false)
		os.Exit(0)
	}

	if listBalloons {
		outputList(am.ListBalloons(isThink))
		os.Exit(0)
	}

	// Determine pony selection and quote handling
	selectedPony := ""
	includeStandard := true
	includeExtra := false

	if len(nonMLPPonies) > 0 {
		selectedPony = nonMLPPonies[rnd.Intn(len(nonMLPPonies))]
		includeStandard = false
		includeExtra = true
	} else if len(anyPonies) > 0 {
		selectedPony = anyPonies[rnd.Intn(len(anyPonies))]
		includeStandard = true
		includeExtra = true
	} else if len(ponyFiles) > 0 {
		selectedPony = ponyFiles[rnd.Intn(len(ponyFiles))]
		includeStandard = true
		includeExtra = false
	}

	var message string
	allMsgArgs := append(messageArgs, flagSet.Args()...)

	// Quote mode processing
	if quoteModeRequested || len(quotePonies) > 0 {
		qp := quotePonies
		if len(qp) == 0 && selectedPony != "" {
			qp = []string{selectedPony}
		}
		pName, qText, err := am.GetPonyQuote(qp)
		if err == nil {
			message = qText
			selectedPony = pName
			includeStandard = true
			includeExtra = true
		}
	}

	if message == "" {
		if len(allMsgArgs) > 0 {
			message = strings.Join(allMsgArgs, " ")
		} else {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				input, err := io.ReadAll(os.Stdin)
				if err == nil {
					message = strings.TrimRight(string(input), "\r\n")
				}
			}
		}
	}

	if message == "" && !ponyOnly {
		message = "I am just the cutest pony!"
	}

	if compress {
		lines := strings.Split(message, "\n")
		var clean []string
		for _, l := range lines {
			if strings.TrimSpace(l) != "" {
				clean = append(clean, l)
			}
		}
		message = strings.Join(clean, "\n")
	}

	if msgColor != "" {
		message = color.ApplyColor(message, msgColor)
	}

	// Load Pony file
	realPonyName, ponyContent, err := am.GetPonyFile(selectedPony, includeStandard, includeExtra)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p, err := pony.ParsePony(realPonyName, ponyContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing pony: %v\n", err)
		os.Exit(1)
	}

	// Info mode handling
	if infoLevel > 0 {
		fmt.Printf("NAME: %s\n", p.Name)
		for k, v := range p.Metadata {
			if infoLevel == 2 {
				fmt.Printf("\x1b[1m%s\x1b[0m: %s\n", k, v)
			} else {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
		os.Exit(0)
	}

	if ponyOnly {
		fmt.Println(p.RenderPonyOnly())
		os.Exit(0)
	}

	// Determine wrap column
	termWidth := term.GetTerminalWidth()
	if wrapCol <= 0 {
		wrapCol = termWidth - 10
		if wrapCol < 20 {
			wrapCol = 20
		}
	}

	// Load Balloon style
	balloonContent, err := am.GetBalloonContent(balloonStyle, isThink)
	if err != nil {
		balloonContent = ""
	}
	b := balloon.ParseBalloon(balloonContent, isThink)

	// Wrap and format message
	lines := balloon.WrapText(message, wrapCol)
	balloonLines := b.FormatBalloon(lines, 0, 0, balloonColor)

	// Render combined result
	rendered := p.RenderPonyWithBalloon(balloonLines, b.Link, linkColor)
	fmt.Println(rendered)
}

func printHelp(isThink bool) {
	cmdName := "ponysay"
	if isThink {
		cmdName = "ponythink"
	}
	fmt.Printf("%s - cowsay reimplementation for ponies (Go port)\n\n", cmdName)
	fmt.Println("Usage:")
	fmt.Printf("  %s [-f PONY] [-b STYLE] [-W COLUMN] [message]\n", cmdName)
	fmt.Printf("  %s -q [PONY]*\n", cmdName)
	fmt.Printf("  %s update | -u | --update\n", cmdName)
	fmt.Printf("  %s -l | -L | +l | +L | -A | +A | -B | --quoters | -i | -v | -h\n\n", cmdName)
	fmt.Println("Options:")
	fmt.Println("  -f, --file PONY    Select a pony by name or file.")
	fmt.Println("  +f PONY            Select a non-MLP pony.")
	fmt.Println("  -F PONY            Select any pony.")
	fmt.Println("  --f, --files       Variadic selection among ponies.")
	fmt.Println("  -q, --quote [PONY] Select a pony quote.")
	fmt.Println("  --q, --quotes      Variadic quote selection.")
	fmt.Println("  -b, --bubble STYLE Select balloon style (cowsay, unicode, ascii, round, etc.).")
	fmt.Println("  -W, --wrap COLUMN  Specify max wrapping width.")
	fmt.Println("  -l, --list         List pony names.")
	fmt.Println("  -L, --symlist      List pony names with alternative names.")
	fmt.Println("  +l                 List non-MLP pony names.")
	fmt.Println("  +L                 List non-MLP pony names with alternative names.")
	fmt.Println("  -A, --all          List all pony names.")
	fmt.Println("  +A                 List all pony names with alternative names.")
	fmt.Println("  -B, --bubblelist   List balloon styles.")
	fmt.Println("  --quoters          List ponies that have quotes.")
	fmt.Println("  --onelist          Format list output in a single line.")
	fmt.Println("  -c, --compress     Compress empty lines in message.")
	fmt.Println("  -i, --info         Print pony metadata info.")
	fmt.Println("  +i, ++info         Print pony metadata info formatted with colors.")
	fmt.Println("  -o, --pony-only    Print only the pony artwork.")
	fmt.Println("  -X, --256-colours  256 color mode.")
	fmt.Println("  -V, --tty-colours  TTY 16 color mode.")
	fmt.Println("  -K, --kms-colours  KMS color mode.")
	fmt.Println("  -u, --update       Update to the latest release from GitHub.")
	fmt.Println("  update             Update to the latest release from GitHub.")
	fmt.Println("  -v, --version      Print version information.")
	fmt.Println("  -h, --help         Print this help message.")
}
