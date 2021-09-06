// g++ -g -O3 main.cpp `llvm-config --cxxflags --ldflags --system-libs --libs core` -o compiler.out

#include "llvm/ADT/APFloat.h"
#include "llvm/ADT/STLExtras.h"
#include "llvm/IR/BasicBlock.h"
#include "llvm/IR/Constants.h"
#include "llvm/IR/DerivedTypes.h"
#include "llvm/IR/Function.h"
#include "llvm/IR/IRBuilder.h"
#include "llvm/IR/LLVMContext.h"
#include "llvm/IR/Module.h"
#include "llvm/IR/Type.h"
#include "llvm/IR/Verifier.h"
#include <algorithm>
#include <cctype>
#include <cstdio>
#include <cstdlib>
#include <map>
#include <memory>
#include <string>
#include <vector>

using namespace llvm;

//===----------------------------------------------------------------------===//
// Lexer
//===----------------------------------------------------------------------===//

// The lexer returns tokens [0-255] if it is an unknown character, otherwise one
// of these for known things.
enum struct Token {
    VAR, DEF, RETURN,
    ID, NUMBER,
    F, THEN, ELSE,
    FOR, DO,
    LEFTBRACKET, RIGHTBRACKET, SEMICOLON, COMMA,
    ASSIGN, PLUS, MINUS, MUL, DIV,
    LESS, MORE,
    EOFTOKEN,
};

static std::string IdentifierStr; // Filled in if ID
static double NumVal;             // Filled in if NUMBER

/// getToken - Return the next token from standard input.
static Token getToken() {
    static int LastChar = ' ';
    // Skip any whitespace.
    while (isspace(LastChar)) {
        LastChar = getchar();
    }
    if (isalpha(LastChar)) { // identifier: [a-zA-Z][a-zA-Z0-9]*
        IdentifierStr = LastChar;
        while (isalnum((LastChar = getchar()))) {
            IdentifierStr += LastChar;
        }
        if (IdentifierStr == "var") {
            return Token::VAR;
        }
        if (IdentifierStr == "def") {
            return Token::DEF;
        }
        if (IdentifierStr == "return") {
            return Token::RETURN;
        }
        if (IdentifierStr == "if") {
            return Token::IF;
        }
        if (IdentifierStr == "then") {
            return Token::THEN;
        }
        if (IdentifierStr == "else") {
            return Token::ELSE;
        }
        if (IdentifierStr == "for") {
            return Token::FOR;
        }
        if (IdentifierStr == "do") {
            return Token::DO;
        }
        return Token::ID;
    }
    if (isdigit(LastChar)) { // Number: [0-9][0-9]*\.?[0-9]*
        std::string NumStr;
        bool pointMet = false;
        do {
            if (LastChar == '.') {
                if (pointMet) {
                    // throw "unexpected .";
                    return Token::EOFTOKEN;
                }
                pointMet = true;
            }
            NumStr += LastChar;
            LastChar = getchar();
        } while (isdigit(LastChar) || LastChar == '.');
        NumVal = strtod(NumStr.c_str(), nullptr);
        return Token::NUMBER;
    }
    if (LastChar == '#') {
        // Comment until end of line.
        do {
            LastChar = getchar();
        } while (LastChar != EOF && LastChar != '\n' && LastChar != '\r');
        if (LastChar != EOF) {
            return getToken();
        }
    }
    switch (LastChar) {
        case '=':
            LastChar = getchar();
            return Token::ASSIGN;
        case '+':
            LastChar = getchar();
            return Token::PLUS;
        case '-':
            LastChar = getchar();
            return Token::MINUS;
        case '*':
            LastChar = getchar();
            return Token::MUL;
        case '/':
            LastChar = getchar();
            return Token::DIV;
        case '<':
            LastChar = getchar();
            return Token::LESS;
        case '>':
            LastChar = getchar();
            return Token::MORE;
        case '(':
            LastChar = getchar();
            return Token::LEFTBRACKET;
        case ')':
            LastChar = getchar();
            return Token::RIGHTBRACKET;
        case ';':
            LastChar = getchar();
            return Token::SEMICOLON;
        case ',':
            LastChar = getchar();
            return Token::COMMA;
        case EOF:
            return Token::EOFTOKEN;
    }
    // throw "unexpected symbol";
    return Token::EOFTOKEN;
}

