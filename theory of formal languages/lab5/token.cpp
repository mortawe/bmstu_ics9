#include <string>

enum class TokenType {
    ID, NUMBER, STRING, PATH, VALUE, LIBRARY, TYPE
};

class Token {
private:
    TokenType type;
    std::string value;

public:
    Token(TokenType type, std::string value): type(type), value(value) {}
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
};