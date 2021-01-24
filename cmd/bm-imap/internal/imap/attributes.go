package imap

import (
	"regexp"
	"strings"
)

type Attribute struct {
	Name    string
	Section string
	Headers []string
	Not     bool
	Peek    bool
}

func (a Attribute) ToString() string {
	ret := a.Name

	if len(a.Headers) == 0 {
		if a.Section != "" {
			ret += "." + a.Section
		}
		return ret
	}

	ret += "["+a.Section+" ("

	for _, h := range a.Headers {
		ret += "\"" + h + "\" "
	}

	ret += ")]"

	return ret
}

// When encounting one of these macro's, replace with the contents instead
var attributeMacros = map[string][]string{
	"ALL": {"FLAGS", "INTERNALDATE", "RFC822.SIZE", "ENVELOPE"},
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
					Section: "HEADER",
				})
				continue
			case "SIZE":
				attrs = append(attrs, Attribute{
					Name:    "RFC822",
					Section: "SIZE",
				})
				continue
			case "TEXT":
				attrs = append(attrs, Attribute{
					Name:    "BODY",
					Section: "TEXT",
				})
				continue
			}
		}

		if strings.HasPrefix(field, "BODY.PEEK[") {
			attr := Attribute{
				Name: "BODY",
				Peek: true,
			}

			if field == "BODY.PEEK[TEXT]" {
				attr.Section = "TEXT"
				continue
			}

			re := regexp.MustCompile("\\[([\\S]+) \\((.+)\\)\\]")
			match := re.FindStringSubmatch(field)
			if len(match) != 3 {
				continue
			}

			if strings.Contains(match[1], ".NOT") {
				strings.Replace(match[1], ".NOT", "", 1)
				attr.Not = true
			}

			attr.Section = match[1]
			attr.Headers = strings.Split(match[2], " ")

			attrs = append(attrs, attr)
			continue
		}

		if strings.HasPrefix(field, "BODY[") {
			attr := Attribute{
				Name: "BODY",
				Peek: false,
			}

			if field == "BODY[TEXT]" {
				attr.Section = "TEXT"
				continue
			}
			if field == "BODY[MIME]" {
				attr.Section = "MIME"
				continue
			}

			re := regexp.MustCompile("\\[([\\S]+) \\((.+)\\)\\]")
			match := re.FindStringSubmatch(field)
			if len(match) != 3 {
				continue
			}

			if strings.Contains(match[1], ".NOT") {
				strings.Replace(match[1], ".NOT", "", 1)
				attr.Not = true
			}

			attr.Section = match[1]
			attr.Headers = strings.Split(strings.ToUpper(match[2]), " ")

			attrs = append(attrs, attr)
			continue
		}

		// Check for "regular" fields
		switch field {
		case "BODY":
			attrs = append(attrs, Attribute{
				Name:    "BODY",
			})
		case "BODYSTRUCTURE":
			attrs = append(attrs, Attribute{
				Name:    "BODYSTRUCTURE",
			})

		case "ENVELOPE":
			attrs = append(attrs, Attribute{
				Name:    "ENVELOPE",
			})

		case "FLAGS":
			attrs = append(attrs, Attribute{
				Name:    "FLAGS",
			})
		case "INTERNALDATE":
			attrs = append(attrs, Attribute{
				Name:    "INTERNALDATE",
			})
		case "RFC822":
			attrs = append(attrs, Attribute{
				Name:    "BODY",
			})
		case "UID":
			attrs = append(attrs, Attribute{
				Name: "UID",
			})
		}
	}

	return attrs
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
	for (idx < len(s)) {
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