//===----------------------------------------------------------------------===//
// Parse Tree
//===----------------------------------------------------------------------===//

namespace {

/// ExprAST - Base class for all expression nodes.
class ExprAST {
public:
    virtual ~ExprAST() = default;
    virtual Value* codegen() = 0;
};

/// NumberExprAST - Expression class for numeric literals like "1.0".
class NumberExprAST : public ExprAST {
    double Val;
public:
    NumberExprAST(double Val) : Val(Val) {}
    Value* codegen() override;
};

/// VariableExprAST - Expression class for referencing a variable, like "a".
class VariableExprAST : public ExprAST {
    std::string Name;
public:
    VariableExprAST(const std::string &Name) : Name(Name) {}
    Value* codegen() override;
    std::string getName() { return this->Name; }
};

/// LetAST - Expression class for referencing a variable declaration, like "let a = 0".
class LetAST : public ExprAST {
    std::string Name;
    std::unique_ptr<ExprAST> Expr;
public:
    LetAST(const std::string &Name)
        : Name(Name), Expr(std::make_unique<NumberExprAST>(0)) {}
    LetAST(const std::string &Name, std::unique_ptr<ExprAST> Expr)
        : Name(Name), Expr(std::move(Expr)) {}
    AllocaInst* codegen() override;
};

/// ReturnAST - Expression class for referencing a return expression, like "return 0".
class ReturnAST : public ExprAST {
    std::unique_ptr<ExprAST> Expr;
public:
    ReturnAST(std::unique_ptr<ExprAST> Expr): Expr(std::move(Expr)) {}
    Value* codegen() override;
};

/// BinaryExprAST - Expression class for a binary operator.
class BinaryExprAST : public ExprAST {
    Token Op;
    std::unique_ptr<ExprAST> LHS, RHS;
public:
    BinaryExprAST(Token Op, std::unique_ptr<ExprAST> LHS, std::unique_ptr<ExprAST> RHS)
        : Op(Op), LHS(std::move(LHS)), RHS(std::move(RHS)) {}
    Value* codegen() override;
};

/// CallExprAST - Expression class for function calls.
class CallExprAST : public ExprAST {
    std::string Callee;
    std::vector<std::unique_ptr<ExprAST>> Args;
public:
    CallExprAST(const std::string &Callee, std::vector<std::unique_ptr<ExprAST>> Args)
        : Callee(Callee), Args(std::move(Args)) {}
    Value* codegen() override;
};

/// PrototypeAST - This class represents the "prototype" for a function,
/// which captures its name, and its argument names (thus implicitly the number
/// of arguments the function takes).
class PrototypeAST {
    std::string Name;
    std::vector<std::string> Args;
public:
    PrototypeAST(const std::string &Name, std::vector<std::string> Args)
        : Name(Name), Args(std::move(Args)) {}
    Function* codegen();
    const std::string &getName() const { return Name; }
};

/// FunctionAST - This class represents a function definition itself.
class FunctionAST {
    std::unique_ptr<PrototypeAST> Proto;
    std::vector<std::unique_ptr<ExprAST>> Body;
public:
    FunctionAST(std::unique_ptr<PrototypeAST> Proto, std::vector<std::unique_ptr<ExprAST>> Body)
        : Proto(std::move(Proto)), Body(std::move(Body)) {}
    Function* codegen();
};

/// IfExprAST - Expression class for if/then/else.
class IfExprAST : public ExprAST {
    std::unique_ptr<ExprAST> Cond, Then, Else;
public:
    IfExprAST(std::unique_ptr<ExprAST> Cond,
              std::unique_ptr<ExprAST> Then,
              std::unique_ptr<ExprAST> Else)
        : Cond(std::move(Cond)), Then(std::move(Then)), Else(std::move(Else)) {}
    Value* codegen() override;
};

/// ForExprAST - Expression class for for/in.
class ForExprAST : public ExprAST {
    std::string VarName;
    std::unique_ptr<ExprAST> Start, End, Step, Body;
public:
    ForExprAST(const std::string &VarName, std::unique_ptr<ExprAST> Start,
               std::unique_ptr<ExprAST> End, std::unique_ptr<ExprAST> Step,
               std::unique_ptr<ExprAST> Body)
        : VarName(VarName), Start(std::move(Start)), End(std::move(End)),
          Step(std::move(Step)), Body(std::move(Body)) {}
    Value* codegen() override;
};

} // end anonymous namespace

