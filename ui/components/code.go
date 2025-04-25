package components

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

type CodeHighlighter struct {
	Language string
	Source   string
	Props    []h.I
}

func gettokclass(typ chroma.TokenType) string {
	switch typ {
	case chroma.PreWrapper:
		return "whitespace-pre-wrap"

	case chroma.Line:
		return "flex"

	case chroma.LineTable:
		return "border-spacing-0 p-0 m-0 border-0"
	case chroma.LineTableTD:
		return "align-top p-0 m-0 border-0"

	case chroma.LineLink:
		return "outline-0 decoration-0 text-inherit"

	case chroma.LineNumbers,
		chroma.LineNumbersTable:
		return "whitespace-pre-wrap select-none mr-4 py-2 text-neutral-400 dark:text-neutral-600"

	case chroma.LineHighlight:
		return "bg-sky-600/40"

	case chroma.Error:
		return "text-red-500 dark:text-red-400 bg-red-950 dark:bg-red-800"

	case chroma.Keyword,
		chroma.KeywordConstant,
		chroma.KeywordDeclaration,
		chroma.KeywordPseudo,
		chroma.KeywordReserved,
		chroma.KeywordType:
		return "text-rose-600 dark:text-rose-400"

	case chroma.KeywordNamespace:
		return "text-indigo-600 dark:text-indigo-400"

	case chroma.NameAttribute,
		chroma.NameClass,
		chroma.NameConstant,
		chroma.NameDecorator,
		chroma.NameEntity,
		chroma.NameException,
		chroma.NameFunction,
		chroma.NameFunctionMagic,
		chroma.NameKeyword,
		chroma.NameOperator:
		return "text-blue-600 dark:text-blue-400"

	case chroma.NameOther:
		return "text-muted-foreground"

	case chroma.Name,
		chroma.NameBuiltin,
		chroma.NameBuiltinPseudo,
		chroma.NameLabel,
		chroma.NameNamespace,
		chroma.NamePseudo,
		chroma.NameProperty,
		chroma.NameTag,
		chroma.NameVariable,
		chroma.NameVariableAnonymous,
		chroma.NameVariableClass,
		chroma.NameVariableGlobal,
		chroma.NameVariableInstance,
		chroma.NameVariableMagic:
		return "text-neutral-900 dark:text-neutral-100"

	case chroma.LiteralString,
		chroma.LiteralStringAffix,
		chroma.LiteralStringAtom,
		chroma.LiteralStringBacktick,
		chroma.LiteralStringBoolean,
		chroma.LiteralStringChar,
		chroma.LiteralStringDelimiter,
		chroma.LiteralStringDoc,
		chroma.LiteralStringDouble,
		chroma.LiteralStringEscape,
		chroma.LiteralStringHeredoc,
		chroma.LiteralStringInterpol,
		chroma.LiteralStringName,
		chroma.LiteralStringOther,
		chroma.LiteralStringRegex,
		chroma.LiteralStringSingle,
		chroma.LiteralStringSymbol:
		return "text-emerald-700 dark:text-emerald-300"

	case chroma.Literal,
		chroma.LiteralNumber,
		chroma.LiteralDate,
		chroma.LiteralOther,
		chroma.LiteralNumberBin,
		chroma.LiteralNumberFloat,
		chroma.LiteralNumberHex,
		chroma.LiteralNumberInteger,
		chroma.LiteralNumberIntegerLong,
		chroma.LiteralNumberOct,
		chroma.LiteralNumberByte:
		return "text-orange-600 dark:text-orange-400"

	case chroma.Operator, chroma.OperatorWord:
		return "text-neutral-200 dark:text-neutral-400"

	case chroma.Punctuation:
		return "text-muted-foreground"

	case chroma.Comment,
		chroma.CommentHashbang,
		chroma.CommentMultiline,
		chroma.CommentSingle,
		chroma.CommentSpecial,
		chroma.CommentPreproc,
		chroma.CommentPreprocFile:
		return "text-neutral-400 dark:text-neutral-600"

	case chroma.GenericEmph:
		return "italic"
	case chroma.GenericStrong:
		return "font-bold"
	case chroma.GenericUnderline:
		return "underline"
	default:
		return ""
	}
}

var htmlformatter = chroma.FormatterFunc(func(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	for token := range iterator.Stdlib() {
		class := gettokclass(token.Type)
		el := h.Span(
			h.If(len(class) > 0, h.Class(class)),

			h.Text(token.Value),
		)
		ctx := context.Background()
		ctxwer, ok := w.(*ctxwriter)
		if ok {
			ctx = ctxwer.ctx
		}

		err := el.Render(ctx, w)
		if err != nil {
			return err
		}
	}

	return nil
})

type ctxwriter struct {
	io.Writer
	ctx context.Context
}

var codecache = map[string][]byte{}

func (c *CodeHighlighter) Render(ctx context.Context, w io.Writer) error {
	hash := sha256.New()
	hash.Write([]byte(c.Source))
	sum := hex.EncodeToString(hash.Sum(nil))

	cache, ok := codecache[sum]
	if ok {
		_, err := w.Write(cache)
		return err
	}

	lexer := lexers.Get(c.Language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	iterator, err := lexer.Tokenise(nil, c.Source)
	if err != nil {
		return err
	}

	var buf = new(bytes.Buffer)
	if err := htmlformatter.Format(&ctxwriter{Writer: buf, ctx: ctx}, styles.Fallback, iterator); err != nil {
		return err
	}

	codecache[sum] = buf.Bytes()
	_, err = io.Copy(w, buf)
	return err
}

func (c *CodeHighlighter) NodeType() h.NodeType {
	return h.NodeElement
}

// TODO: Add a prop to include the copy button
func Code(src, lang string, props ...h.I) h.Node {
	return h.Div(
		Join(
			props,
			h.Class("relative overflow-x-auto"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				h.Class("!size-7 rounded-sm absolute top-2 right-2 md:top-3 md:right-3"),
				On("click", fn("$clipboard", src)),

				icons.Clipboard(h.Class("size-3")),
				h.Span(h.Class("sr-only"), h.Text("Copy to Clipboard")),
			),
			h.Pre(
				h.Class("p-3 pr-10 md:p-6 md:pr-12 overflow-auto bg-accent/50 dark:bg-accent/20 rounded-lg text-muted-foreground selection:bg-sky-600/20 whitespace-pre-wrap"),

				h.Cde(
					&CodeHighlighter{Language: lang, Source: src, Props: props},
				),
			),
		)...,
	)
}
