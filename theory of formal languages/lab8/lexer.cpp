#include <fstream>
#include <iostream>
#include <regex>
#include <string>
#include "token.cpp"
#include <vector>


const std::regex spacesRegex(" {2,}(?=([^']*'[^']*')*[^']*$)");
const std::regex whitespacesRegex(R"(^ +|[\r\n\t\f\v]+(?=([^']*'[^']*')*[^']*$))");
const std::regex oneLineCommentsRegex(R"(\/\/.*\n)");
const std::regex multiLineCommentsRegex("\\/\\*(.|\n)*?\\*\\/");
const std::regex idRegex("^[a-zA-Z][a-zA-Z0-9]*");
const std::regex intRegex("^[0-9][0-9]*");
const std::regex pathRegex(R"(^<[\.\,\_\/\-a-zA-Z0-9]*>)");
const std::regex charRegex("^'[a-zA-Z0-9]'");

class Lexer {
private:
    std::string code;
    std::vector<Token *> tokens;

    void normalize(std::string &code) {
        code = regex_replace(code, oneLineCommentsRegex, "");
        code = regex_replace(code, multiLineCommentsRegex, "");
        code = regex_replace(code, whitespacesRegex, " ");
        code = regex_replace(code, spacesRegex, " ");
    }

    uint LibraryName(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        std::smatch match;
        std::string substr = code.substr(current);
        if (regex_search(substr, match, pathRegex)) {
            tokens.push_back(new Token(TokenType::LIBRARY, match.str()));
            current += match.length();
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;

        }
        return current + 1;
    }

    uint _getTypeAndName(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        if (code.substr(current, 3) == "int") {
            tokens.push_back(new Token(TokenType::TYPE, "int"));
            current += 3;
        } else {
            if (code.substr(current, 4) == "char") {
                tokens.push_back(new Token(TokenType::TYPE, "char"));
                current += 4;
            } else {
                char err[1024];
                std::strcpy(err, std::to_string(current).c_str());
                throw err;
            }
        }
        for (; code[current] == ' '; ++current);
        std::smatch match;
        std::string substr = code.substr(current);
        if (regex_search(substr, match, idRegex)) {
            tokens.push_back(new Token(TokenType::ID, match.str()));
            current += match.length();
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;
        }
        return current;
    }

    uint ArgList(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        if (code.substr(current, 1) == "(") {
            tokens.push_back(new Token(TokenType::VALUE, "("));
            current += 1;
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;
        }
        for (; code[current] == ' '; ++current);
        bool first = true;
        while (code.substr(current, 1) != ")") {
            if (!first && code.substr(current, 1) == ",") {
                tokens.push_back(new Token(TokenType::VALUE, ","));
                current++;
            } else {
                if (!first) {
                    char err[1024];
                    std::strcpy(err, std::to_string(current).c_str());
                    throw err;
                }
            }
            first = false;
            current = _getTypeAndName(code, current, tokens);
        }
        tokens.push_back(new Token(TokenType::VALUE, ")"));
        current += 1;
        return current;

    }