//===----------------------------------------------------------------------===//
// Parser
//===----------------------------------------------------------------------===//

/// CurrentToken/getNextToken - Provide a simple token buffer.  CurrentToken is the current
/// token the parser is looking at.  getNextToken reads another token from the
/// lexer and updates CurrentToken with its results.
static Token CurrentToken;
static Token getNextToken() { return CurrentToken = getToken(); }

/// BinopPrecedence - This holds the precedence for each binary operator that is
/// defined.
static std::map<Token, int> BinopPrecedence;

/// GetTokPrecedence - Get the precedence of the pending binary operator token.
static int GetTokPrecedence() {
    if (CurrentToken != Token::PLUS && CurrentToken != Token::MINUS &&
        CurrentToken != Token::MUL && CurrentToken != Token::DIV &&
        CurrentToken != Token::LESS && CurrentToken != Token::MORE &&
        CurrentToken != Token::ASSIGN
    ) {
        return -1;
    }
    return BinopPrecedence[CurrentToken];
}

/// LogError* - These are little helper functions for error handling.
std::unique_ptr<ExprAST> LogError(const char* Str) {
    fprintf(stderr, "Error: %s\n", Str);
    return nullptr;
}

std::unique_ptr<PrototypeAST> LogErrorP(const char* Str) {
    LogError(Str);
    return nullptr;
}

std::unique_ptr<FunctionAST> LogErrorF(const char* Str) {
    LogError(Str);
    return nullptr;
}

static std::unique_ptr<ExprAST> ParseExpression();

/// numberexpr ::= number
static std::unique_ptr<ExprAST> ParseNumberExpr() {
    auto Result = std::make_unique<NumberExprAST>(NumVal);
    getNextToken(); // consume the number
    return std::move(Result);
}

/// parenexpr ::= '(' expression ')'
static std::unique_ptr<ExprAST> ParseParenExpr() {
    getNextToken(); // eat (.
    auto V = ParseExpression();
    if (!V) {
        return nullptr;
    }
    if (CurrentToken != Token::RIGHTBRACKET) {
        return LogError("expected ')'");
    }
    getNextToken(); // eat ).
    return V;
}

/// identifierexpr
///   ::= identifier
///   ::= identifier '(' expression* ')'
static std::unique_ptr<ExprAST> ParseIdentifierExpr() {
    std::string IdName = IdentifierStr;
    getNextToken(); // eat identifier.
    if (CurrentToken != Token::LEFTBRACKET) {// Simple variable ref.
        return std::make_unique<VariableExprAST>(IdName);
    }
    // Call.
    getNextToken(); // eat (
    std::vector<std::unique_ptr<ExprAST>> Args;
    if (CurrentToken != Token::RIGHTBRACKET) {
        while (true) {
            if (auto Arg = ParseExpression()) {
                Args.push_back(std::move(Arg));
            } else {
                return nullptr;
            }
            if (CurrentToken == Token::RIGHTBRACKET) {
                break;
            }
            if (CurrentToken != Token::COMMA) {
                return LogError("Expected ')' or ',' in argument list");
            }
            getNextToken();
        }
    }
    // Eat the ')'.
    getNextToken();
    return std::make_unique<CallExprAST>(IdName, std::move(Args));
}

/// ifexpr ::= 'if' expression 'then' expression 'else' expression
static std::unique_ptr<ExprAST> ParseIfExpr() {
    getNextToken();  // eat the if.
    // condition.
    auto Cond = ParseExpression();
    if (!Cond) {
        return nullptr;
    }
    if (CurrentToken != Token::THEN) {
        return LogError("expected then");
    }
    getNextToken();  // eat the then
    auto Then = ParseExpression();
    if (!Then) {
        return nullptr;
    }
    if (CurrentToken != Token::ELSE) {
        return LogError("expected else");
    }
    getNextToken();
    auto Else = ParseExpression();
    if (!Else) {
        return nullptr;
    }
    return std::make_unique<IfExprAST>(std::move(Cond), std::move(Then), std::move(Else));
}

