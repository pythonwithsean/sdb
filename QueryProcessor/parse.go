package queryprocessor

import (
	"errors"
	"fmt"
)

type Column struct {
	Type string
	Name string
}

type Condition struct {
	Column string
	Value  string
}

type CreateDatabase struct {
	Name string
}

type CreateTable struct {
	Name    string
	Columns []Column
}

type Select struct {
	All     bool
	Columns []string
	Table   string
	Where   *Condition
}

type Insert struct {
	Table  string
	Values []string
}

type Update struct {
	Table string
	Set   map[string]string
	Where *Condition
}

type Delete struct {
	Table string
	Where *Condition
}

type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) current() Token {
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) atEnd() bool {
	return p.pos >= len(p.tokens)
}

type Statement interface {
	statement()
}

func (*CreateDatabase) statement() {}
func (*CreateTable) statement()    {}
func (*Select) statement()         {}
func (*Insert) statement()         {}
func (*Update) statement()         {}
func (*Delete) statement()         {}

func (p *Parser) expect(t TokenType, value string) error {
	if p.atEnd() {
		return fmt.Errorf("expected %q, got end of input", value)
	}
	if p.current().Type != t || p.current().Value != value {
		return fmt.Errorf("expected %q, got %q", value, p.current().Value)
	}
	p.advance()
	return nil
}

func (p *Parser) expectKeyword(value string) error {
	return p.expect(KEYWORD, value)
}

func (p *Parser) expectSymbol(value string) error {
	return p.expect(SYMBOL, value)
}

func (p *Parser) expectIdentifier() (string, error) {
	if p.atEnd() || p.current().Type != IDENTIFIER {
		return "", fmt.Errorf("expected identifier")
	}
	v := p.current().Value
	p.advance()
	return v, nil
}

func (p *Parser) expectData() (string, error) {
	if p.atEnd() || p.current().Type != DATA {
		return "", fmt.Errorf("expected value")
	}
	v := p.current().Value
	p.advance()
	return v, nil
}

func (p *Parser) parseWhere() (*Condition, error) {
	if p.atEnd() || p.current().Value != "WHERE" {
		return nil, nil
	}
	p.advance()
	col, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.expectSymbol("="); err != nil {
		return nil, err
	}
	val, err := p.expectData()
	if err != nil {
		return nil, err
	}
	return &Condition{Column: col, Value: val}, nil
}

func (p *Parser) parseCreateDatabase() (*CreateDatabase, error) {
	if err := p.expectKeyword("DATABASE"); err != nil {
		return nil, err
	}
	name, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.expectSymbol(";"); err != nil {
		return nil, err
	}
	return &CreateDatabase{Name: name}, nil
}

func (p *Parser) parseCreateTable() (*CreateTable, error) {
	if err := p.expectKeyword("TABLE"); err != nil {
		return nil, err
	}
	name, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.expectSymbol("("); err != nil {
		return nil, err
	}
	columns := []Column{}
	for {
		if p.atEnd() || p.current().Type != KEYWORD {
			return nil, fmt.Errorf("expected column type")
		}
		colType := p.current().Value
		if p.current().Type != KEYWORD || (colType != "INT" && colType != "STRING") {
			return nil, fmt.Errorf("expected column type INT or STRING, got %q", colType)
		}
		p.advance()
		colName, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		columns = append(columns, Column{Type: colType, Name: colName})
		if p.atEnd() {
			return nil, errors.New("expected , or )")
		}
		if p.current().Value == "," {
			p.advance()
			continue
		}
		if p.current().Value == ")" {
			p.advance()
			break
		}
		return nil, fmt.Errorf("expected , or ), got %q", p.current().Value)
	}
	if err := p.expectSymbol(";"); err != nil {
		return nil, err
	}
	return &CreateTable{Name: name, Columns: columns}, nil
}

