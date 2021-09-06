#include <iostream>
#include "lexer.cpp"
#include "tree.cpp"
#include <vector>

void throwCode(int code) {

    char err[1024];
    std::strcpy(err, std::to_string(code).c_str());
    throw err;
}

class Parser {
private:
    std::vector<Token *> tokens;
    TreeNode *tree;
    Lexer lexer;

    void printTree(TreeNode *node) {
        if (node->isTerminal) {
            std::cout << node->id << "[label=\"" << node->token->toString() << "\"];\n";
        } else {
            std::cout << node->id << "[label=\"" << node->rule << "\"];\n";
        }
        for (TreeNode *n : node->children) {
            std::cout << node->id << " -- " << n->id << ";\n";
        }
        for (TreeNode *n : node->children) {

            printTree(n);
        }
    }

    uint SimpleSt(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "SimpleSt");
        current = Term(current, node);
        if (*tokens[current] == Token(TokenType::VALUE, "||")) {
            node->append(true, tokens[current]);
            current++;
            current = SimpleSt(current, node);
        } else {
            if (*tokens[current] == Token(TokenType::VALUE, "&&")) {
                node->append(true, tokens[current]);
                current++;
                current = SimpleSt(current, node);
            }
        }
        parent->append(node);
        return current;
    }

    uint Values(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "Values");
        current = SimpleSt(current, node);
        if (*tokens[current] == Token(TokenType::VALUE, ",")) {
            node->append(true, tokens[current]);
            current++;
            current = Values(current, node);
        }
        parent->append(node);
        return current;
    }


    uint Comparable(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "Comparable");
        current = Addable(current, node);
        if (*tokens[current] == Token(TokenType::VALUE, "+")) {
            node->append(true, tokens[current]);
            current++;
            current = Comparable(current, node);
        } else {
            if (*tokens[current] == Token(TokenType::VALUE, "-")) {
                node->append(true, tokens[current]);
                current++;
                current = Comparable(current, node);
            }
        }
        parent->append(node);
        return current;
    }

    uint Term(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "Term");
        current = Comparable(current, node);
        if (*tokens[current] == Token(TokenType::VALUE, "<")) {
            node->append(true, tokens[current]);
            current++;
            current = Term(current, node);
        } else {
            if (*tokens[current] == Token(TokenType::VALUE, ">")) {
                node->append(true, tokens[current]);
                current++;
                current = Term(current, node);
            }
        }
        parent->append(node);
        return current;
    }

    uint AssignSt(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "AssignSt");
        if (tokens[current]->type == TokenType::ID) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        if (*tokens[current] == Token(TokenType::VALUE, "=")) {
            node->append(true, tokens[current]);
            current++;
            current = SimpleSt(current, node);
        } else {
            throwCode(current);
        }
        parent->append(node);
        return current;
    }

    uint ReturnSt(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "ReturnSt");
        if (*tokens[current] == Token(TokenType::VALUE, "return")) {
            node->append(true, tokens[current]);
            current++;
            current = SimpleSt(current, node);
        } else {
            throwCode(current);
        }
        parent->append(node);
        return current;
    }

    uint CondSt(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "CondSt");
        if (*tokens[current] == Token(TokenType::VALUE, "(")) {
            node->append(true, tokens[current]);
            current++;
            current = SimpleSt(current, node);
            if (*tokens[current] != Token(TokenType::VALUE, ")")) {
                throwCode(current);
            }
            node->append(true, tokens[current]);
            current++;
        }

        parent->append(node);
        return current;
    }

    uint While(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "While");
        if (*tokens[current] == Token(TokenType::VALUE, "while")) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        current = CondSt(current, node);
        if (*tokens[current] == Token(TokenType::VALUE, "{")) {
            node->append(true, tokens[current]);
        } else {
            throwCode(current);
        }
        current++;

        if (*tokens[current] == Token(TokenType::VALUE, "}")) {
            node->append(true, tokens[current]);
        } else {
            current = FunctionBody(current, node);
        }
        current++;
        parent->append(node);
        return current;
    }

    uint IfSt(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "If");
        if (*tokens[current] == Token(TokenType::VALUE, "if")) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        current = CondSt(current, node);
        if (*tokens[current] == Token(TokenType::VALUE, "{")) {
            node->append(true, tokens[current]);
        } else {
            throwCode(current);
        }
        current++;

        if (*tokens[current] == Token(TokenType::VALUE, "}")) {
            node->append(true, tokens[current]);
        } else {
            current = FunctionBody(current, node);
        }
        current++;
        parent->append(node);
        return current;
    }

    uint DefineList(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "DefineList");
        if (tokens[current]->type == TokenType::ID) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        if (*tokens[current] == Token(TokenType::VALUE, ",")) {
            node->append(true, tokens[current]);
            current++;
            current = DefineList(current, node);
        } else {
            if (*tokens[current] == Token(TokenType::VALUE, "=")) {
                node->append(true, tokens[current]);
                current++;
                current = SimpleSt(current, node);
                if (*tokens[current] == Token(TokenType::VALUE, ",")) {
                    node->append(true, tokens[current]);
                    current++;
                    current = DefineList(current, node);
                }
            }
        }
        parent->append(node);
        return current;
    }

    uint DefineSt(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "DefineSt");
        if (tokens[current]->type == TokenType::TYPE) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        current = DefineList(current, node);
        parent->append(node);
        return current;
    }


    uint Addable(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "Addable");
        if (tokens[current]->type == TokenType::NUMBER ||
            tokens[current]->type == TokenType::STRING) {
            node->append(true, tokens[current]);
            current++;
        } else {
            if (tokens[current]->type == TokenType::ID) {
                node->append(true, tokens[current]);
                current++;
                if (*tokens[current] == Token(TokenType::VALUE, "(")) {
                    node->append(true, tokens[current]);
                    current++;
                    current = Values(current, node);
                    if (*tokens[current] != Token(TokenType::VALUE, ")")) {
                        throwCode(current);
                    }
                    node->append(true, tokens[current]);
                    current++;
                }
            } else {
                if (*tokens[current] == Token(TokenType::VALUE, "(")) {
                    node->append(true, tokens[current]);
                    current++;
                    current = SimpleSt(current, node);
                    if (*tokens[current] != Token(TokenType::VALUE, ")")) {
                        throwCode(current);
                    }
                    node->append(true, tokens[current]);
                    current++;
                }

            }

        }
        parent->append(node);
        return current;
    }

    uint Statement(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "Statement");
        bool works = true;
        try {
            current = IfSt(current, node);
            parent->append(node);
            return current;
        }
        catch (const char *e) {
            node->children.clear();
            works = false;
        }
        if (!works) {
            works = true;
            try {
                current = While(current, node);
                parent->append(node);
                return current;
            }
            catch (const char *e) {
                node->children.clear();
                works = false;
            }
            if (!works) {
                works = true;
                try {
                    current = AssignSt(current, node);
                }
                catch (const char *e) {
                    node->children.clear();
                    works = false;
                }
                if (!works) {
                    works = true;
                    try {
                        current = ReturnSt(current, node);
                    }
                    catch (const char *e) {
                        node->children.clear();
                        works = false;
                    }

                    if (!works) {
                        works = true;
                        try {
                            current = DefineSt(current, node);
                        }
                        catch (const char *e) {
                            node->children.clear();
                            works = false;
                        }
                        if (!works) {
                            works = true;
                            try {
                                current = SimpleSt(current, node);
                            }
                            catch (const char *e) {
                                node->children.clear();
                                works = false;
                            }
                            if (!works) {
                                throwCode(current);
                            }
                        }

                    }
                }
            }
        }
        if (*tokens[current] == Token(TokenType::VALUE, ";")) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        parent->append(node);
        return current;
    }

    uint FunctionBody(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "FunctionBody");
        current = Statement(current, node);
        if (*tokens[current] != Token(TokenType::VALUE, "}")) {
            current = FunctionBody(current, node);
        }
        parent->append(node);
        return current;
    }

    uint ArgList(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "ArgList");
        if (tokens[current]->type == TokenType::TYPE) {
            node->append(true, tokens[current]);
        } else {
            throwCode(current);
        }
        current++;
        if (tokens[current]->type == TokenType::ID) {
            node->append(true, tokens[current]);
        } else {
            throwCode(current);
        }
        current++;
        if (*tokens[current] == Token(TokenType::VALUE, ",")) {
            node->append(true, tokens[current]);
            current++;
            current = ArgList(current, node);
        }
        parent->append(node);
        return current;
    }

    uint Function(uint current, TreeNode *parent) {
        TreeNode *node = new TreeNode(false, "Function");
        if (tokens[current]->type == TokenType::TYPE) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        if (tokens[current]->type == TokenType::ID) {
            node->append(true, tokens[current]);
            current++;
        } else {
            throwCode(current);
        }
        if (*tokens[current] == Token(TokenType::VALUE, "(")) {
            node->append(true, tokens[current]);
        } else {
            throwCode(current);
        }
        current++;
        if (*tokens[current] == Token(TokenType::VALUE, ")")) {
            node->append(true, tokens[current]);
        } else {
            current = ArgList(current, node);
        }
        current++;
        if (*tokens[current] == Token(TokenType::VALUE, "{")) {
            node->append(true, tokens[current]);
        } else {
            throwCode(current);
        }
        current++;

        if (*tokens[current] == Token(TokenType::VALUE, "}")) {
            node->append(true, tokens[current]);
        } else {
            current = FunctionBody(current, node);
        }
        current++;
        parent->append(node);
        return current;
    }

    uint Program(uint current, TreeNode *parent) {
        if (current >= tokens.size()) {
            return current;
        }
        TreeNode *node = new TreeNode(false, "Program");
        if (*tokens[current] == Token(TokenType::VALUE, "include")) {
            node->append(true, tokens[current]);
            current++;
            if (tokens[current]->type != TokenType::LIBRARY) {
                throwCode(current);
            }
            node->append(true, tokens[current]);
            current++;
        } else {
            if (tokens[current]->type == TokenType::TYPE) {
                current = Function(current, node);
            } else {
                return current;
            }
        }
        current = Program(current, node);
        parent->append(node);
        return current;
    }

public:
    Parser() : lexer(Lexer()) {}

    TreeNode *analyze(std::string inFileName) {
        this->tokens = this->lexer.analyze(inFileName);
        for (Token *t : tokens) {
            std::cout << t->toString() << "\n";
        }
        uint current = 0;
        TreeNode *node = new TreeNode(false, "S");
        try {
            current = Program(current, node);
            this->tree = node;
        } catch (const char *e) {
            this->tree = node;
            std::cerr << e << std::endl;
        }
        this->printTree(this->tree);
        return this->tree;
    }
};