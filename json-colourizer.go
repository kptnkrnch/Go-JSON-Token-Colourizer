package main

import (
	"fmt"
	"os"
	"bufio"
	"unicode"
)

func main() {
	if len(os.Args) >= 2 {
		tokens := scanner(os.Args[1])
		colorizer(tokens)
	} else {
		fmt.Println("Missing the name of a valid JSON file.")
	}
}

const defaultExitToken string = ""

// Various JSON token types
const TypeString string = "string"
const TypeBracket string = "bracket"
const TypeColon string = "colon"
const TypeComma string = "comma"
const TypeSquareBracket string = "squarebracket"
const TypeNumber string = "number"
const TypeKeyword string = "keyword"

// Structure for storing token type and its data
type Token struct {
	tokenType string
	data string
	color string // unused
}
func (t Token) String() string {
	return t.data
}

// Scanner function that reads tokens from the json file (character by character)
func scanner(filename string) []Token {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	var tokenBuffer string
	var exitToken string
	var tokens []Token

	for scanner.Scan() {
		tokenBuffer = ""
		var token Token
		character := scanner.Text()
		exitToken = defaultExitToken
		
		// checking what type of token a character belongs too
		if isSquareBracket(character) {
			tokenBuffer = character
			token.tokenType = TypeSquareBracket
			token.data = tokenBuffer
			token.color = ""
			tokens = append(tokens, token)
		} else if isColon(character) {
			tokenBuffer = character
			token.tokenType = TypeColon
			token.data = tokenBuffer
			token.color = ""
			tokens = append(tokens, token)
		} else if isComma(character) {
			tokenBuffer = character
			token.tokenType = TypeComma
			token.data = tokenBuffer
			token.color = ""
			tokens = append(tokens, token)
		} else if isStringStart(character) {
			tokenBuffer, exitToken = readString(character, scanner)
			token.tokenType = TypeString
			token.data = tokenBuffer
			token.color = ""
			tokens = append(tokens, token)
		} else if isBracket(character) {
			tokenBuffer = character
			token.tokenType = TypeBracket
			token.data = tokenBuffer
			token.color = ""
			tokens = append(tokens, token)
		} else if isNumberStart(character) {
			tokenBuffer, exitToken = readNumber(character, scanner)
			token.tokenType = TypeNumber
			token.data = tokenBuffer
			token.color = ""
			tokens = append(tokens, token)
		} else if isKeywordStart(character) {
			tokenBuffer, _ = readKeyword(character, scanner)
			token.tokenType = TypeKeyword
			token.data = tokenBuffer
			token.color = ""
			tokens = append(tokens, token)
		}

		// loop to keep on reading tokens if the last token read is not whitespace
		for exitToken != defaultExitToken {
			if isSquareBracket(exitToken) {
				tokenBuffer = exitToken
				token.tokenType = TypeSquareBracket
				token.data = tokenBuffer
				token.color = ""
				tokens = append(tokens, token)
				exitToken = defaultExitToken
			} else if isColon(exitToken) {
				tokenBuffer = exitToken
				token.tokenType = TypeColon
				token.data = tokenBuffer
				token.color = ""
				tokens = append(tokens, token)
				exitToken = defaultExitToken
			} else if isComma(exitToken) {
				tokenBuffer = exitToken
				token.tokenType = TypeComma
				token.data = tokenBuffer
				token.color = ""
				tokens = append(tokens, token)
				exitToken = defaultExitToken
			} else if isStringStart(exitToken) {
				tokenBuffer, exitToken = readString(exitToken, scanner)
				token.tokenType = TypeString
				token.data = tokenBuffer
				token.color = ""
				tokens = append(tokens, token)
			} else if isBracket(exitToken) {
				tokenBuffer = exitToken
				token.tokenType = TypeBracket
				token.data = tokenBuffer
				token.color = ""
				tokens = append(tokens, token)
				exitToken = defaultExitToken
			} else if isNumberStart(exitToken) {
				tokenBuffer, exitToken = readNumber(exitToken, scanner)
				token.tokenType = TypeNumber
				token.data = tokenBuffer
				token.color = ""
				tokens = append(tokens, token)
			} else if isKeywordStart(exitToken) {
				tokenBuffer, _ = readKeyword(exitToken, scanner)
				token.tokenType = TypeKeyword
				token.data = tokenBuffer
				token.color = ""
				tokens = append(tokens, token)
				exitToken = defaultExitToken
			} else if exitToken == " " {
				exitToken = defaultExitToken
			} else {
				exitToken = defaultExitToken
			}
		}
	}
	file.Close()

	return tokens
}

