package core

import "github.com/gookit/color"

var asciiLogo = " ____  _ _   __  __            _\n" +
		 "|  _ \\(_) | |  \\/  |          | |\n" +
		 "| |_) |_| |_| \\  / | __ _  ___| |_   _ _ __ ___\n" +
		 "|  _ <| | __| |\\/| |/ _` |/ _ \\ | | | | '_ ` _ \\\n" +
		 "| |_) | | |_| |  | | (_| |  __/ | |_| | | | | | |\n" +
		 "|____/|_|\\__|_|  |_|\\__,_|\\___|_|\\__,_|_| |_| |_|\n" +
		 "   P r i v a c y   i s   y o u r s   a g a i n\n"

var rainbowAsciiLogo = "\033[31m ____  _ _   __  __            _\n" +
	"\033[32m|  _ \\(_) | |  \\/  |          | |\n" +
	"\033[33m| |_) |_| |_| \\  / | __ _  ___| |_   _ _ __ ___\n" +
	"\033[34m|  _ <| | __| |\\/| |/ _` |/ _ \\ | | | | '_ ` _ \\\n" +
	"\033[35m| |_) | | |_| |  | | (_| |  __/ | |_| | | | | | |\n" +
	"\033[36m|____/|_|\\__|_|  |_|\\__,_|\\___|_|\\__,_|_| |_| |_|\n" +
	"\033[37m   P r i v a c y   i s   y o u r s   a g a i n\n" +
	"\033[0m"

var rainbow256AsciiLogo = "\033[38;5;208m ____  _ _   __  __            _\n" +
	"\033[38;5;209m|  _ \\(_) | |  \\/  |          | |\n" +
	"\033[38;5;210m| |_) |_| |_| \\  / | __ _  ___| |_   _ _ __ ___\n" +
	"\033[38;5;211m|  _ <| | __| |\\/| |/ _` |/ _ \\ | | | | '_ ` _ \\\n" +
	"\033[38;5;212m| |_) | | |_| |  | | (_| |  __/ | |_| | | | | | |\n" +
	"\033[38;5;213m|____/|_|\\__|_|  |_|\\__,_|\\___|_|\\__,_|_| |_| |_|\n" +
	"\033[38;5;214m   P r i v a c y   i s   y o u r s   a g a i n\n" +
	"\033[0m"

func GetMonochromeAsciiLogo() string {
	return asciiLogo
}

// Returns ASCII logo with or without colors depending on your console settings
func GetAsciiLogo() string {
	// Ooh. Nice and shiny terminal! Display a cool colorscheme
	if color.IsSupport256Color() || color.IsSupportTrueColor() {
		return rainbow256AsciiLogo
	}

	// A fresh rainbow color made out of the standard 16 ANSI colors
	if color.IsSupportColor() {
		return rainbowAsciiLogo
	}

	// No color. Lame :/
	return asciiLogo

}