/// forexpr ::= 'for' identifier '=' expr ',' expr (',' expr)? 'in' expression
static std::unique_ptr<ExprAST> ParseForExpr() {
    getNextToken();  // eat the for.
    if (CurrentToken != Token::ID) {
        return LogError("expected identifier after for");
    }
    std::string IdName = IdentifierStr;
    getNextToken();  // eat identifier.
    if (CurrentToken != Token::ASSIGN) {
        return LogError("expected '=' after for");
    }
    getNextToken();  // eat '='.
    auto Start = ParseExpression();
    if (!Start) {
        return nullptr;
    }
    if (CurrentToken != Token::COMMA) {
        return LogError("expected ',' after for start value");
    }
    getNextToken();
    auto End = ParseExpression();
    if (!End) {
        return nullptr;
    }
    // The step value is optional.
    std::unique_ptr<ExprAST> Step;
    if (CurrentToken == Token::COMMA) {
        getNextToken();
        Step = ParseExpression();
        if (!Step) {
            return nullptr;
        }
    }
    if (CurrentToken != Token::DO) {
        return LogError("expected 'do' after for");
    }
    getNextToken();  // eat 'do'.
    auto Body = ParseExpression();
    if (!Body) {
        return nullptr;
    }
    return std::make_unique<ForExprAST>(IdName, std::move(Start), std::move(End),
                                        std::move(Step), std::move(Body));
}

/// primary
///   ::= identifierexpr
///   ::= numberexpr
///   ::= parenexpr
static std::unique_ptr<ExprAST> ParsePrimary() {
    switch (CurrentToken) {
    default:
        return LogError("unknown token when expecting an expression");
    case Token::ID:
        return ParseIdentifierExpr();
    case Token::NUMBER:
        return ParseNumberExpr();
    case Token::LEFTBRACKET:
        return ParseParenExpr();
    case Token::IF:
        return ParseIfExpr();
    case Token::FOR:
        return ParseForExpr();
    }
}

/// binoprhs
///   ::= ('+' primary)*
static std::unique_ptr<ExprAST> ParseBinOpRHS(int ExprPrec, std::unique_ptr<ExprAST> LHS) {
    // If this is a binop, find its precedence.
    while (true) {
        int TokPrec = GetTokPrecedence();
        // If this is a binop that binds at least as tightly as the current binop,
        // consume it, otherwise we are done.
        if (TokPrec < ExprPrec) {
            return LHS;
        }
        // Okay, we know this is a binop.
        Token BinOp = CurrentToken;
        getNextToken(); // eat binop
        // Parse the primary expression after the binary operator.
        auto RHS = ParsePrimary();
        if (!RHS) {
            return nullptr;
        }
        // If BinOp binds less tightly with RHS than the operator after RHS, let
        // the pending operator take RHS as its LHS.
        int NextPrec = GetTokPrecedence();
        if (TokPrec < NextPrec) {
            RHS = ParseBinOpRHS(TokPrec + 1, std::move(RHS));
            if (!RHS) {
                return nullptr;
            }
        }
        // Merge LHS/RHS.
        LHS = std::make_unique<BinaryExprAST>(BinOp, std::move(LHS), std::move(RHS));
    }
}

/// expression
///   ::= primary binoprhs
///   ::= identifierexpr '=' expression
static std::unique_ptr<ExprAST> ParseExpression() {
    auto LHS = ParsePrimary();
    if (!LHS) {
        return nullptr;
    }
    return ParseBinOpRHS(0, std::move(LHS));
}

/// prototype
///   ::= id '(' id* ')'
static std::unique_ptr<PrototypeAST> ParsePrototype() {
    if (CurrentToken != Token::ID) {
        return LogErrorP("Expected function name in prototype");
    }
    std::string FnName = IdentifierStr;
    getNextToken();
    if (CurrentToken != Token::LEFTBRACKET) {
        return LogErrorP("Expected '(' in prototype");
    }
    std::vector<std::string> ArgNames;
    while (getNextToken() == Token::ID) {
        ArgNames.push_back(IdentifierStr);
    }
    if (CurrentToken != Token::RIGHTBRACKET) {
        return LogErrorP("Expected ')' in prototype");
    }
    // success.
    getNextToken(); // eat ')'.
    return std::make_unique<PrototypeAST>(FnName, std::move(ArgNames));
}

