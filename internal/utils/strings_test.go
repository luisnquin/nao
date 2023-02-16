package utils_test

import (
	"testing"

	"github.com/luisnquin/nao/v3/internal/utils"
)

func TestToPascalCase(t *testing.T) {
	checks := []struct {
		in, out string
	}{
		{
			in:  "Hello world",
			out: "HelloWorld",
		},
		{
			in:  "pascal case",
			out: "PascalCase",
		},
		{
			in:  "What do you think of the cookies?",
			out: "WhatDoYouThinkOfTheCookies?",
		},
		{
			in:  "space odissey",
			out: "SpaceOdissey",
		},
		{
			in:  "you can't see me anymore!",
			out: "YouCan'tSeeMeAnymore!",
		},
	}

	for _, expected := range checks {
		if out := utils.ToPascalCase(expected.in); out != expected.out {
			t.Errorf("expected '%s', but got '%s' from '%s'", expected.out, out, expected.in)
			t.Fail()
		}
	}
}

func TestToCamelCase(t *testing.T) {
	checks := []struct {
		in, out string
	}{
		{
			in:  "Hello world",
			out: "helloWorld",
		},
		{
			in:  "abc de",
			out: "abcDe",
		},
		{
			in:  "your birthday in two days",
			out: "yourBirthdayInTwoDays",
		},
		{
			in:  "nice spread",
			out: "niceSpread",
		},
		{
			in:  "the poetry where you practically said wait for yesterday for tomorrow",
			out: "thePoetryWhereYouPracticallySaidWaitForYesterdayForTomorrow",
		},
	}

	for _, expected := range checks {
		if out := utils.ToCamelCase(expected.in); out != expected.out {
			t.Errorf("expected '%s', but got '%s' from '%s'", expected.out, out, expected.in)
			t.Fail()
		}
	}
}