// checks if a character is a square bracket token
func isSquareBracket(input string) bool {
	if input == "[" || input == "]" {
		return true
	} else {
		return false
	}
}

// checks if a character is a curly bracket token
func isBracket(input string) bool {
	if input == "{" || input == "}" {
		return true
	} else {
		return false
	}
}

// checks if a character is a colon token
func isColon(input string) bool {
	if input == ":" {
		return true
	} else {
		return false
	}
}

// checks if a character is a comma token
func isComma(input string) bool {
	if input == "," {
		return true
	} else {
		return false
	}
}

// checks if a character the beginning of a number token
func isNumberStart(input string) bool {
	if input == "-" || unicode.IsDigit(rune(input[0])) {
		return true
	} else {
		return false
	}
}

// reads subsequent characters until a character that violates the JSON number specification is found
func readNumber(initialCharacter string, scanner *bufio.Scanner) (string, string) {
	var number string = initialCharacter

	var encounteredDot bool = false
	var secondCharacter string
	var encounteredSecond bool = false
	var expectingDot bool = false
	var expectingSign bool = false
	var encounteredSign bool = false

	for scanner.Scan() {
		token := scanner.Text()
		if !encounteredSecond {
			secondCharacter = token
			encounteredSecond = true
		}

		if !encounteredDot && !expectingSign && !encounteredSign && !expectingDot && initialCharacter == "0" && secondCharacter == "." {
			number += token
			encounteredDot = true
		} else if !encounteredDot && !expectingSign && !encounteredSign && !expectingDot && initialCharacter == "-" && secondCharacter == "0" {
			number += token
			expectingDot = true
		} else if unicode.IsDigit(rune(token[0])) && !expectingDot && !encounteredSign {
			number += token
		} else if token == "." && !encounteredDot && !encounteredSign {
			number += token
			encounteredDot = true
		} else if encounteredDot && !encounteredSign && !expectingSign && unicode.IsDigit(rune(token[0])) {
			number += token
		} else if !encounteredSign && !expectingSign && token == "e" || token == "E" {
			number += token
			expectingSign = true
		} else if !encounteredSign && expectingSign && (token == "-" || token == "+" || unicode.IsDigit(rune(token[0]))) {
			number += token
			encounteredSign = true
		} else if encounteredSign && unicode.IsDigit(rune(token[0])) {
			number += token
		} else {
			return number, token
		}
	}
	return number, defaultExitToken
}

// reads string characters until an ending quote is found (not an escaped quote)
func readString(initialCharacter string, scanner *bufio.Scanner) (string, string) {
	var stringData string = initialCharacter
	var prevCharacter string = initialCharacter
	for scanner.Scan() {
		token := scanner.Text()
		if token[0] != '"' {
			stringData += token
		} else if token[0] == '"' && len(prevCharacter) > 0 && prevCharacter == "\\" {
			stringData += token
		} else {
			stringData += token
			if token[0] == '"' {
				token = defaultExitToken
			}
			return stringData, token
		}
		if len(prevCharacter) > 0 && prevCharacter == "\\" && token == "\\" {
			prevCharacter = ""
		} else {
			prevCharacter = token
		}
	}
	return stringData, defaultExitToken
}

// checks if the character is a quote that marks the start of a string
func isStringStart(input string) bool {
	if input == "\"" {
		return true
	} else {
		return false
	}
}

// checks if the character is the start of a keyword (true, false, null)
func isKeywordStart(input string) bool {
	if input == "t" || input == "n" || input == "f" {
		return true
	} else {
		return false
	}
}

// reads the rest of the three keywords (true, false, null)
func readKeyword(initialCharacter string, scanner *bufio.Scanner) (string, bool) {
	var keyword string = initialCharacter
	var validKeyword bool = true

	for scanner.Scan() {
		token := scanner.Text()

		if initialCharacter == "t" {
			keyword += token
			if len(keyword) >= 4 {
				if keyword != "true" {
					validKeyword = false
				}
				break
			}
		} else if initialCharacter == "f" {
			keyword += token
			if len(keyword) >= 5 {
				if keyword != "false" {
					validKeyword = false
				}
				break
			}
		} else if initialCharacter == "n" {
			keyword += token
			if len(keyword) >= 4 {
				if keyword != "null" {
					validKeyword = false
				}
				break
			}
		} else {
			keyword += token
			return keyword, false
		}
	}
	return keyword, validKeyword
}

