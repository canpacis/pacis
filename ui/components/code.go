package components

import (
	"context"
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

var htmlformatter = chroma.FormatterFunc(func(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	for token := range iterator.Stdlib() {
		var class string

		switch token.Type {
		case chroma.PreWrapper:
			class = "whitespace-pre-wrap"

		case chroma.Line:
			class = "flex"

		case chroma.LineTable:
			class = "border-spacing-0 p-0 m-0 border-0"
		case chroma.LineTableTD:
			class = "align-top p-0 m-0 border-0"

		case chroma.LineLink:
			class = "outline-0 decoration-0 text-inherit"

		case chroma.LineNumbers,
			chroma.LineNumbersTable:
			class = "whitespace-pre-wrap select-none mr-4 py-2 text-neutral-400 dark:text-neutral-600"

		case chroma.LineHighlight:
			class = "bg-sky-600/40"

		case chroma.Error:
			class = "text-red-500 dark:text-red-400 bg-red-950 dark:bg-red-800"

		case chroma.Other:

		case chroma.Keyword,
			chroma.KeywordConstant,
			chroma.KeywordDeclaration,
			chroma.KeywordPseudo,
			chroma.KeywordReserved,
			chroma.KeywordType:
			class = "text-rose-600 dark:text-rose-400"

		case chroma.KeywordNamespace:
			class = "text-indigo-600 dark:text-indigo-400"

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
			class = "text-blue-600 dark:text-blue-400"

		case chroma.NameOther:
			class = "text-muted-foreground"

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
			class = "text-neutral-900 dark:text-neutral-100"

		case chroma.Literal:
		case chroma.LiteralDate:

		case chroma.LiteralOther:

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
			class = "text-emerald-700 dark:text-emerald-300"

		case chroma.LiteralNumber,
			chroma.LiteralNumberBin,
			chroma.LiteralNumberFloat,
			chroma.LiteralNumberHex,
			chroma.LiteralNumberInteger,
			chroma.LiteralNumberIntegerLong,
			chroma.LiteralNumberOct,
			chroma.LiteralNumberByte:
			class = "text-orange-600 dark:text-orange-400"

		case chroma.Operator, chroma.OperatorWord:
			class = "text-neutral-200 dark:text-neutral-400"

		case chroma.Punctuation:
			class = "text-muted-foreground"

		case chroma.Comment,
			chroma.CommentHashbang,
			chroma.CommentMultiline,
			chroma.CommentSingle,
			chroma.CommentSpecial,
			chroma.CommentPreproc,
			chroma.CommentPreprocFile:
			class = "text-neutral-400 dark:text-neutral-600"

		case chroma.GenericEmph:
			class = "italic"
		case chroma.GenericStrong:
			class = "font-bold"
		case chroma.GenericUnderline:
			class = "underline"
		}

		el := h.Span(
			h.If(len(class) > 0, h.Class(class)),

			h.Text(token.Value),
		)
		err := el.Render(context.Background(), w)
		if err != nil {
			return err
		}
	}

	return nil
})

func (c *CodeHighlighter) Render(ctx context.Context, w io.Writer) error {
	lexer := lexers.Get(c.Language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// formatter := html.New(
	// 	html.WithClasses(true),
	// )
	style := styles.Fallback

	iterator, err := lexer.Tokenise(nil, c.Source)
	if err != nil {
		return err
	}

	return htmlformatter.Format(w, style, iterator)
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
