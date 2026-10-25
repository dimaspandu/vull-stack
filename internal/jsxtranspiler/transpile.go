package jsxtranspiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type transpiler struct {
	src     []rune
	i       int
	factory string
}

var pragmaRe = regexp.MustCompile(`(?m)^\s*/\*\*?\s*@jsx\s+([A-Za-z_$][\w$]*)\s*\*/\s*$`)

func Compile(source, factory string) (string, error) {
	if factory == "" {
		factory = "d"
		if m := pragmaRe.FindStringSubmatch(source); m != nil {
			factory = m[1]
			source = pragmaRe.ReplaceAllString(source, "")
		}
	}
	t := &transpiler{
		src:     []rune(source),
		factory: factory,
	}
	return t.transformJS()
}

func (t *transpiler) errorf(format string, args ...any) error {
	return fmt.Errorf("%s at position %d", fmt.Sprintf(format, args...), t.i)
}

func (t *transpiler) skipWhitespace() {
	for t.i < len(t.src) && (t.src[t.i] == ' ' || t.src[t.i] == '\t' || t.src[t.i] == '\n' || t.src[t.i] == '\r' || t.src[t.i] == '\f' || t.src[t.i] == '\v') {
		t.i++
	}
}

func (t *transpiler) readIdentifier() (string, error) {
	start := t.i
	for t.i < len(t.src) {
		c := t.src[t.i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '$' || c == ':' || c == '.' || c == '-' {
			t.i++
			continue
		}
		break
	}
	if start == t.i {
		return "", t.errorf("Expected identifier")
	}
	return string(t.src[start:t.i]), nil
}

func (t *transpiler) readQuotedString() (string, error) {
	if t.i >= len(t.src) || (t.src[t.i] != '"' && t.src[t.i] != '\'') {
		return "", t.errorf("Expected quoted string")
	}
	quote := t.src[t.i]
	t.i++
	var value strings.Builder
	for t.i < len(t.src) {
		ch := t.src[t.i]
		t.i++
		if ch == '\\' {
			value.WriteRune(ch)
			if t.i < len(t.src) {
				value.WriteRune(t.src[t.i])
				t.i++
			}
			continue
		}
		if ch == quote {
			return strconv.Quote(value.String()), nil
		}
		value.WriteRune(ch)
	}
	return "", t.errorf("Unclosed string")
}

func (t *transpiler) readExpression() (string, error) {
	if t.i >= len(t.src) || t.src[t.i] != '{' {
		return "", t.errorf("Expected {")
	}
	t.i++
	start := t.i
	depth := 1
	var quote rune
	template := false

	for t.i < len(t.src) {
		ch := t.src[t.i]
		if quote != 0 {
			if ch == '\\' {
				t.i += 2
				continue
			}
			if ch == quote {
				quote = 0
			}
			t.i++
			continue
		}
		if template {
			if ch == '\\' {
				t.i += 2
				continue
			}
			if ch == '`' {
				template = false
				t.i++
				continue
			}
			t.i++
			continue
		}
		if ch == '"' || ch == '\'' {
			quote = ch
			t.i++
			continue
		}
		if ch == '`' {
			template = true
			t.i++
			continue
		}
		if ch == '{' {
			depth++
		}
		if ch == '}' {
			depth--
			if depth == 0 {
				expr := strings.TrimSpace(string(t.src[start:t.i]))
				t.i++
				return Compile(expr, t.factory)
			}
		}
		t.i++
	}
	return "", t.errorf("Unclosed expression")
}

func (t *transpiler) readAttributeValue() (string, error) {
	t.skipWhitespace()
	if t.i < len(t.src) && (t.src[t.i] == '"' || t.src[t.i] == '\'') {
		return t.readQuotedString()
	}
	if t.i < len(t.src) && t.src[t.i] == '{' {
		return t.readExpression()
	}
	return "true", nil
}

type attr struct {
	spread string
	name   string
	value  string
}

func (t *transpiler) readAttributes() ([]attr, bool, error) {
	attrs := []attr{}
	for t.i < len(t.src) {
		t.skipWhitespace()
		if t.i+1 < len(t.src) && t.src[t.i] == '/' && t.src[t.i+1] == '>' {
			t.i += 2
			return attrs, true, nil
		}
		if t.i < len(t.src) && t.src[t.i] == '>' {
			t.i++
			return attrs, false, nil
		}
		if t.i < len(t.src) && t.src[t.i] == '{' {
			expr, err := t.readExpression()
			if err != nil {
				return nil, false, err
			}
			expr = strings.TrimSpace(expr)
			if !strings.HasPrefix(expr, "...") {
				return nil, false, t.errorf("Expected spread attribute")
			}
			attrs = append(attrs, attr{spread: strings.TrimSpace(expr[3:])})
			continue
		}
		name, err := t.readIdentifier()
		if err != nil {
			return nil, false, err
		}
		t.skipWhitespace()
		value := "true"
		if t.i < len(t.src) && t.src[t.i] == '=' {
			t.i++
			v, err := t.readAttributeValue()
			if err != nil {
				return nil, false, err
			}
			value = v
		}
		attrs = append(attrs, attr{name: name, value: value})
	}
	return nil, false, t.errorf("Unclosed opening tag")
}

func (t *transpiler) buildProps(attrs []attr) string {
	if len(attrs) == 0 {
		return "null"
	}
	parts := make([]string, 0, len(attrs))
	for _, a := range attrs {
		if a.spread != "" {
			parts = append(parts, "...("+a.spread+")")
			continue
		}
		parts = append(parts, fmt.Sprintf("%s: %s", strconv.Quote(a.name), a.value))
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

func normalizeText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	// collapse newlines followed by spaces
	for strings.Contains(text, "\n ") {
		text = strings.ReplaceAll(text, "\n ", "\n")
	}
	for strings.Contains(text, " \n") {
		text = strings.ReplaceAll(text, " \n", "\n")
	}
	fields := strings.Fields(text)
	return strings.Join(fields, " ")
}

func (t *transpiler) parseElement() (string, error) {
	if t.i >= len(t.src) || t.src[t.i] != '<' {
		return "", t.errorf("Expected <")
	}
	t.i++
	// fragment
	if t.i < len(t.src) && t.src[t.i] == '>' {
		t.i++
		children, err := t.parseChildren("")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.fragment(%s)", t.factory, strings.Join(children, ", ")), nil
	}
	tag, err := t.readIdentifier()
	if err != nil {
		return "", err
	}
	attrs, selfClosing, err := t.readAttributes()
	if err != nil {
		return "", err
	}
	props := t.buildProps(attrs)
	tagRef := strconv.Quote(tag)
	if tag[0] >= 'A' && tag[0] <= 'Z' {
		tagRef = tag
	}
	if selfClosing {
		return fmt.Sprintf("%s(%s, %s)", t.factory, tagRef, props), nil
	}
	children, err := t.parseChildren(tag)
	if err != nil {
		return "", err
	}
	args := []string{tagRef, props}
	args = append(args, children...)
	return fmt.Sprintf("%s(%s)", t.factory, strings.Join(args, ", ")), nil
}

func (t *transpiler) parseChildren(parentTag string) ([]string, error) {
	children := []string{}
	for t.i < len(t.src) {
		// closing tag
		if t.i+1 < len(t.src) && t.src[t.i] == '<' && t.src[t.i+1] == '/' {
			t.i += 2
			if t.i < len(t.src) && t.src[t.i] == '>' {
				t.i++
				return children, nil
			}
			closing, err := t.readIdentifier()
			if err != nil {
				return nil, err
			}
			t.skipWhitespace()
			if t.i >= len(t.src) || t.src[t.i] != '>' {
				return nil, t.errorf("Expected >")
			}
			t.i++
			if parentTag != "" && closing != parentTag {
				return nil, t.errorf("Expected </%s> but found </%s>", parentTag, closing)
			}
			return children, nil
		}
		if t.i < len(t.src) && t.src[t.i] == '<' {
			el, err := t.parseElement()
			if err != nil {
				return nil, err
			}
			children = append(children, el)
			continue
		}
		if t.i < len(t.src) && t.src[t.i] == '{' {
			expr, err := t.readExpression()
			if err != nil {
				return nil, err
			}
			if strings.TrimSpace(expr) != "" {
				children = append(children, expr)
			}
			continue
		}
		start := t.i
		for t.i < len(t.src) && t.src[t.i] != '<' && t.src[t.i] != '{' {
			t.i++
		}
		text := normalizeText(string(t.src[start:t.i]))
		if text != "" {
			children = append(children, strconv.Quote(text))
		}
	}
	if parentTag != "" {
		return nil, t.errorf("Missing </%s>", parentTag)
	}
	return children, nil
}

func (t *transpiler) transformJS() (string, error) {
	var out strings.Builder
	cursor := 0
	for cursor < len(t.src) {
		ch := t.src[cursor]
		// skip strings
		if ch == '"' || ch == '\'' || ch == '`' {
			quote := ch
			end := cursor + 1
			for end < len(t.src) {
				if t.src[end] == '\\' {
					end += 2
					continue
				}
				if t.src[end] == quote {
					end++
					break
				}
				end++
			}
			out.WriteString(string(t.src[cursor:end]))
			cursor = end
			continue
		}
		// potential JSX
		if ch == '<' && cursor+1 < len(t.src) {
			next := t.src[cursor+1]
			if (next >= 'A' && next <= 'Z') || (next >= 'a' && next <= 'z') || next == '>' {
				t.i = cursor
				jsx, err := t.parseElement()
				if err != nil {
					return "", err
				}
				out.WriteString(jsx)
				cursor = t.i
				continue
			}
		}
		out.WriteRune(ch)
		cursor++
	}
	return out.String(), nil
}
