package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"iter"
	"slices"
	"strings"
)

var (
	ErrNotAPacisDirective    = errors.New("given value is not a pacis directive")
	ErrInvalidDirectiveType  = errors.New("invalid directive type")
	ErrInvalidDireciveParams = errors.New("invalid directive params")
	ErrMissingDirectiveParam = errors.New("directive is missing a parameter")
)

const (
	dirprefix = "//pacis:"
)

type DirectiveType int

const (
	InvalidDirective = DirectiveType(iota)
	PageDirective
	LayoutDirective
	RedirectDirective
	ActionDirective
	MiddlewareDirective
	LanguageDirective
	AuthenticationDirective
)

func (dt DirectiveType) String() string {
	switch dt {
	case PageDirective:
		return "page"
	case LayoutDirective:
		return "layout"
	case RedirectDirective:
		return "redirect"
	case ActionDirective:
		return "action"
	case MiddlewareDirective:
		return "middleware"
	case LanguageDirective:
		return "language"
	case AuthenticationDirective:
		return "authentication"
	default:
		return "unknown"
	}
}

var validdirs = []string{"page", "layout", "redirect", "action", "middleware", "language", "authentication"}

type Directive struct {
	Type     DirectiveType
	Params   map[string]string
	Position token.Position
	Node     ast.Decl
}

func parseparams(list []string) (map[string]string, error) {
	params := make(map[string]string)

	for _, item := range list {
		parts := strings.Split(item, "=")
		if len(parts) != 2 {
			return nil, errors.New("invalid number of param parts")
		}
		params[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return params, nil
}

func ParseComment(comment string, pos token.Position, decl ast.Decl) (*Directive, error) {
	after, ok := strings.CutPrefix(comment, dirprefix)
	if !ok {
		return nil, ErrNotAPacisDirective
	}
	parts := strings.Split(after, " ")
	if len(parts) == 0 {
		return nil, ErrInvalidDirectiveType
	}

	if !slices.Contains(validdirs, parts[0]) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDirectiveType, parts[0])
	}

	var typ DirectiveType
	switch parts[0] {
	case "page":
		typ = PageDirective
	case "layout":
		typ = LayoutDirective
	case "redirect":
		typ = RedirectDirective
	case "action":
		typ = ActionDirective
	case "middleware":
		typ = MiddlewareDirective
	case "language":
		typ = LanguageDirective
	case "authentication":
		typ = AuthenticationDirective
	default:
		typ = InvalidDirective
	}

	params := map[string]string{}
	if len(parts) >= 2 {
		var err error
		params, err = parseparams(parts[1:])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidDireciveParams, err)
		}
	}

	return &Directive{
		Type:     typ,
		Params:   params,
		Position: pos,
		Node:     decl,
	}, nil
}

type AstIter struct {
	file *ast.File
}

func (ai *AstIter) Comments() iter.Seq[*ast.CommentGroup] {
	return func(yield func(*ast.CommentGroup) bool) {
		var stopped bool
		ast.Inspect(ai.file, func(n ast.Node) bool {
			if stopped || n == nil {
				return false
			}
			comgrp, ok := n.(*ast.CommentGroup)
			if ok {
				if !yield(comgrp) {
					stopped = true
				}
			}
			return true
		})
	}
}

type DirectiveList struct {
	Page           []*Directive
	Action         []*Directive
	Layout         []*Directive
	Redirect       []*Directive
	Middleware     []*Directive
	Language       []*Directive
	Authentication []*Directive
}

func decl(pos token.Pos, file *ast.File) ast.Decl {
	for _, dcl := range file.Decls {
		fn, ok := dcl.(*ast.FuncDecl)
		if !ok {
			gen, ok := dcl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if gen.Tok != token.VAR && gen.Tok != token.CONST && gen.Tok != token.FUNC {
				continue
			}
			if gen.Doc == nil || len(gen.Doc.List) == 0 {
				continue
			}
		} else {
			if fn.Doc == nil || len(fn.Doc.List) == 0 {
				continue
			}
		}
		// Get start and end position of the declaration
		start := dcl.Pos()
		end := dcl.End()

		if pos >= start && pos <= end {
			return dcl
		}
	}
	return nil
}

func ParseDir(dir string) (*DirectiveList, error) {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	pkg, ok := pkgs["app"]
	if !ok {
		return nil, errors.New("failed to locate the app package")
	}

	list := &DirectiveList{}

	for _, file := range pkg.Files {
		iter := AstIter{file}

		for group := range iter.Comments() {
			for _, comment := range group.List {
				if strings.HasPrefix(comment.Text, dirprefix) {
					dir, err := ParseComment(comment.Text, fset.Position(comment.Pos()), decl(group.End()+1, file))
					if err != nil {
						return nil, err
					}

					switch dir.Type {
					case PageDirective:
						list.Page = append(list.Page, dir)
					case LayoutDirective:
						list.Layout = append(list.Layout, dir)
					case RedirectDirective:
						list.Redirect = append(list.Redirect, dir)
					case ActionDirective:
						list.Action = append(list.Action, dir)
					case MiddlewareDirective:
						list.Middleware = append(list.Middleware, dir)
					case LanguageDirective:
						list.Language = append(list.Language, dir)
					case AuthenticationDirective:
						list.Authentication = append(list.Authentication, dir)
					}
				}
			}
		}
	}

	return list, nil
}
