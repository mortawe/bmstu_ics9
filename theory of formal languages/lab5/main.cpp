#include <iostream>
#include <fstream>
#include <regex>
#include <string>
#include "token.cpp"
#include <vector>

using namespace std;

const string inFileName = "example.c";
const string outFileName = "output.txt";
const regex spacesRegex(" {2,}(?=([^']*'[^']*')*[^']*$)");
const regex whitespacesRegex(R"(^ +|[\r\n\t\f\v]+(?=([^']*'[^']*')*[^']*$))");
const regex oneLineCommentsRegex(R"(\/\/.*\n)");
const regex multiLineCommentsRegex("\\/\\*(.|\n)*?\\*\\/");
const regex idRegex("^[a-zA-Z][a-zA-Z0-9]*");
const regex intRegex("^[0-9][0-9]*");
const regex pathRegex(R"(^<[\.\,\_\/\-a-zA-Z0-9]*>)");
const regex charRegex("^'[a-zA-Z0-9]'");

uint SimpleSt(string &code, uint current, vector<Token *> &tokens);

uint Addable(string &code, uint current, vector<Token *> &tokens);

uint Function(string &code, uint current, vector<Token *> &tokens);

uint Function_body(string &code, uint current, vector<Token *> &tokens);


void normalize(string &code) {
    code = regex_replace(code, oneLineCommentsRegex, "");
    code = regex_replace(code, multiLineCommentsRegex, "");
    code = regex_replace(code, whitespacesRegex, " ");
    code = regex_replace(code, spacesRegex, " ");
}

uint LibraryName(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    smatch match;
    string substr = code.substr(current);
    if (regex_search(substr, match, pathRegex)) {
        tokens.push_back(new Token(TokenType::LIBRARY, match.str()));
        current += match.length();
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token";
    }
    return current + 1;
}

uint _getTypeAndName(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 3) == "int") {
        tokens.push_back(new Token(TokenType::TYPE, "int"));
        current += 3;
    } else {
        if (code.substr(current, 4) == "char") {
            tokens.push_back(new Token(TokenType::TYPE, "char"));
            current += 4;
        } else {
            cerr << "current : " << current << endl;
            throw "Unexpected token";
        }
    }
    for (; code[current] == ' '; ++current);
    smatch match;
    string substr = code.substr(current);
    if (regex_search(substr, match, idRegex)) {
        tokens.push_back(new Token(TokenType::ID, match.str()));
        current += match.length();
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token";
    }
    return current;
}

uint ArgList(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 1) == "(") {
        tokens.push_back(new Token(TokenType::VALUE, "("));
        current += 1;
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token";
    }
    for (; code[current] == ' '; ++current);
    bool first = true;
    while (code.substr(current, 1) != ")") {
        if (!first && code.substr(current, 1) == ",") {
            current++;
        } else {
            if (!first) {
                cerr << "current : " << current << endl;
                throw "Unexpected token";
            }
        }
        first = false;
        current = _getTypeAndName(code, current, tokens);
    }
    tokens.push_back(new Token(TokenType::VALUE, ")"));
    current += 1;
    return current;

}

uint _getFuncIDAndArgs(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    smatch match;
    string substr = code.substr(current);
    if (regex_search(substr, match, idRegex)) {
        tokens.push_back(new Token(TokenType::ID, match.str()));
        current += match.length();
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token";
    }
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 1) == "(") {
        tokens.push_back(new Token(TokenType::VALUE, "("));
        current++;
        for (; code[current] == ' '; ++current);
        while (code.substr(current, 1) != ")") {
            current = SimpleSt(code, current, tokens);
            if (code.substr(current, 1) == ",") {
                tokens.push_back(new Token(TokenType::VALUE, ","));
                current++;
            }
            for (; code[current] == ' '; ++current);
        }
    }
    tokens.push_back(new Token(TokenType::VALUE, ")"));
    current++;
    return current;
}

uint Comparable(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    current = Addable(code, current, tokens);
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 1) == "-") {
        for (; code[current] == ' '; ++current);
        tokens.push_back(new Token(TokenType::VALUE, "-"));
        current++;
        current = Comparable(code, current, tokens);
    } else {
        if (code.substr(current, 1) == "+") {
            for (; code[current] == ' '; ++current);
            tokens.push_back(new Token(TokenType::VALUE, "+"));
            current++;
            current = Comparable(code, current, tokens);
        }
    }
    return current;
}

