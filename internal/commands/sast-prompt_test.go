package commands

import (
	"fmt"
	"testing"
)

const expectedOutputFormat = "<span style=\"color: regular;\">**CONFIDENCE:**</span><span style=\"color: grey; font-style: italic;\"> " +
	"A score between 0 (low) and 100 (high) indicating the degree of confidence in the exploitability of this vulnerability in the context of your code. " +
	"<br></span>%s<span style=\"color: regular;\">**EXPLANATION:**</span><span style=\"color: grey; font-style: italic;\"> " +
	"An OpenAI generated description of the vulnerability. <br></span>%s<span style=\"color: " +
	"regular;\">**PROPOSED REMEDIATION:**</span><span style=\"color: grey; font-style: italic;\"> " +
	"A customized snippet, generated by OpenAI, that can be used to remediate the vulnerability in your code. <br></span>%s"

func getExpectedOutput(confidenceNumber, explanationText, fixText string) string {
	return fmt.Sprintf(expectedOutputFormat, confidenceNumber, explanationText, fixText)
}

func TestAddDescriptionForIdentifiers(t *testing.T) {
	input := confidence + " 35 " + explanation + " this is a short explanation." + fix + " a fixed snippet"
	expected := getExpectedOutput(" 35 ", " this is a short explanation.", " a fixed snippet")
	output := getActual(input, t)

	if output[len(output)-1] != expected {
		t.Errorf("Expected %q, but got %q", expected, output)
	}
}

func TestAddNewlinesIfNecessarySomeNewlines(t *testing.T) {
	input := confidence + " 35 " + explanation + " this is a short explanation.\n" + fix + " a fixed snippet"
	expected := getExpectedOutput(" 35 ", " this is a short explanation.\n", " a fixed snippet")

	output := getActual(input, t)

	if output[len(output)-1] != expected {
		t.Errorf("Expected %q, but got %q", expected, output)
	}
}

func TestAddNewlinesIfNecessaryAllNewlines(t *testing.T) {
	input := confidence + " 35\n " + explanation + " this is a short explanation.\n" + fix + " a fixed snippet"
	expected := getExpectedOutput(" 35\n ", " this is a short explanation.\n", " a fixed snippet")

	output := getActual(input, t)

	if output[len(output)-1] != expected {
		t.Errorf("Expected %q, but got %q", expected, output)
	}
}

func getActual(input string, t *testing.T) []string {
	someText := "some text"
	response := []string{someText, someText, input}
	output := addDescriptionForIdentifier(response)
	for i := 0; i < len(output)-1; i++ {
		if output[i] != response[i] {
			t.Errorf("All strings except last expected to stay the same")
		}
	}
	return output
}