func (p *Parser) parseSelect() (*Select, error) {
	if err := p.expectKeyword("SELECT"); err != nil {
		return nil, err
	}
	sel := &Select{}
	if !p.atEnd() && p.current().Value == "*" {
		sel.All = true
		p.advance()
	} else {
		for {
			col, err := p.expectIdentifier()
			if err != nil {
				return nil, err
			}
			sel.Columns = append(sel.Columns, col)
			if p.atEnd() || p.current().Value != "," {
				break
			}
			p.advance()
		}
	}
	if err := p.expectKeyword("FROM"); err != nil {
		return nil, err
	}
	table, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	sel.Table = table
	where, err := p.parseWhere()
	if err != nil {
		return nil, err
	}
	sel.Where = where
	if err := p.expectSymbol(";"); err != nil {
		return nil, err
	}
	return sel, nil
}

func (p *Parser) parseInsert() (*Insert, error) {
	if err := p.expectKeyword("INSERT"); err != nil {
		return nil, err
	}
	if err := p.expectKeyword("INTO"); err != nil {
		return nil, err
	}
	table, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.expectKeyword("VALUES"); err != nil {
		return nil, err
	}
	if err := p.expectSymbol("("); err != nil {
		return nil, err
	}
	ins := &Insert{Table: table}
	for {
		v, err := p.expectData()
		if err != nil {
			return nil, err
		}
		ins.Values = append(ins.Values, v)
		if p.atEnd() || p.current().Value != "," {
			break
		}
		p.advance()
	}
	if err := p.expectSymbol(")"); err != nil {
		return nil, err
	}
	if err := p.expectSymbol(";"); err != nil {
		return nil, err
	}
	return ins, nil
}

func (p *Parser) parseUpdate() (*Update, error) {
	if err := p.expectKeyword("UPDATE"); err != nil {
		return nil, err
	}
	table, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.expectKeyword("SET"); err != nil {
		return nil, err
	}
	upd := &Update{Table: table, Set: map[string]string{}}
	for {
		col, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		if err := p.expectSymbol("="); err != nil {
			return nil, err
		}
		v, err := p.expectData()
		if err != nil {
			return nil, err
		}
		upd.Set[col] = v
		if p.atEnd() || p.current().Value != "," {
			break
		}
		p.advance()
	}
	where, err := p.parseWhere()
	if err != nil {
		return nil, err
	}
	upd.Where = where
	if err := p.expectSymbol(";"); err != nil {
		return nil, err
	}
	return upd, nil
}

func (p *Parser) parseDelete() (*Delete, error) {
	if err := p.expectKeyword("DELETE"); err != nil {
		return nil, err
	}
	if err := p.expectKeyword("FROM"); err != nil {
		return nil, err
	}
	table, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	del := &Delete{Table: table}
	where, err := p.parseWhere()
	if err != nil {
		return nil, err
	}
	del.Where = where
	if err := p.expectSymbol(";"); err != nil {
		return nil, err
	}
	return del, nil
}

func Parse(input string) (Statement, error) {
	tokens := Tokenize(input)
	for _, t := range tokens {
		if len(t.Value) > 0 && t.Value[0] == '"' && (len(t.Value) < 2 || t.Value[len(t.Value)-1] != '"') {
			return nil, fmt.Errorf("unterminated string: %s", t.Value)
		}
	}
	if len(tokens) == 0 {
		return nil, errors.New("empty input")
	}
	p := &Parser{tokens: tokens}
	if p.current().Type != KEYWORD {
		return nil, fmt.Errorf("expected a statement keyword, got %q", p.current().Value)
	}
	switch p.current().Value {
	case "CREATE":
		p.advance()
		if p.atEnd() {
			return nil, errors.New("expected DATABASE or TABLE")
		}
		if p.current().Value == "DATABASE" {
			return p.parseCreateDatabase()
		}
		if p.current().Value == "TABLE" {
			return p.parseCreateTable()
		}
		return nil, fmt.Errorf("expected DATABASE or TABLE, got %q", p.current().Value)
	case "SELECT":
		return p.parseSelect()
	case "INSERT":
		return p.parseInsert()
	case "UPDATE":
		return p.parseUpdate()
	case "DELETE":
		return p.parseDelete()
	}
	return nil, fmt.Errorf("unknown statement starting with %q", tokens[0].Value)
}