uint Addable(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    smatch match;
    string substr = code.substr(current);
    if (regex_search(substr, match, intRegex)) {
        tokens.push_back(new Token(TokenType::NUMBER, match.str()));
        current = current + match.length();
    } else {

        if (regex_search(substr, match, idRegex)) {
            auto temp = current + match.length();
            for (; code[temp] == ' '; ++temp);
            if (code.substr(temp, 1) == "(") {
                current = _getFuncIDAndArgs(code, current, tokens);
            } else {
                tokens.push_back(new Token(TokenType::ID, match.str()));
                current = current + match.length();
            }
        } else {
            if (regex_search(substr, match, charRegex)) {
                tokens.push_back(new Token(TokenType::STRING, match.str()));
                current = current + match.length();
            } else {
                if (code.substr(current, 1) == "(") {
                    tokens.push_back(new Token(TokenType::VALUE, "("));
                    current++;
                    current = SimpleSt(code, current, tokens);
                    for (; code[current] == ' '; ++current);
                    if (code.substr(current, 1) == ")") {
                        tokens.push_back(new Token(TokenType::VALUE, ")"));
                        current++;
                    } else {
                        cerr << "current : " << current << endl;
                        throw "Unexpected token";
                    }
                }
            }
        }
    }
    return current;
}

uint Term(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    current = Comparable(code, current, tokens);
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 1) == "<") {
        for (; code[current] == ' '; ++current);
        tokens.push_back(new Token(TokenType::VALUE, "<"));
        current++;
        current = Term(code, current, tokens);
    } else {
        if (code.substr(current, 1) == ">") {
            for (; code[current] == ' '; ++current);
            tokens.push_back(new Token(TokenType::VALUE, ">"));
            current++;
            current = Term(code, current, tokens);
        }
    }
    return current;
}

uint SimpleSt(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    current = Term(code, current, tokens);
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 2) == "||") {
        for (; code[current] == ' '; ++current);
        tokens.push_back(new Token(TokenType::VALUE, "||"));
        current += 2;
        current = SimpleSt(code, current, tokens);
    } else {
        if (code.substr(current, 2) == "&&") {
            for (; code[current] == ' '; ++current);
            tokens.push_back(new Token(TokenType::VALUE, "&&"));
            current += 2;
            current = SimpleSt(code, current, tokens);
        }
    }
    return current;
}

uint DefSt(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    current = _getTypeAndName(code, current, tokens);
    for (; code[current] == ' '; ++current);
    while (code.substr(current, 1) != ";") {
        if (code.substr(current, 1) == "=") {
            tokens.push_back(new Token(TokenType::VALUE, "="));
            current++;
            current = SimpleSt(code, current, tokens);
        } else {
            if (code.substr(current, 1) == ",") {
                tokens.push_back(new Token(TokenType::VALUE, ","));
                current++;
                for (; code[current] == ' '; ++current);
                smatch match;
                string substr = code.substr(current);
                if (regex_search(substr, match, idRegex)) {
                    tokens.push_back(new Token(TokenType::ID, match.str()));
                    current += match.length();
                } else {
                    cerr << "current : " << current << endl;
                    throw "Unexpected token";
                }
            } else {
                cerr << "current : " << current << endl;
                throw "Unexpected token";
            }
        }
        for (; code[current] == ' '; ++current);
    }
    tokens.push_back(new Token(TokenType::VALUE, ";"));
    current++;
    return current;
}

uint CondSt(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 1) == "(") {
        tokens.push_back(new Token(TokenType::VALUE, "("));
        current++;
        for (; code[current] == ' '; ++current);
        if (code.substr(current, 1) != ")") {
            current = SimpleSt(code, current, tokens);
            for (; code[current] == ' '; ++current);
            if (code.substr(current, 1) != ")") {
                cerr << "current : " << current << endl;
                throw "Unexpected token not )";
            }
        }
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token not (";
    }
    tokens.push_back(new Token(TokenType::VALUE, ")"));
    current++;
    for (; code[current] == ' '; ++current);
    return current;
}

uint WhileCycle(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    current = CondSt(code, current, tokens);
    for (; code[current] == ' '; ++current);
    current = Function_body(code, current, tokens);
    return current;
}

uint IfSt(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    current = CondSt(code, current, tokens);
    for (; code[current] == ' '; ++current);
    current = Function_body(code, current, tokens);
    return current;
}