/// let ::= 'let' variable ( '=' expression )?
static std::unique_ptr<ExprAST> ParseLet() {
    getNextToken(); // eat const.
    std::string ConstName = IdentifierStr;
    getNextToken();
    if (CurrentToken != Token::ASSIGN) {
        return std::make_unique<LetAST>(ConstName);
    }
    getNextToken(); // eat '='
    if (auto E = ParseExpression()) {
        return std::make_unique<LetAST>(ConstName, std::move(E));
    }
    return nullptr;
}

/// return ::= 'return' expression
static std::unique_ptr<ExprAST> ParseReturn() {
    getNextToken(); // eat return.
    if (CurrentToken == Token::SEMICOLON) {
        return std::make_unique<ReturnAST>(nullptr);
    }
    if (auto E = ParseExpression()) {
        return std::make_unique<ReturnAST>(std::move(E));
    }
    return nullptr;
}

/// definition ::= 'func' prototype ( expression | let )+ return
static std::unique_ptr<FunctionAST> ParseDefinition() {
    getNextToken(); // eat func.
    auto Proto = ParsePrototype();
    if (!Proto) {
        return nullptr;
    }
    std::vector<std::unique_ptr<ExprAST>> body;
    while (CurrentToken != Token::RETURN) {
        if (CurrentToken == Token::VAR) {
            if (auto C = ParseLet()) {
                body.push_back(std::move(C));
            } else {
                return LogErrorF("Error while parsing const declaration");
            }
        } else {
            if (auto E = ParseExpression()) {
                body.push_back(std::move(E));
            } else {
                return LogErrorF("Error while parsing expression");
            }
        }
        if (CurrentToken == Token::SEMICOLON) {
            getNextToken();
        }
    }
    if (auto R = ParseReturn()) {
        body.push_back(std::move(R));
    } else {
        return LogErrorF("Error while parsing return expression");
    }
    if (CurrentToken == Token::SEMICOLON) {
        getNextToken();
    }
    return std::make_unique<FunctionAST>(std::move(Proto), std::move(body));
}

/// toplevelexpr
///   ::= expression
///   ::= 'let' variable ( '=' expression )?
static std::unique_ptr<FunctionAST> ParseTopLevelExpr() {
    if (CurrentToken == Token::VAR) {
        if (auto C = ParseLet()) {
            if (CurrentToken == Token::SEMICOLON) {
                getNextToken();
            }
            // Make an anonymous proto.
            auto Proto = std::make_unique<PrototypeAST>("__anon_expr", std::vector<std::string>());
            std::vector<std::unique_ptr<ExprAST>> tmp;
            tmp.push_back(std::move(C));
            return std::make_unique<FunctionAST>(std::move(Proto), std::move(tmp));
        }
        return nullptr;
    }
    if (auto E = ParseExpression()) {
        if (CurrentToken == Token::SEMICOLON) {
            getNextToken();
        }
        // Make an anonymous proto.
        auto Proto = std::make_unique<PrototypeAST>("__anon_expr", std::vector<std::string>());
        std::vector<std::unique_ptr<ExprAST>> tmp;
        tmp.push_back(std::move(E));
        return std::make_unique<FunctionAST>(std::move(Proto), std::move(tmp));
    }
    return nullptr;
}

//===----------------------------------------------------------------------===//
// Code Generation
//===----------------------------------------------------------------------===//

static std::unique_ptr<LLVMContext> TheContext;
static std::unique_ptr<Module> TheModule;
static std::unique_ptr<IRBuilder<>> Builder;
static std::map<std::string, Value*> NamedValues;
static std::map<std::string, AllocaInst*> DeclaredValues;

Value* LogErrorV(const char* Str) {
    LogError(Str);
    return nullptr;
}

