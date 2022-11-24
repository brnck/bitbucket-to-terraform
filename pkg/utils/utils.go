package utils

import "strings"

type replacementList struct {
	word        string
	replaceWith string
}

var changeList = []replacementList{
	{"0", "zero_"},
	{"1", "one_"},
	{"2", "two_"},
	{"3", "three_"},
	{"4", "four_"},
	{"5", "five_"},
	{"6", "six_"},
	{"7", "seven_"},
	{"8", "eight_"},
	{"9", "nine_"},
}

// TransformStringToBeTFCompliant takes string and transforms it
// in order to be compliant with TF resource name requirements
func TransformStringToBeTFCompliant(n string) string {
	n = strings.TrimSpace(n)
	n = strings.ToLower(n)
	n = strings.ReplaceAll(n, " ", "_")
	n = strings.ReplaceAll(n, "-", "_")
	n = strings.ReplaceAll(n, ".", "_")

	for _, v := range changeList {
		if n[0:1] != v.word {
			continue
		}

		n = v.replaceWith + n[1:]
	}

	return n
}

func StringExistsInList(list []string, item string) bool {
	for _, listItem := range list {
		if item == listItem {
			return true
		}
	}

	return false
}
