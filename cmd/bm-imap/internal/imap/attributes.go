package imap

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Attribute struct {
	Name     string
	SubName  string
	Section  string
	Headers  []string
	Not      bool
	Peek     bool
	MinRange int
	MaxRange int
}

func (a Attribute) ToString() string {
	ret := a.Name

	if a.SubName != "" {
		ret += "." + a.SubName
	}


	if len(a.Headers) == 0 {
		if a.Section != "" {
			ret += "[" + a.Section + "]"
		}
	} else {
		ret += "[" + a.Section + " ("

		for _, h := range a.Headers {
			ret += "\"" + h + "\" "
		}

		ret += ")]"
	}

	if a.MaxRange > 0 {
		ret += "<" + strconv.Itoa(a.MinRange) + "." + strconv.Itoa(a.MaxRange) + ">"
	}

	return ret
}

// When encounting one of these macro's, replace with the contents instead
var attributeMacros = map[string][]string{
	"ALL":  {"FLAGS", "INTERNALDATE", "RFC822.SIZE", "ENVELOPE"},
	"FAST": {"FLAGS", "INTERNALDATE", "RFC822.SIZE"},
	"FULL": {"FLAGS", "INTERNALDATE", "RFC822.SIZE", "ENVELOPE", "BODY"},
}

func ParseAttributes(s string) []Attribute {
	var attrs []Attribute

	for _, field := range getFields(s) {
		// Check for RFC822.<item>
		if strings.HasPrefix(field, "RFC822.") {
			parts := strings.Split(field, ".")
			switch parts[1] {
			case "HEADER":
				attrs = append(attrs, Attribute{
					Name:    "BODY",
					SubName: "HEADER",
				})
				continue
			case "SIZE":
				attrs = append(attrs, Attribute{
					Name:    "RFC822",
					SubName: "SIZE",
				})
				continue
			case "TEXT":
				attrs = append(attrs, Attribute{
					Name:    "BODY",
					SubName: "TEXT",
				})
				continue
			}
		}

		// Handle BODY commands
		if strings.HasPrefix(field, "BODY.PEEK[") || strings.HasPrefix(field, "BODY[") {
			attr, err := getBodyAttribute(field)
			if err != nil {
				continue
			}

			attrs = append(attrs, *attr)
			continue
		}

		// Check for "regular" commands
		switch field {
		case "BODY":
			attrs = append(attrs, Attribute{
				Name: "BODY",
			})
		case "BODYSTRUCTURE":
			attrs = append(attrs, Attribute{
				Name: "BODYSTRUCTURE",
			})

		case "ENVELOPE":
			attrs = append(attrs, Attribute{
				Name: "ENVELOPE",
			})

		case "FLAGS":
			attrs = append(attrs, Attribute{
				Name: "FLAGS",
			})
		case "INTERNALDATE":
			attrs = append(attrs, Attribute{
				Name: "INTERNALDATE",
			})
		case "RFC822":
			attrs = append(attrs, Attribute{
				Name: "BODY",
			})
		case "UID":
			attrs = append(attrs, Attribute{
				Name: "UID",
			})
		}
	}

	return attrs
}

func getBodyAttribute(field string) (*Attribute, error) {
	attr := Attribute{
		Name: "BODY",
	}

	// if field == "BODY.PEEK[TEXT]" {
	// 	attr.Section = "TEXT"
	// 	return &attr, nil
	// }

	re := regexp.MustCompile("(.+)\\[([^\\(]+)(?: \\((.+)\\))?\\](?:\\<(\\d+).(\\d+)\\>)?")
	parts := re.FindStringSubmatch(field)
	if len(parts) != 6 {
		return nil, errors.New("incorrect field format")
	}

	if parts[1] == "BODY.PEEK" {
		attr.Peek = true
	}

	if strings.Contains(parts[2], ".NOT") {
		strings.Replace(parts[2], ".NOT", "", 1)
		attr.Not = true
	}

	attr.Section = parts[2]
	attr.Headers = strings.Split(strings.ToLower(parts[3]), " ")
	if len(attr.Headers) == 1 && attr.Headers[0] == "" {
		// Empty string means empty slice.
		attr.Headers = []string{}
	}

	attr.MinRange, _ = strconv.Atoi(parts[4])
	attr.MaxRange, _ = strconv.Atoi(parts[5])

	return &attr, nil
}

func getFields(s string) []string {
	s = strings.Trim(s, "()")

	// Expand macro if found
	for macro, attrs := range attributeMacros {
		if s == macro {
			return attrs
		}
	}

	// Simple state machine that deals with brackets
	idx := 0
	inSquareBracket := false
	var fields = []string{}

	ret := ""
	for idx < len(s) {
		if s[idx] == ' ' && !inSquareBracket {
			fields = append(fields, ret)
			ret = ""
			idx++
			continue
		}
		if s[idx] == '[' {
			inSquareBracket = true
		}
		if s[idx] == ']' {
			inSquareBracket = false
		}

		ret += string(s[idx])
		idx++
	}
	fields = append(fields, ret)

	return fields
}
