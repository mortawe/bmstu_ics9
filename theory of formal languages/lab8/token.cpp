#include <string>

enum class TokenType {
    ID, NUMBER, STRING, PATH, VALUE, LIBRARY, TYPE
};

struct Token {
    TokenType type;
    std::string value;
    Token(TokenType type, std::string value) : type(type), value(value) {}

    std::string toString() {
        std::string res = "<";
        switch (this->type) {
            case TokenType::ID:
                res += "id, " + this->value;
                break;
            case TokenType::NUMBER:
                res += "num, " + this->value;
                break;
            case TokenType::STRING:
                res += "str, " + this->value;
                break;
            case TokenType::PATH:
                res += "path, " + this->value;
                break;
            case TokenType::VALUE:
                res += this->value;
                break;
            case TokenType::LIBRARY:
                res += "library path, " + this->value;
                break;

            case TokenType::TYPE:
                res += "type, " + this->value;
                break;

        }
        return res + ">";
    }
    friend bool operator==(const Token& left, const Token& right);
    friend bool operator!=(const Token& left, const Token& right);
};

bool operator==(const Token &left, const Token &right) {
    return left.type == right.type && left.value == right.value;
}

bool operator!=(const Token &left, const Token &right) {
    return left.type != right.type || left.value != right.value;
}