// colorizer that outputs the formatted and coloured JSON in HTML form
// takes in a slice of tokens
func colorizer(tokens []Token) {
	var indentationLevel int = 0
	var indentationAmount int = 30
	var printingArray bool = false

	fmt.Printf("<span style=\"font-family:monospace; white-space:pre\">\n")
	for index, token := range tokens {
		if token.tokenType == TypeBracket {
			if token.data == "}" {
				fmt.Printf("</div>")
			}
			if token.data == "{" && indentationLevel > 0 {
				fmt.Printf("<br />")
			}
			fmt.Printf("<span style=\"color:green;\">")
			fmt.Printf("%v", token.data)
			if token.data == "{" && indentationLevel == 0 {
				indentationLevel++
			}
			if token.data == "}" && index + 1 < len(tokens) && tokens[index + 1].tokenType == TypeComma {
				fmt.Printf("</span>")
			} else {
				fmt.Printf("</span><br />")
			}
			if token.data == "{" {
				fmt.Printf("<div style=\"margin-left: %vpx;\">", indentationAmount * indentationLevel)
			}
		} else if token.tokenType == TypeKeyword {
			fmt.Printf("<span style=\"color:red;\">")
			fmt.Printf("%v", token.data)
			fmt.Printf("</span>")
		} else if token.tokenType == TypeNumber {
			fmt.Printf("<span style=\"color:blue;\">")
			fmt.Printf("%v", token.data)
			fmt.Printf("</span>")
		} else if token.tokenType == TypeString {
			fmt.Printf("<span style=\"color:purple;\">")
			var backslash bool = false
			var hexadecimal bool = false
			var hexCounter int = 0
			for _, character := range token.data {
				if hexadecimal {
					hexCounter++
					fmt.Printf("%c", character)
					if hexCounter == 4 {
						backslash = false
						hexCounter = 0
						hexadecimal = false
						fmt.Printf("</span>")
					}
				} else if backslash && !hexadecimal {
					if character == '"' {
						fmt.Printf("%s", "&quot;")
					} else if character == '<' {
						fmt.Printf("%s", "&lt;")
					} else if character == '>' {
						fmt.Printf("%s", "&gt;")
					} else if character == '&' {
						fmt.Printf("%s", "&amp;")
					} else if character == '\'' {
						fmt.Printf("%s", "&apos;")
					} else {
						fmt.Printf("%c", character)
					}
					
					if character == '"' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == '\\' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == '/' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == 'b' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == 'f' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == 'n' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == 'r' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == 't' {
						backslash = false
						fmt.Printf("</span>")
					} else if character == 'u' {
						hexadecimal = true
					}
				} else if character == '\\' {
					backslash = true
					fmt.Printf("<span style=\"color:orange;\">")
					fmt.Printf("%c", character)
				} else {
					if character == '"' {
						fmt.Printf("%s", "&quot;")
					} else if character == '<' {
						fmt.Printf("%s", "&lt;")
					} else if character == '>' {
						fmt.Printf("%s", "&gt;")
					} else if character == '&' {
						fmt.Printf("%s", "&amp;")
					} else if character == '\'' {
						fmt.Printf("%s", "&apos;")
					} else {
						fmt.Printf("%c", character)
					}
				}
			}
			fmt.Printf("</span>")
		} else if token.tokenType == TypeComma {
			fmt.Printf("<span style=\"color:brown;\">")
			fmt.Printf("%v", token.data)
			fmt.Printf("</span>")
			if !printingArray {
				fmt.Printf("<br />")
			} else {
				fmt.Printf("&nbsp;")
			}
		} else if token.tokenType == TypeColon {
			fmt.Printf("<span style=\"color:teal;\">")
			fmt.Printf("&nbsp;%v&nbsp;", token.data)
			fmt.Printf("</span>")
		} else if token.tokenType == TypeSquareBracket {
			if token.data == "[" {
				printingArray = true
			} else {
				printingArray = false
			}
			fmt.Printf("<span style=\"color:cyan;\">")
			fmt.Printf("%v", token.data)
			fmt.Printf("</span>")
		}
	}
	fmt.Printf("</span>")
}

// helper function to determine if characters are valid hexidecimal
func isHexademicalCharacter(hexString string) bool {
	for _, character := range hexString {
		if character != '0' &&
		   character != '1' &&
		   character != '2' &&
		   character != '3' &&
		   character != '4' &&
		   character != '5' &&
		   character != '6' &&
		   character != '7' &&
		   character != '8' &&
		   character != '9' &&
		   character != 'a' &&
		   character != 'b' &&
		   character != 'c' &&
		   character != 'd' &&
		   character != 'e' &&
		   character != 'f' &&
		   character != 'A' &&
		   character != 'B' &&
		   character != 'C' &&
		   character != 'D' &&
		   character != 'E' &&
		   character != 'F' {
		   	return false
		   }
	}
	return true
}