package jsxtranspiler

import (
	"strings"
	"testing"
)

func compileForTest(t *testing.T, source string, factory any) string {
	t.Helper()

	var (
		got string
		err error
	)

	switch f := factory.(type) {
	case Factory:
		got, err = CompileWithFactory(source, f)
	case string:
		got, err = Compile(source, f)
	default:
		got, err = Compile(source, "")
	}

	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	return got
}

func TestCompile(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		factory any
		want    string
	}{
		// Default factory.
		{
			name:   "basic element",
			source: `const x = <div>Hello</div>;`,
			want:   `const x = document.createElement("div", null, "Hello");`,
		},
		{
			name:   "element with className",
			source: `<div className="foo" />`,
			want:   `document.createElement("div", { "className": "foo" })`,
		},
		{
			name:   "self-closing attributes",
			source: `<img src="x.png" alt="img" />`,
			want:   `document.createElement("img", { "src": "x.png", "alt": "img" })`,
		},
		{
			name:   "nested elements",
			source: `<div><span>A</span><span>B</span></div>`,
			want:   `document.createElement("div", null, document.createElement("span", null, "A"), document.createElement("span", null, "B"))`,
		},
		{
			name:   "expression child",
			source: `<h1>{title}</h1>`,
			want:   `document.createElement("h1", null, title)`,
		},
		{
			name:   "mixed children",
			source: `<div>Text <span>{val}</span></div>`,
			want:   `document.createElement("div", null, "Text", document.createElement("span", null, val))`,
		},
		{
			name:   "boolean attribute",
			source: `<input disabled />`,
			want:   `document.createElement("input", { "disabled": true })`,
		},
		{
			name:   "spread attributes",
			source: `<div {...props} />`,
			want:   `document.createElement("div", { ...(props) })`,
		},
		{
			name:   "spread attributes with whitespace after brace",
			source: `<StoryListItem { ...story } />`,
			want:   `document.createElement(StoryListItem, { ...(story) })`,
		},
		{
			name:   "spread attributes without whitespace after brace",
			source: `<StoryListItem {...story} />`,
			want:   `document.createElement(StoryListItem, { ...(story) })`,
		},
		{
			name:   "object expression attribute",
			source: `<button on={{ click: handler }}>Click</button>`,
			want:   `document.createElement("button", { "on": { click: handler } }, "Click")`,
		},
		{
			name:   "template literal expression",
			source: "<strong>{`Count: ${count}`}</strong>",
			want:   "document.createElement(\"strong\", null, `Count: ${count}`)",
		},
		{
			name:   "template literal attribute",
			source: "<div style={`color: red`} />",
			want:   "document.createElement(\"div\", { \"style\": `color: red` })",
		},
		{
			name: "jsx inside js function",
			source: `
    return (
      <div>Hello</div>
    );
  `,
			want: `
    return (
      document.createElement("div", null, "Hello")
    );
  `,
		},

		// Uppercase tags are component references.
		{
			name:   "uppercase component without attributes",
			source: `<Item />`,
			want:   `document.createElement(Item, null)`,
		},
		{
			name:   "uppercase component with attributes",
			source: `<Item text="Hello" />`,
			want:   `document.createElement(Item, { "text": "Hello" })`,
		},
		{
			name:   "uppercase component with children",
			source: `<Item>text</Item>`,
			want:   `document.createElement(Item, null, "text")`,
		},
		{
			name:   "uppercase component with expression child",
			source: `<Item>{value}</Item>`,
			want:   `document.createElement(Item, null, value)`,
		},
		{
			name:   "mixed lowercase and uppercase",
			source: `<div><Item text="x" /></div>`,
			want:   `document.createElement("div", null, document.createElement(Item, { "text": "x" }))`,
		},
		{
			name:   "multi-word component",
			source: `<FooBar>text</FooBar>`,
			want:   `document.createElement(FooBar, null, "text")`,
		},
		{
			name: "components in conditional and call",
			source: `import elementBuilder from "./factories/elementBuilder.js";

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
`,
			factory: "elementBuilder",
			want: `import elementBuilder from "./factories/elementBuilder.js";

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
`,
		},

		// String factory parameter.
		{
			name:    "custom factory",
			source:  `<div />`,
			factory: "h",
			want:    `h("div", null)`,
		},
		{
			name:    "custom factory with children",
			source:  `<div>Hello <span>x</span></div>`,
			factory: "jsx",
			want:    `jsx("div", null, "Hello", jsx("span", null, "x"))`,
		},
		{
			name:    "fragment with custom factory",
			source:  `<><span>x</span></>`,
			factory: "h",
			want:    `h.fragment(h("span", null, "x"))`,
		},

		// Pragmas.
		{
			name: "@jsx pragma",
			source: `
    /** @jsx h */
    const x = <div />;
  `,
			want: `
    const x = h("div", null);
  `,
		},
		{
			name: "@jsx pragma with children",
			source: `
    /** @jsx h */
    <div>Hello</div>
  `,
			want: `
    h("div", null, "Hello")
  `,
		},
		{
			name:   "@jsx pragma with dotted path factory",
			source: `/** @jsx document.createElement */<div />`,
			want:   `document.createElement("div", null)`,
		},
		{
			name:   "@jsx pragma with dotted path factory and children",
			source: `/** @jsx document.createElement */<div>Hello</div>`,
			want:   `document.createElement("div", null, "Hello")`,
		},
		{
			name:   "@jsxRuntime pragma",
			source: `/** @jsxRuntime classic */<div />`,
			want:   `document.createElement("div", null)`,
		},
		{
			name: "@jsxRuntime pragma strips comment",
			source: `/** @jsxRuntime classic */
<div />`,
			want: `
document.createElement("div", null)`,
		},
		{
			name:   "@jsxFrag pragma with dotted path",
			source: `/** @jsxFrag document.createDocumentFragment */<><div /></>`,
			want:   `document.createDocumentFragment(document.createElement("div", null))`,
		},
		{
			name:   "all three pragmas combined",
			source: `/** @jsxRuntime classic */ /** @jsx document.createElement */ /** @jsxFrag document.createDocumentFragment */ <><div>A</div></>`,
			want:   `document.createDocumentFragment(document.createElement("div", null, "A"))`,
		},
		{
			name:    "parameter wins over @jsx pragma",
			source:  `/** @jsx pragmaFactory */<div />`,
			factory: "paramFactory",
			want:    `paramFactory("div", null)`,
		},
		{
			name:    "@jsxFrag pragma with string factory",
			source:  `/** @jsxFrag frag */ <><div /></>`,
			factory: "h",
			want:    `frag(h("div", null))`,
		},

		// Object factory parameter.
		{
			name:    "object factory parameter",
			source:  `<div />`,
			factory: Factory{JSX: "h", JSXFrag: "h.fragment"},
			want:    `h("div", null)`,
		},
		{
			name:    "object factory parameter with children",
			source:  `<div>Hello</div>`,
			factory: Factory{JSX: "jsx"},
			want:    `jsx("div", null, "Hello")`,
		},
		{
			name:    "object factory parameter with custom fragment",
			source:  `<><span>x</span></>`,
			factory: Factory{JSX: "h", JSXFrag: "h.fragment"},
			want:    `h.fragment(h("span", null, "x"))`,
		},
		{
			name:    "object factory parameter with dotted paths",
			source:  `<div>Hello</div>`,
			factory: Factory{JSX: "document.createElement", JSXFrag: "document.createDocumentFragment"},
			want:    `document.createElement("div", null, "Hello")`,
		},
		{
			name:    "object factory parameter fragment uses jsxFrag",
			source:  `<><div>A</div></>`,
			factory: Factory{JSX: "document.createElement", JSXFrag: "document.createDocumentFragment"},
			want:    `document.createDocumentFragment(document.createElement("div", null, "A"))`,
		},
		{
			name:    "object factory parameter only jsx specified",
			source:  `<div />`,
			factory: Factory{JSX: "custom"},
			want:    `custom("div", null)`,
		},
		{
			name:    "object factory parameter only jsxFrag specified",
			source:  `<><div /></>`,
			factory: Factory{JSXFrag: "frag"},
			want:    `frag(document.createElement("div", null))`,
		},

		// Fragments.
		{
			name:   "fragment",
			source: `<><div>A</div><div>B</div></>`,
			want:   `document.createDocumentFragment(document.createElement("div", null, "A"), document.createElement("div", null, "B"))`,
		},

		// Comments inside opening tags.
		{
			name: "line comment between attributes",
			source: `<input
  value={query}
  // TODO: use onInput for live search
  onInput={(e) => vm.search(e.target.value)}
/>`,
			want: `document.createElement("input", { "value": query, "onInput": (e) => vm.search(e.target.value) })`,
		},
		{
			name:   "block comment between attributes",
			source: `<div /* layout */ className="x" />`,
			want:   `document.createElement("div", { "className": "x" })`,
		},
		{
			name: "line comment as first thing inside opening tag",
			source: `<div // first
  className="x" />`,
			want: `document.createElement("div", { "className": "x" })`,
		},
		{
			name:   "comment immediately before closing bracket",
			source: `<div className="x" /* end */>hi</div>`,
			want:   `document.createElement("div", { "className": "x" }, "hi")`,
		},
		{
			name:   "comment immediately before self-closing bracket",
			source: `<div className="x" /* end */ />`,
			want:   `document.createElement("div", { "className": "x" })`,
		},
		{
			name: "line comment immediately before self-closing bracket",
			source: `<div className="x" // end
/>`,
			want: `document.createElement("div", { "className": "x" })`,
		},
		{
			name: "multiple consecutive comments",
			source: `<div
  // one
  // two
  /* three */
  className="x"
/>`,
			want: `document.createElement("div", { "className": "x" })`,
		},
		{
			name:   "comment between spread and normal attribute",
			source: `<div {...props} /* spread */ className="x" />`,
			want:   `document.createElement("div", { ...(props), "className": "x" })`,
		},
		{
			name: "line comment between spread and normal attribute",
			source: `<div {...props}
  // spread
  className="x" />`,
			want: `document.createElement("div", { ...(props), "className": "x" })`,
		},
		{
			name:   "comment before spread attribute",
			source: `<div /* spread */ {...props} />`,
			want:   `document.createElement("div", { ...(props) })`,
		},
		{
			name:   "double slash in quoted attribute value is not a comment",
			source: `<div title="// not a comment" />`,
			want:   `document.createElement("div", { "title": "// not a comment" })`,
		},
		{
			name:   "block comment marker in quoted attribute value is not a comment",
			source: `<div title="/* not a comment */" />`,
			want:   `document.createElement("div", { "title": "/* not a comment */" })`,
		},
		{
			name:   "comment expression child passes through unchanged",
			source: `<div>{/* a comment */}</div>`,
			want:   `document.createElement("div", null, /* a comment */)`,
		},

		// Comments outside JSX.
		{
			name: "block comment before jsx is ignored",
			source: `/* a plain block comment */
const x = <div>Hi</div>;`,
			want: `/* a plain block comment */
const x = document.createElement("div", null, "Hi");`,
		},
		{
			name: "line comment before jsx is ignored",
			source: `// a plain line comment
const x = <div>Hi</div>;`,
			want: `// a plain line comment
const x = document.createElement("div", null, "Hi");`,
		},
		{
			name: "apostrophe inside block comment does not swallow jsx",
			source: `/*
 * The stateless createView builds a container and embeds a componentBuilder
 * instance. The ViewModel notifies listeners; the view's listener calls
 */
function createView(viewModel) {
  return <div>{viewModel.count}</div>;
}`,
			want: `/*
 * The stateless createView builds a container and embeds a componentBuilder
 * instance. The ViewModel notifies listeners; the view's listener calls
 */
function createView(viewModel) {
  return document.createElement("div", null, viewModel.count);
}`,
		},
		{
			name: "apostrophe inside line comment does not swallow jsx",
			source: `// the view's listener calls
const x = <div>Hi</div>;`,
			want: `// the view's listener calls
const x = document.createElement("div", null, "Hi");`,
		},
		{
			name: "multiple top-level comments before jsx",
			source: `/* first */
// second
/* third */
const x = <div>Hi</div>;`,
			want: `/* first */
// second
/* third */
const x = document.createElement("div", null, "Hi");`,
		},
		{
			name: "block comment between two jsx elements",
			source: `const a = <div>A</div>;
/* the view's listener */
const b = <div>B</div>;`,
			want: `const a = document.createElement("div", null, "A");
/* the view's listener */
const b = document.createElement("div", null, "B");`,
		},
		{
			name: "apostrophe in string before jsx is still skipped",
			source: `const name = "view's listener";
const x = <div>{name}</div>;`,
			want: `const name = "view's listener";
const x = document.createElement("div", null, name);`,
		},
		{
			name: "double slash inside block comment is not a line comment",
			source: `/* see https://example.com/docs */
const x = <div>Hi</div>;`,
			want: `/* see https://example.com/docs */
const x = document.createElement("div", null, "Hi");`,
		},

		// Larger samples.
		{
			name: "stateful component inside stateless view",
			source: `
    function createView(viewModel) {
      const stateful = componentBuilder(() => (
        <div>
          <h3>{viewModel.count}</h3>
          <button onClick={() => viewModel.increment()}>
            UPDATE
          </button>
        </div>
      ));

      const container = (
        <div>
          <h2>Local stateful inside stateless</h2>
          <p>Outer container is static; inner subtree rerenders via componentBuilder.</p>
          {stateful.render()}
        </div>
      );

      viewModel.addListener(() => stateful.rerender());

      return container;
    }
  `,
			want: `
    function createView(viewModel) {
      const stateful = componentBuilder(() => (
        document.createElement("div", null, document.createElement("h3", null, viewModel.count), document.createElement("button", { "onClick": () => viewModel.increment() }, "UPDATE"))
      ));

      const container = (
        document.createElement("div", null, document.createElement("h2", null, "Local stateful inside stateless"), document.createElement("p", null, "Outer container is static; inner subtree rerenders via componentBuilder."), stateful.render())
      );

      viewModel.addListener(() => stateful.rerender());

      return container;
    }
  `,
		},
		{
			name: "card-like component",
			source: `
    function card(title, description, content) {
      return (
        <section className="card">
          <h2>{title}</h2>
          <small>{description}</small>
          {content}
        </section>
      );
    }
  `,
			want: `
    function card(title, description, content) {
      return (
        document.createElement("section", { "className": "card" }, document.createElement("h2", null, title), document.createElement("small", null, description), content)
      );
    }
  `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compileForTest(t, tt.source, tt.factory)
			if got != tt.want {
				t.Errorf("Compile() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tt.want)
			}
		})
	}
}

func TestCompileErrors(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{"unmatched closing tag", `<div></span>`, "Expected </div> but found </span>"},
		{"missing closing tag", `<div>`, "Missing </div>"},
		{"unterminated block comment", `<div /* unterminated className="x" />`, "Unterminated comment"},
		{"unterminated expression", `<div>{value</div>`, "Unclosed expression"},
		{"unclosed string", `<div title="x />`, "Unclosed string"},
		{"expected identifier", `<div class="x" =y />`, "Expected identifier"},
		{"expected greater than", `<div></div`, "Expected >"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compile(tt.source, "")
			if err == nil {
				t.Fatalf("Compile() error = nil, want it to contain %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Compile() error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}