Value* NumberExprAST::codegen() {
    return ConstantFP::get(*TheContext, APFloat(Val));
}

Value* VariableExprAST::codegen() {
    // Look this variable up in the function.
    Value* V = NamedValues[Name];
    if (!V) {
        return LogErrorV("Unknown variable name");
    }
    return V;
}

Value* BinaryExprAST::codegen() {
    Value* L = LHS->codegen();
    Value* R = RHS->codegen();
    if (!L || !R) {
        return nullptr;
    }
    switch (Op) {
    case Token::ASSIGN:
        VariableExprAST* V;
        if (!(V = dynamic_cast<VariableExprAST*>(LHS.get())) || !DeclaredValues[V->getName()]) {
            return LogErrorV("unknown variable assignment");
        }
        NamedValues[V->getName()] = R;
        return Builder->CreateStore(R, DeclaredValues[V->getName()]);
    case Token::PLUS:
        return Builder->CreateFAdd(L, R, "addtmp");
    case Token::MINUS:
        return Builder->CreateFSub(L, R, "subtmp");
    case Token::MUL:
        return Builder->CreateFMul(L, R, "multmp");
    case Token::DIV:
        return Builder->CreateFDiv(L, R, "divtmp");
    case Token::LESS:
        L = Builder->CreateFCmpULT(L, R, "cmptmp");
        // Convert bool 0/1 to double 0.0 or 1.0
        return Builder->CreateUIToFP(L, Type::getDoubleTy(*TheContext), "booltmp");
    case Token::MORE:
        L = Builder->CreateFCmpUGT(L, R, "cmptmp");
        // Convert bool 0/1 to double 0.0 or 1.0
        return Builder->CreateUIToFP(L, Type::getDoubleTy(*TheContext), "booltmp");
    default:
        return LogErrorV("invalid binary operator");
    }
}

Value* CallExprAST::codegen() {
    // Look up the name in the global module table.
    Function* CalleeF = TheModule->getFunction(Callee);
    if (!CalleeF) {
        return LogErrorV("Unknown function referenced");
    }
    // If argument mismatch error.
    if (CalleeF->arg_size() != Args.size()) {
        return LogErrorV("Incorrect # arguments passed");
    }
    std::vector<Value*> ArgsV;
    for (unsigned i = 0, e = Args.size(); i != e; ++i) {
        ArgsV.push_back(Args[i]->codegen());
        if (!ArgsV.back()) {
            return nullptr;
        }
    }
    return Builder->CreateCall(CalleeF, ArgsV, "calltmp");
}

Function* PrototypeAST::codegen() {
    // Make the function type:  double(double,double) etc.
    std::vector<Type*> Doubles(Args.size(), Type::getDoubleTy(*TheContext));
    FunctionType* FT = FunctionType::get(Type::getDoubleTy(*TheContext), Doubles, false);
    Function* F = Function::Create(FT, Function::ExternalLinkage, Name, TheModule.get());
    // Set names for all arguments.
    unsigned Idx = 0;
    for (auto &Arg : F->args()) {
        Arg.setName(Args[Idx++]);
    }
    return F;
}

Function* FunctionAST::codegen() {
    Function* TheFunction = Proto->codegen();
    if (!TheFunction) {
        return nullptr;
    }
    // Create a new basic block to start insertion into.
    BasicBlock* BB = BasicBlock::Create(*TheContext, "entry", TheFunction);
    Builder->SetInsertPoint(BB);
    // Record the function arguments in the NamedValues map.
    NamedValues.clear();
    for (auto &Arg : TheFunction->args()) {
        NamedValues[std::string(Arg.getName())] = &Arg;
    }
    DeclaredValues.clear();
    for (uint i = 0; i < Body.size() - 1; ++i) {
        if (auto E = Body[i]->codegen()) {
            continue;
        } else {
            // Error reading body, remove function.
            TheFunction->eraseFromParent();
            return nullptr;
        }
    }
    // Finish off the function.
    Builder->CreateRet(Body[Body.size() - 1]->codegen());
    // Validate the generated code, checking for consistency.
    verifyFunction(*TheFunction);
    return TheFunction;
}