    uint _getFuncIDAndArgs(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        std::smatch match;
        std::string substr = code.substr(current);
        if (regex_search(substr, match, idRegex)) {
            tokens.push_back(new Token(TokenType::ID, match.str()));
            current += match.length();
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;

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

    uint Comparable(std::string &code, uint current, std::vector<Token *> &tokens) {
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

    uint Addable(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        std::smatch match;
        std::string substr = code.substr(current);
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
                            char err[1024];
                            std::strcpy(err, std::to_string(current).c_str());
                            throw err;

                        }
                    }
                }
            }
        }
        return current;
    }

    uint Term(std::string &code, uint current, std::vector<Token *> &tokens) {
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

    uint SimpleSt(std::string &code, uint current, std::vector<Token *> &tokens) {
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

    uint DefSt(std::string &code, uint current, std::vector<Token *> &tokens) {
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
                    std::smatch match;
                    std::string substr = code.substr(current);
                    if (regex_search(substr, match, idRegex)) {
                        tokens.push_back(new Token(TokenType::ID, match.str()));
                        current += match.length();
                    } else {
                        char err[1024];
                        std::strcpy(err, std::to_string(current).c_str());
                        throw err;

                    }
                } else {
                    char err[1024];
                    std::strcpy(err, std::to_string(current).c_str());
                    throw err;

                }
            }
            for (; code[current] == ' '; ++current);
        }
        tokens.push_back(new Token(TokenType::VALUE, ";"));
        current++;
        return current;
    }

    uint CondSt(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        if (code.substr(current, 1) == "(") {
            tokens.push_back(new Token(TokenType::VALUE, "("));
            current++;
            for (; code[current] == ' '; ++current);
            if (code.substr(current, 1) != ")") {
                current = SimpleSt(code, current, tokens);
                for (; code[current] == ' '; ++current);
                if (code.substr(current, 1) != ")") {
                    char err[1024];
                    std::strcpy(err, std::to_string(current).c_str());
                    throw err;

                }
            }
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;

        }
        tokens.push_back(new Token(TokenType::VALUE, ")"));
        current++;
        for (; code[current] == ' '; ++current);
        return current;
    }

    uint WhileCycle(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        current = CondSt(code, current, tokens);
        for (; code[current] == ' '; ++current);
        current = Function_body(code, current, tokens);
        return current;
    }

    uint IfSt(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
        current = CondSt(code, current, tokens);
        for (; code[current] == ' '; ++current);
        current = Function_body(code, current, tokens);
        return current;
    }

    uint AssignSt(std::string &code, uint current, std::vector<Token *> &tokens) {
        std::smatch match;
        std::string substr = code.substr(current);
        if (regex_search(substr, match, idRegex)) {
            tokens.push_back(new Token(TokenType::ID, match.str()));
            current += match.length();
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;

        }
        for (; code[current] == ' '; ++current);
        if (code.substr(current, 1) == "=") {
            tokens.push_back(new Token(TokenType::VALUE, "="));
            current++;
            for (; code[current] == ' '; ++current);
            current = SimpleSt(code, current, tokens);
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;

        }
        if (code.substr(current, 1) != ";") {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;

        }
        tokens.push_back(new Token(TokenType::VALUE, ";"));
        current++;
        return current;
    }

    uint Statement(std::string &code, uint current, std::vector<Token *> &tokens) {
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
                    char err[1024];
                    std::strcpy(err, std::to_string(current).c_str());
                    throw err;

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
                    int before_size = tokens.size();
                    try {
                        current = DefSt(code, current, tokens);
                    }
                    catch (const char *e) {
                        found = false;
                        tokens.resize(before_size);
                    }
                    if (!found) {
                        found = true;
                        try {
                            current = AssignSt(code, current, tokens);
                            for (; code[current] == ' '; ++current);
                        }
                        catch (const char *e) {
                            found = false;
                            tokens.resize(before_size);
                        }
                        if (!found) {
                            current = SimpleSt(code, current, tokens);
                            for (; code[current] == ' '; ++current);
                            if (code.substr(current, 1) != ";") {
                                char err[1024];
                                std::strcpy(err, std::to_string(current).c_str());
                                throw err;

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

    uint Function_body(std::string &code, uint current, std::vector<Token *> &tokens) {
        if (code.substr(current, 1) == "{") {
            tokens.push_back(new Token(TokenType::VALUE, "{"));
            current++;
            while (code.substr(current, 1) != "}") {
                current = Statement(code, current, tokens);
                for (; code[current] == ' '; ++current);
            }
        } else {
            char err[1024];
            std::strcpy(err, std::to_string(current).c_str());
            throw err;

        }
        tokens.push_back(new Token(TokenType::VALUE, "}"));
        current++;
        return current;
    }

    uint Function(std::string &code, uint current, std::vector<Token *> &tokens) {
        for (; code[current] == ' '; ++current);
//        tokens.push_back(new Token(TokenType::VALUE, "function"));
        current = _getTypeAndName(code, current, tokens);
        current = ArgList(code, current, tokens);
        for (; code[current] == ' '; ++current);
        current = Function_body(code, current, tokens);
        return current;
    }

    void tokenize(std::string &code, std::vector<Token *> &tokens) {
        uint current = 0;
        try {
            while (current != code.length()) {
                for (; code[current] == ' '; ++current);
                if (code.substr(current, 8) == "#include") {
                    tokens.push_back(new Token(TokenType::VALUE, "include"));
                    current = LibraryName(code, current + 8, tokens);
                } else {
                    current = Function(code, current, tokens);
                    for (; code[current] == ' '; ++current);
                }
            }
        } catch (std::out_of_range &e) {
            std::cerr << "Unexpected end of code\n";
//        tokens = vector<Token *>();
        } catch (const char *e) {
            std::cerr << e << std::endl;
//        tokens = vector<Token *>();
        }
    }
public:
    std::vector<Token *> analyze(std::string inFileName) {
        std::ifstream inFile(inFileName);
        if (!inFile.is_open()) {
            std::cout << "shit";
            return std::vector<Token *>();
        }
        this->code = std::string((std::istreambuf_iterator<char>(inFile)), std::istreambuf_iterator<char>());
        this->normalize(this->code);
        // std::cout << this->code << "\n";
        this->tokens = std::vector<Token *>();
        this->tokenize(this->code, this->tokens);
        // for (Token* t : this->tokens) {
        // std::cout << t->toString() << "\n";
        // }
        inFile.close();
        return this->tokens;
    }
};