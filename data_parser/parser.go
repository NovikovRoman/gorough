package data_parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	Command = iota
	Number
	EOD
)

type pathToken struct {
	tokenType int
	text      string
	number    float64
}

type Segment struct {
	Key  string
	Data []float64
}

var (
	commandParams = map[string]int{
		"A": 7, "a": 7, "C": 6, "c": 6, "H": 1, "h": 1, "L": 2, "l": 2, "M": 2, "m": 2, "Q": 4, "q": 4, "S": 4,
		"s": 4, "T": 2, "t": 2, "V": 1, "v": 1, "Z": 0, "z": 0}

	reGarbage = regexp.MustCompile(`(?s)^([ \t\r\n,]+)`)
	reCommand = regexp.MustCompile(`(?s)^([aAcChHlLmMqQsStTvVzZ])`)
	reNumber  = regexp.MustCompile(`(?s)^(([-+]?[0-9]+(\.[0-9]*)?|[-+]?\.[0-9]+)([eE][-+]?[0-9]+)?)`)
)

func (p pathToken) isType(t int) bool {
	return p.tokenType == t
}

func ParsePath(d string) (segments []Segment, err error) {
	var (
		tokens []pathToken
		params []float64
	)
	if tokens, err = tokenize(d); err != nil {
		return
	}

	mode := "BOD"
	index := 0
	token := tokens[index]

	for !token.isType(EOD) {
		paramsCount := 0
		params = []float64{}

		if mode == "BOD" {
			if token.text == "M" || token.text == "m" {
				index++
				paramsCount = commandParams[token.text]
				mode = token.text

			} else {
				return ParsePath("M0,0" + d)
			}

		} else if token.isType(Number) {
			paramsCount = commandParams[mode]

		} else {
			index++
			paramsCount = commandParams[token.text]
			mode = token.text
		}

		if (index + paramsCount) >= len(tokens) {
			err = errors.New("Path data ended short. ")
			return
		}

		for i := index; i < index+paramsCount; i++ {
			numToken := tokens[i]
			if !numToken.isType(Number) {
				err = fmt.Errorf("Param not a number: %s, %s ", mode, numToken.text)
				return
			}
			params = append(params, numToken.number)
		}

		if _, ok := commandParams[mode]; !ok {
			err = errors.New("Bad segment: " + mode)
			return
		}

		segments = append(segments, Segment{
			Key:  mode,
			Data: params,
		})
		index += paramsCount
		token = tokens[index]

		if mode == "M" {
			mode = "L"

		} else if mode == "m" {
			mode = "l"
		}
	}

	return
}

func Serialize(segments []Segment) string {
	var tokens []string
	for _, s := range segments {
		tokens = append(tokens, s.Key)
		switch s.Key {
		case "C", "c":
			tokens = append(tokens,
				fmt.Sprintf("%g %g, %g %g, %g %g", s.Data[0], s.Data[1], s.Data[2], s.Data[3], s.Data[4], s.Data[5]))

		case "S", "s", "Q", "q":
			tokens = append(tokens,
				fmt.Sprintf("%g %g, %g %g", s.Data[0], s.Data[1], s.Data[2], s.Data[3]))

		default:
			for _, d := range s.Data {
				tokens = append(tokens, strconv.FormatFloat(d, 'f', -1, 64))
			}
		}
	}

	return strings.Join(tokens, " ")
}

func tokenize(d string) (tokens []pathToken, err error) {
	num := 0.0
	for d != "" {
		m := reGarbage.FindStringSubmatch(d)
		if len(m) > 0 {
			d = strings.Replace(d, m[1], "", 1)
			continue
		}

		m = reCommand.FindStringSubmatch(d)
		if len(m) > 0 {
			tokens = append(tokens, pathToken{tokenType: Command, text: m[1]})
			d = strings.Replace(d, m[1], "", 1)
			continue
		}

		m = reNumber.FindStringSubmatch(d)
		if len(m) > 0 {
			if num, err = strconv.ParseFloat(m[1], 64); err != nil {
				err = fmt.Errorf("%s: %s", err, d)
				return
			}
			tokens = append(tokens, pathToken{tokenType: Number, text: m[1], number: num})
			d = strings.Replace(d, m[1], "", 1)
			continue
		}

		err = errors.New("Unknown error: " + d)
		return
	}

	tokens = append(tokens, pathToken{tokenType: EOD})
	return
}
