package gosshtun

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptForPassword ask
func PromptForPassword(username string) (pw string, err error) {
	start := getNewPasswordStarter()

	end := ""
	const numTry = 3
	for i := 0; i < numTry; i++ {
		switch i {
		case 1:
			fmt.Printf("\n%s\n... no problem, try it again (attempt %v of %v):\n\n", err, i+1, numTry)
		case numTry - 1:
			fmt.Printf("\n%s\n... arg, still not right. One last try:\n\n", err)
		}
		fmt.Printf("adding user '%s'...\n\nThe first part of your new passphrase is '%s'. Add a memorable end to the sentence (between 3 - 100 characters) to complete it. For a strong passphrase, add five(5) or more words on top of the three we start you with\n\n%s",
			username, start, start)
		reader := bufio.NewReader(os.Stdin)
		end, err = reader.ReadString('\n')
		panicOn(err)
		end = strings.Trim(end, "\n\r")
		if len(end) < 3 {
			err = fmt.Errorf("alert! completion of phrase too short, must be 3-100 characters")
			continue
		}
		if len(end) > 100 {
			err = fmt.Errorf("alert! completion of phrase too long, must be 3-100 characters")
			continue
		}
		pw = start + end
		fmt.Printf("\n Your new passphrase for account '%s' is '%s' (without the single quotes). Type the phrase in full to confirm you have it\n\n", username, pw)

		var pw2 string
		pw2, err = reader.ReadString('\n')
		panicOn(err)
		pw2 = strings.Trim(pw2, "\n\r")
		if pw != pw2 {
			err = fmt.Errorf("alert! passphrases don't match at position %v: '%s' versus '%s'", firstDiff(pw, pw2)+1, pw, pw2)
			continue
		}

		fmt.Printf("\n Phrases match! Success!\n")
		// success
		err = nil
		return
	}
	// fail
	pw = ""
	return
}

func firstDiff(a, b string) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	i := 0
	for ; i < len(a); i++ {
		if i >= len(b) {
			return i
		}
		if a[i] != b[i] {
			return i
		}
	}
	return i
}
