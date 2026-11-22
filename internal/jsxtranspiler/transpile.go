package jsxtranspiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Factory names the functions used to build elements and fragments.
type Factory struct {
	JSX     string
	JSXFrag string
}

const (
	defaultJSX     = "document.createElement"
	defaultJSXFrag = "document.createDocumentFragment"
)

// pragmaRe matches @jsx, @jsxRuntime and @jsxFrag pragma comments anywhere in
// the source, including the whitespace around them so they can be stripped.
var pragmaRe = regexp.MustCompile(`(?:\r?\n)?[ \t]*/\*\*?\s*@(jsx|jsxRuntime|jsxFrag)\s+([^\*]*?)\*/[ \t]*`)

// Compile transpiles the JSX in source. factory is the element factory name;
// when it is empty an @jsx pragma is used, then the default factory.
func Compile(source, factory string) (string, error) {
	return compile(source, factory)
}

// CompileWithFactory transpiles the JSX in source with explicit element and
// fragment factories. Empty fields fall back to a pragma, then to the default.
func CompileWithFactory(source string, factory Factory) (string, error) {
	return compile(source, factory)
}

type transpiler struct {
	src     []rune
	i       int
	jsx     string
	jsxFrag string
	factory any
}

func compile(source string, factory any) (string, error) {
	jsx, jsxFrag := resolveFactories(source, factory)

	// Pragmas are directives, so they never reach the output.
	source = pragmaRe.ReplaceAllString(source, "")

	t := &transpiler{
		src:     []rune(source),
		jsx:     jsx,
		jsxFrag: jsxFrag,
		factory: factory,
	}

	return t.transformJS()
}

// resolveFactories applies the same precedence as the JavaScript compiler:
// an explicit parameter wins over a pragma, a pragma wins over the defaults,
// and a string parameter makes the fragment fall back to "<jsx>.fragment".
func resolveFactories(source string, factory any) (string, string) {
	jsx, jsxFrag := "", ""
	stringFactory := false

	switch f := factory.(type) {
	case string:
		jsx = f
		stringFactory = f != ""
	case Factory:
		jsx = f.JSX
		jsxFrag = f.JSXFrag
	}

	for _, m := range pragmaRe.FindAllStringSubmatch(source, -1) {
		switch {
		case m[1] == "jsx" && jsx == "":
			jsx = strings.TrimSpace(m[2])
		case m[1] == "jsxFrag" && jsxFrag == "":
			jsxFrag = strings.TrimSpace(m[2])
		}
	}

	if jsx == "" {
		jsx = defaultJSX
	}

	if jsxFrag == "" {
		if stringFactory {
			jsxFrag = jsx + ".fragment"
		} else {
			jsxFrag = defaultJSXFrag
		}
	}

	return jsx, jsxFrag
}

func (t *transpiler) errorf(format string, args ...any) error {
	return fmt.Errorf("%s at position %d", fmt.Sprintf(format, args...), t.i)
}

func (t *transpiler) skipWhitespace() {
	for t.i < len(t.src) {
		switch t.src[t.i] {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			t.i++
		default:
			return
		}
	}
}

func (t *transpiler) skipTrivia() error {
	for {
		t.skipWhitespace()

		if t.i >= len(t.src) {
			return nil
		}

		if t.src[t.i] == '/' && t.i+1 < len(t.src) && t.src[t.i+1] == '/' {
			for t.i < len(t.src) && t.src[t.i] != '\n' {
				t.i++
			}
			continue
		}

		if t.src[t.i] == '/' && t.i+1 < len(t.src) && t.src[t.i+1] == '*' {
			t.i += 2

			end := indexRunes(t.src, "*/", t.i)
			if end < 0 {
				return t.errorf("Unterminated comment")
			}

			t.i = end + 2
			continue
		}

		return nil
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

		// Quoted string.
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

		// Template literal.
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
				expression := strings.TrimSpace(string(t.src[start:t.i]))
				t.i++

				// JSX may appear inside the expression, for example in either
				// branch of a conditional, so compile it with the same factory.
				return compile(expression, t.factory)
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
		if err := t.skipTrivia(); err != nil {
			return nil, false, err
		}

		// End of the opening tag.
		if t.i+1 < len(t.src) && t.src[t.i] == '/' && t.src[t.i+1] == '>' {
			t.i += 2
			return attrs, true, nil
		}

		if t.i < len(t.src) && t.src[t.i] == '>' {
			t.i++
			return attrs, false, nil
		}

		// Spread: {...props}
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

	// Collapse the whitespace a JSX text node may carry.
	fields := strings.Fields(text)

	return strings.Join(fields, " ")
}

func (t *transpiler) parseElement() (string, error) {
	if t.i >= len(t.src) || t.src[t.i] != '<' {
		return "", t.errorf("Expected <")
	}

	t.i++

	// Fragment: <>.
	if t.i < len(t.src) && t.src[t.i] == '>' {
		t.i++

		children, err := t.parseChildren("")
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s(%s)", t.jsxFrag, strings.Join(children, ", ")), nil
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
		return fmt.Sprintf("%s(%s, %s)", t.jsx, tagRef, props), nil
	}

	children, err := t.parseChildren(tag)
	if err != nil {
		return "", err
	}

	args := append([]string{tagRef, props}, children...)

	return fmt.Sprintf("%s(%s)", t.jsx, strings.Join(args, ", ")), nil
}

func (t *transpiler) parseChildren(parentTag string) ([]string, error) {
	children := []string{}

	for t.i < len(t.src) {
		// Closing tag.
		if t.i+1 < len(t.src) && t.src[t.i] == '<' && t.src[t.i+1] == '/' {
			t.i += 2

			// Fragment closing tag: </>.
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

		// Nested JSX.
		if t.i < len(t.src) && t.src[t.i] == '<' {
			element, err := t.parseElement()
			if err != nil {
				return nil, err
			}
			children = append(children, element)
			continue
		}

		// Expression.
		if t.i < len(t.src) && t.src[t.i] == '{' {
			expression, err := t.readExpression()
			if err != nil {
				return nil, err
			}
			if strings.TrimSpace(expression) != "" {
				children = append(children, expression)
			}
			continue
		}

		// Text.
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

		// Comments are copied verbatim, so an apostrophe or a "<" inside one is
		// never mistaken for a string or for JSX.
		if ch == '/' && cursor+1 < len(t.src) {
			switch t.src[cursor+1] {
			case '/':
				end := cursor + 2
				for end < len(t.src) && t.src[end] != '\n' {
					end++
				}
				out.WriteString(string(t.src[cursor:end]))
				cursor = end
				continue
			case '*':
				end := indexRunes(t.src, "*/", cursor+2)
				if end < 0 {
					end = len(t.src)
				} else {
					end += 2
				}
				out.WriteString(string(t.src[cursor:end]))
				cursor = end
				continue
			}
		}

		// Strings are copied verbatim.
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

		// Potential JSX.
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

func indexRunes(src []rune, sub string, from int) int {
	if from < 0 {
		from = 0
	}

	target := []rune(sub)
	if len(target) == 0 {
		return from
	}

	for i := from; i+len(target) <= len(src); i++ {
		match := true
		for j, r := range target {
			if src[i+j] != r {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}

	return -1
}