uint AssignSt(string &code, uint current, vector<Token *> &tokens) {
    smatch match;
    string substr = code.substr(current);
    if (regex_search(substr, match, idRegex)) {
        tokens.push_back(new Token(TokenType::ID, match.str()));
        current += match.length();
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token";
    }
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 1) == "=") {
        tokens.push_back(new Token(TokenType::VALUE, "="));
        current++;
        for (; code[current] == ' '; ++current);
        current = SimpleSt(code, current, tokens);
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token";
    }
    if (code.substr(current, 1) != ";") {
        cerr << "current : " << current << endl;
        throw "There is no ; in the assign st";
    }
    tokens.push_back(new Token(TokenType::VALUE, ";"));
    current++;
    return current;
}

uint Statement(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    if (code.substr(current, 5) == "while") {
        tokens.push_back(new Token(TokenType::VALUE, "while"));
        current = WhileCycle(code, current + 5, tokens);
    } else {
        if (code.substr(current, 6) == "return") {
            tokens.push_back(new Token(TokenType::VALUE, "return"));
            current = SimpleSt(code, current + 6, tokens);
            for (; code[current] == ' '; ++current);
            if (code.substr(current, 1) != ";") {
                cerr << "current : " << current << endl;
                throw "There is no ; in the return st";
            }
            tokens.push_back(new Token(TokenType::VALUE, ";"));
            current++;
            for (; code[current] == ' '; ++current);

        } else {
            if (code.substr(current, 2) == "if") {
                tokens.push_back(new Token(TokenType::VALUE, "if"));
                current += 2;
                current = IfSt(code, current, tokens);
            } else {
                bool found = true;
                try {
                    current = DefSt(code, current, tokens);
                }
                catch (const char *e) {
                    found = false;
                }
                if (!found) {
                    found = true;
                    try {
                        current = AssignSt(code, current, tokens);
                        for (; code[current] == ' '; ++current);

                    }
                    catch (const char *e) {
                        found = false;
                    }
                    if (!found) {
                        current = SimpleSt(code, current, tokens);
                        for (; code[current] == ' '; ++current);
                        if (code.substr(current, 1) != ";") {
                            cerr << "current : " << current << endl;
                            throw "There is no ; in the simple st st";
                        }
                        tokens.push_back(new Token(TokenType::VALUE, ";"));
                        current++;
                    }
                }
            }
        }
    }

    return current;

}

uint Function_body(string &code, uint current, vector<Token *> &tokens) {
    if (code.substr(current, 1) == "{") {
        tokens.push_back(new Token(TokenType::VALUE, "{"));
        current++;
        while (code.substr(current, 1) != "}") {
            current = Statement(code, current, tokens);
        }
    } else {
        cerr << "current : " << current << endl;
        throw "Unexpected token";
    }
    tokens.push_back(new Token(TokenType::VALUE, "}"));
    current++;
    return current;
}

uint Function(string &code, uint current, vector<Token *> &tokens) {
    for (; code[current] == ' '; ++current);
    tokens.push_back(new Token(TokenType::VALUE, "function"));
    current = _getTypeAndName(code, current, tokens);
    current = ArgList(code, current, tokens);
    for (; code[current] == ' '; ++current);
    current = Function_body(code, current, tokens);
    return current;
}


void tokenize(string &code, vector<Token *> &tokens) {
    uint current = 0;
    try {
        while (current != code.length()) {
            for (; code[current] == ' '; ++current);
            if (code.substr(current, 8) == "#include") {
                tokens.push_back(new Token(TokenType::VALUE, "library"));
                current = LibraryName(code, current + 8, tokens);
            } else {
                current = Function(code, current, tokens);
                for (; code[current] == ' '; ++current);
            }
        }
    } catch (out_of_range &e) {
        cerr << "Unexpected end of code\n";
//        tokens = vector<Token *>();
    } catch (const char *e) {
        cerr << e << endl;
        cerr << "cur pos: " << current << "\n";
//        tokens = vector<Token *>();
    }
}

int main() {
    ifstream inFile(inFileName);
    if (!inFile.is_open()) {
        return 3;
    }
    ofstream outFile(outFileName);
    if (!outFile.is_open()) {
        return 2;
    }
    string code((istreambuf_iterator<char>(inFile)), istreambuf_iterator<char>());
    normalize(code);
    cout << code << "\n";
    vector<Token *> tokens = vector<Token *>();
    tokenize(code, tokens);
    for (Token *t : tokens) {
        cout << t->toString() << "\n";
    }
    inFile.close();
    outFile.close();
    return 0;
}