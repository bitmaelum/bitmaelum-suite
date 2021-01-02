// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package internal

import (
	"os"
	"runtime"

	"github.com/gookit/color"
)

var asciiLogo = " ____  _ _   __  __            _\n" +
	"|  _ \\(_) | |  \\/  |          | |\n" +
	"| |_) |_| |_| \\  / | __ _  ___| |_   _ _ __ ___\n" +
	"|  _ <| | __| |\\/| |/ _` |/ _ \\ | | | | '_ ` _ \\\n" +
	"| |_) | | |_| |  | | (_| |  __/ | |_| | | | | | |\n" +
	"|____/|_|\\__|_|  |_|\\__,_|\\___|_|\\__,_|_| |_| |_|\n" +
	"\n" +
	"   P r i v a c y   i s   y o u r s   a g a i n\n"

var rainbowASCIILogo = "\033[31m ____  _ _   __  __            _\n" +
	"\033[32m|  _ \\(_) | |  \\/  |          | |\n" +
	"\033[33m| |_) |_| |_| \\  / | __ _  ___| |_   _ _ __ ___\n" +
	"\033[34m|  _ <| | __| |\\/| |/ _` |/ _ \\ | | | | '_ ` _ \\\n" +
	"\033[35m| |_) | | |_| |  | | (_| |  __/ | |_| | | | | | |\n" +
	"\033[36m|____/|_|\\__|_|  |_|\\__,_|\\___|_|\\__,_|_| |_| |_|\n" +
	"\n" +
	"\033[37m   P r i v a c y   i s   y o u r s   a g a i n\n" +
	"\033[0m"

var rainbow256ASCIILogo = "\033[38;5;208m ____  _ _   __  __            _\n" +
	"\033[38;5;209m|  _ \\(_) | |  \\/  |          | |\n" +
	"\033[38;5;210m| |_) |_| |_| \\  / | __ _  ___| |_   _ _ __ ___\n" +
	"\033[38;5;211m|  _ <| | __| |\\/| |/ _` |/ _ \\ | | | | '_ ` _ \\\n" +
	"\033[38;5;212m| |_) | | |_| |  | | (_| |  __/ | |_| | | | | | |\n" +
	"\033[38;5;213m|____/|_|\\__|_|  |_|\\__,_|\\___|_|\\__,_|_| |_| |_|\n" +
	"\n" +
	"\033[38;5;214m   P r i v a c y   i s   y o u r s   a g a i n\n" +
	"\033[0m"

// GetMonochromeASCIILogo returns the monochrome version of the logo
func GetMonochromeASCIILogo() string {
	return asciiLogo
}

// GetASCIILogo returns ASCII logo with or without colors depending on your console settings
func GetASCIILogo() string {
	// Ooh. Nice and shiny terminal! Display a cool colorscheme
	if color.IsSupport256Color() || color.IsSupportTrueColor() || isPowerShell() {
		return rainbow256ASCIILogo
	}

	// A fresh rainbow color made out of the standard 16 ANSI colors
	if color.IsSupportColor() {
		return rainbowASCIILogo
	}

	// No color. Lame :/
	return asciiLogo
}

// isPowerShell returns true when detecting Windows/Powershell.. All bets are off with terminal/console I guess
func isPowerShell() bool {
	return os.Getenv("WT_SESSION") != "" && runtime.GOOS == "windows"
}
