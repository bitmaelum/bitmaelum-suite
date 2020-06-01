package main

import "github.com/jaytaph/mailv2/client/cmd"

/*
[.]    account
[.]        add
[.]        list
[.]        update
[.]            pow <bits>
[.]            key
[X]                add <key>
[X]                list
[X]                rm <id>
[X]    resolve
[X]        ls <address>
[X]        update
[X]    mail
[X]        open    <account>
[X]        info    <account>
[X]        boxes   <account>
[X]        create	<to> <subject>
 */

func main() {
	cmd.Execute()
}
