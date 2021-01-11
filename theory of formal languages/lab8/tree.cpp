#include <string>
#include <vector>

int counter = 1;

struct TreeNode {
    int id;
    bool isTerminal;
    Token* token;
    std::string rule;
    std::vector<TreeNode*> children;
    TreeNode() {
        id = counter++;
    }
    TreeNode(bool isTerminal, Token* token) {
//        if (!isTerminal) {
//            throw "Token should be leaf node";
//        }
        this->isTerminal = isTerminal;
        this->token = token;
        id = counter++;
    }
    TreeNode(bool isTerminal, std::string rule) {
//        if (isTerminal) {
//            throw "Rule could not be leaf node";
//        }
        this->isTerminal = isTerminal;
        this->rule = rule;
        this->children = std::vector<TreeNode*>();
        id = counter++;
    }
    void append(TreeNode* node) {
//        if (this->isTerminal) {
//            throw "Could not append to leaf node";
//        }
        this->children.push_back(node);
    }
    void append(bool isTerminal, Token* token) {
//        if (this->isTerminal) {
//            throw "Could not append to leaf node";
//        }
        this->children.push_back(new TreeNode(isTerminal, token));
    }
    void append(bool isTerminal, std::string rule) {
//        if (this->isTerminal) {
//            throw "Could not append to leaf node";
//        }
        this->children.push_back(new TreeNode(isTerminal, rule));
    }
};