AllocaInst* LetAST::codegen() {
    AllocaInst* Let = Builder->CreateAlloca(Type::getDoubleTy(*TheContext), nullptr, Name);
    auto gen = Expr->codegen();
    if (!gen) {
        return nullptr;
    }
    Builder->CreateStore(gen, Let);
    NamedValues[Name] = gen;
    DeclaredValues[Name] = Let;
    return Let;
}

Value* ReturnAST::codegen() {
    if (!Expr) {
        return nullptr;
    }
    if (Value* RetVal = Expr->codegen()) {
        return RetVal;
    }
    return nullptr;
}

Value* IfExprAST::codegen() {
    Value* CondV = Cond->codegen();
    if (!CondV) {
        return nullptr;
    }
    // Convert condition to a bool by comparing non-equal to 0.0.
    CondV = Builder->CreateFCmpONE(CondV, ConstantFP::get(*TheContext, APFloat(0.0)), "ifcond");
    Function* TheFunction = Builder->GetInsertBlock()->getParent();
    // Create blocks for the then and else cases.  Insert the 'then' block at the
    // end of the function.
    BasicBlock* ThenBB = BasicBlock::Create(*TheContext, "then", TheFunction);
    BasicBlock* ElseBB = BasicBlock::Create(*TheContext, "else");
    BasicBlock* MergeBB = BasicBlock::Create(*TheContext, "ifcont");
    Builder->CreateCondBr(CondV, ThenBB, ElseBB);
    // Emit then value.
    Builder->SetInsertPoint(ThenBB);
    Value* ThenV = Then->codegen();
    if (!ThenV) {
        return nullptr;
    }
    Builder->CreateBr(MergeBB);
    // Codegen of 'Then' can change the current block, update ThenBB for the PHI.
    ThenBB = Builder->GetInsertBlock();
    // Emit else block.
    TheFunction->getBasicBlockList().push_back(ElseBB);
    Builder->SetInsertPoint(ElseBB);
    Value* ElseV = Else->codegen();
    if (!ElseV) {
        return nullptr;
    }
    Builder->CreateBr(MergeBB);
    // Codegen of 'Else' can change the current block, update ElseBB for the PHI.
    ElseBB = Builder->GetInsertBlock();
    // Emit merge block.
    TheFunction->getBasicBlockList().push_back(MergeBB);
    Builder->SetInsertPoint(MergeBB);
    PHINode* PN = Builder->CreatePHI(Type::getDoubleTy(*TheContext), 2, "iftmp");
    PN->addIncoming(ThenV, ThenBB);
    PN->addIncoming(ElseV, ElseBB);
    return PN;
}

