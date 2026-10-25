package jsxtranspiler

import (
	"strings"
	"testing"
)

func TestCompileMatchesJSAnalyzer(t *testing.T) {
	tests := []struct {
		name, source, factory, want string
	}{
		{"basic element", `const x = <div>Hello</div>;`, "", `const x = d("div", null, "Hello");`},
		{"element with className", `<div className="foo" />`, "", `d("div", { "className": "foo" })`},
		{"self-closing attributes", `<img src="x.png" alt="img" />`, "", `d("img", { "src": "x.png", "alt": "img" })`},
		{"nested elements", `<div><span>A</span><span>B</span></div>`, "", `d("div", null, d("span", null, "A"), d("span", null, "B"))`},
		{"expression child", `<h1>{title}</h1>`, "", `d("h1", null, title)`},
		{"mixed children", `<div>Text <span>{val}</span></div>`, "", `d("div", null, "Text", d("span", null, val))`},
		{"boolean attribute", `<input disabled />`, "", `d("input", { "disabled": true })`},
		{"spread attributes", `<div {...props} />`, "", `d("div", { ...(props) })`},
		{"spread attributes with spaces", `<StoryListItem { ...story } />`, "", `d(StoryListItem, { ...(story) })`},
		{"object expression attribute", `<button on={{ click: handler }}>Click</button>`, "", `d("button", { "on": { click: handler } }, "Click")`},
		{"template literal expression", "<strong>{`Count: ${count}`}</strong>", "", "d(\"strong\", null, `Count: ${count}`)"},
		{"template literal attribute", "<div style={`color: red`} />", "", "d(\"div\", { \"style\": `color: red` })"},
		{"JS function context", "\n    return (\n      <div>Hello</div>\n    );\n  ", "", "\n    return (\n      d(\"div\", null, \"Hello\")\n    );\n  "},
		{"uppercase without attributes", `<Item />`, "", `d(Item, null)`},
		{"uppercase with attributes", `<Item text="Hello" />`, "", `d(Item, { "text": "Hello" })`},
		{"uppercase with children", `<Item>text</Item>`, "", `d(Item, null, "text")`},
		{"uppercase with expression child", `<Item>{value}</Item>`, "", `d(Item, null, value)`},
		{"mixed lowercase and uppercase", `<div><Item text="x" /></div>`, "", `d("div", null, d(Item, { "text": "x" }))`},
		{"multi-word component", `<FooBar>text</FooBar>`, "", `d(FooBar, null, "text")`},
		{"components in conditional and call", `import elementBuilder from "./factories/elementBuilder.js";

function ProgressIndicator() {
  return (
    <>
      <div class="skeleton" />
      <div class="skeleton" />
      <div class="skeleton" />
    </>
  );
}

function renderGreeting() {
  return <h2>Alright, Mate?</h2>
}

function App() {
  return (
    <div>
      <h1>Hello, World!</h1>
      {isLoading ? <ProgressIndicator /> : renderGreeting()}
    </div>
  );
}
`, "elementBuilder", `import elementBuilder from "./factories/elementBuilder.js";

function ProgressIndicator() {
  return (
    elementBuilder.fragment(elementBuilder("div", { "class": "skeleton" }), elementBuilder("div", { "class": "skeleton" }), elementBuilder("div", { "class": "skeleton" }))
  );
}

function renderGreeting() {
  return elementBuilder("h2", null, "Alright, Mate?")
}

function App() {
  return (
    elementBuilder("div", null, elementBuilder("h1", null, "Hello, World!"), isLoading ? elementBuilder(ProgressIndicator, null) : renderGreeting())
  );
}
`},
		{"custom factory", `<div />`, "h", `h("div", null)`},
		{"custom factory with children", `<div>Hello <span>x</span></div>`, "jsx", `jsx("div", null, "Hello", jsx("span", null, "x"))`},
		{"pragma", `
    /** @jsx h */
    const x = <div />;
  `, "", `
    const x = h("div", null);
  `},
		{"pragma with children", `
    /** @jsx h */
    <div>Hello</div>
  `, "", `
    h("div", null, "Hello")
  `},
		{"fragment", `<><div>A</div><div>B</div></>`, "", `d.fragment(d("div", null, "A"), d("div", null, "B"))`},
		{"fragment with custom factory", `<><span>x</span></>`, "h", `h.fragment(h("span", null, "x"))`},
		{"card-like component", `
    function card(title, description, content) {
      return (
        <section className="card">
          <h2>{title}</h2>
          <small>{description}</small>
          {content}
        </section>
      );
    }
  `, "", `
    function card(title, description, content) {
      return (
        d("section", { "className": "card" }, d("h2", null, title), d("small", null, description), content)
      );
    }
  `},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compile(tt.source, tt.factory)
			if err != nil {
				t.Fatalf("Compile() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Compile() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tt.want)
			}
		})
	}
}

func TestCompileRejectsUnmatchedClosingTag(t *testing.T) {
	if _, err := Compile(`<div></span>`, ""); err == nil || !strings.Contains(err.Error(), "Expected </div> but found </span>") {
		t.Fatalf("expected mismatched closing tag error, got %v", err)
	}
}
