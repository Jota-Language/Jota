package ast

type Statement interface {
    Accept(visitor StatementVisitor) interface{}
}

type StatementVisitor interface {
    VisitExpressionStatement(statement *ExpressionStatement) interface{}
    VisitPrintStatement(statement *PrintStatement) interface{}
    VisitVariableStatement(statement *VariableStatement) interface{}
    VisitBlockStatement(statement *BlockStatement) interface{}
    VisitIfStatement(statement *IfStatement) interface{}
    VisitWhileStatement(statement *WhileStatement) interface{}
    VisitFunctionStatement(statement *FunctionStatement) interface{}
}

type ExpressionStatement struct {
    Expression Expression
}

func (es *ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitExpressionStatement(es)
}

type PrintStatement struct {
    Expression Expression
}

func (ps *PrintStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitPrintStatement(ps)
}

type VariableStatement struct {
    Name Token
    Initializer Expression
}

func (vs *VariableStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitVariableStatement(vs)
}

type BlockStatement struct {
    Statements []Statement
}

func (bs *BlockStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitBlockStatement(bs)
}

type IfStatement struct {
    Condition Expression
    ThenBranch Statement
    ElseBranch Statement
}

func (is *IfStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitIfStatement(is)
}

type WhileStatement struct {
    Condition Expression
    Body Statement
}

func (ws *WhileStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitWhileStatement(ws)
}

type FunctionStatement struct {
    Name Token
    Params []Token
    Body []Statement
}

func (fs *FunctionStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitFunctionStatement(fs)
}