// Output for-loop as:
//   ...
//   start = startexpr
//   goto loop
// loop:
//   variable = phi [start, loopheader], [nextvariable, loopend]
//   ...
//   bodyexpr
//   ...
// loopend:
//   step = stepexpr
//   nextvariable = variable + step
//   endcond = endexpr
//   br endcond, loop, endloop
// outloop:
Value* ForExprAST::codegen() {
    // Emit the start code first, without 'variable' in scope.
    Value* StartVal = Start->codegen();
    if (!StartVal) {
        return nullptr;
    }
    // Make the new basic block for the loop header, inserting after current block.
    Function* TheFunction = Builder->GetInsertBlock()->getParent();
    BasicBlock* PreheaderBB = Builder->GetInsertBlock();
    BasicBlock* LoopBB = BasicBlock::Create(*TheContext, "loop", TheFunction);
    // Insert an explicit fall through from the current block to the LoopBB.
    Builder->CreateBr(LoopBB);
    // Start insertion in LoopBB.
    Builder->SetInsertPoint(LoopBB);
    // Start the PHI node with an entry for Start.
    PHINode* Variable = Builder->CreatePHI(Type::getDoubleTy(*TheContext), 2, VarName);
    Variable->addIncoming(StartVal, PreheaderBB);
    // Within the loop, the variable is defined equal to the PHI node.  If it
    // shadows an existing variable, we have to restore it, so save it now.
    Value* OldVal = NamedValues[VarName];
    NamedValues[VarName] = Variable;
    // Emit the body of the loop.  This, like any other expr, can change the
    // current BB.  Note that we ignore the value computed by the body, but don't
    // allow an error.
    if (!Body->codegen()) {
        return nullptr;
    }
    // Emit the step value.
    Value* StepVal = nullptr;
    if (Step) {
        StepVal = Step->codegen();
        if (!StepVal) {
            return nullptr;
        }
    } else {
        // If not specified, use 1.0.
        StepVal = ConstantFP::get(*TheContext, APFloat(1.0));
    }
    Value* NextVar = Builder->CreateFAdd(Variable, StepVal, "nextvar");
    // Compute the end condition.
    Value* EndCond = End->codegen();
    if (!EndCond) {
        return nullptr;
    }
    // Convert condition to a bool by comparing non-equal to 0.0.
    EndCond = Builder->CreateFCmpONE(EndCond, ConstantFP::get(*TheContext, APFloat(0.0)),
                                     "loopcond");
    // Create the "after loop" block and insert it.
    BasicBlock* LoopEndBB = Builder->GetInsertBlock();
    BasicBlock* AfterBB = BasicBlock::Create(*TheContext, "afterloop", TheFunction);
    // Insert the conditional branch into the end of LoopEndBB.
    Builder->CreateCondBr(EndCond, LoopBB, AfterBB);
    // Any new code will be inserted in AfterBB.
    Builder->SetInsertPoint(AfterBB);
    // Add a new entry to the PHI node for the backedge.
    Variable->addIncoming(NextVar, LoopEndBB);
    // Restore the unshadowed variable.
    if (OldVal) {
        NamedValues[VarName] = OldVal;
    } else {
        NamedValues.erase(VarName);
    }
    // for expr always returns 0.0.
    return Constant::getNullValue(Type::getDoubleTy(*TheContext));
}

//===----------------------------------------------------------------------===//
// Top-Level parsing and JIT Driver
//===----------------------------------------------------------------------===//

static void InitializeModule() {
    // Open a new context and module.
    TheContext = std::make_unique<LLVMContext>();
    TheModule = std::make_unique<Module>("", *TheContext);
    // Create a new builder for the module.
    Builder = std::make_unique<IRBuilder<>>(*TheContext);
}

static void HandleDefinition() {
    if (auto FnAST = ParseDefinition()) {
        if (auto* FnIR = FnAST->codegen()) {
            fprintf(stderr, "Read function definition:");
            FnIR->print(errs());
            fprintf(stderr, "\n");
        }
    } else {
        // Skip token for error recovery.
        getNextToken();
    }
}

static void HandleTopLevelExpression() {
    // Evaluate a top-level expression into an anonymous function.
    if (auto FnAST = ParseTopLevelExpr()) {
        if (auto* FnIR = FnAST->codegen()) {
            fprintf(stderr, "Read top-level expression:");
            FnIR->print(errs());
            fprintf(stderr, "\n");
            // Remove the anonymous expression.
            FnIR->eraseFromParent();
        }
    } else {
        // Skip token for error recovery.
        getNextToken();
    }
}

/// top ::= definition | expression | ';'
static void MainLoop() {
    while (true) {
        switch (CurrentToken) {
        case Token::EOFTOKEN:
            return;
        case Token::SEMICOLON: // ignore top-level semicolons.
            getNextToken();
            break;
        case Token::DEF:
            HandleDefinition();
            break;
        default:
            HandleTopLevelExpression();
            break;
        }
    }
}

//===----------------------------------------------------------------------===//
// Main driver code.
//===----------------------------------------------------------------------===//

int main() {
    // Install standard binary operators.
    // 1 is lowest precedence.
    BinopPrecedence[Token::ASSIGN] = 1;
    BinopPrecedence[Token::LESS] = 2;
    BinopPrecedence[Token::MORE] = 2;
    BinopPrecedence[Token::PLUS] = 3;
    BinopPrecedence[Token::MINUS] = 3;
    BinopPrecedence[Token::MUL] = 4;
    BinopPrecedence[Token::DIV] = 4;
    // Prime the first token.
    getNextToken();
    // Make the module, which holds all the code.
    InitializeModule();
    // Run the main "interpreter loop" now.
    MainLoop();
    // Print out all of the generated code.
    TheModule->print(errs(), nullptr);
    return 0;
}