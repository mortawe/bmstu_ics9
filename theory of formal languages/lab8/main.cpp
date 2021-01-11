#include <iostream>
#include <fstream>
#include <regex>
#include <string>
#include <vector>
#include "parser.cpp"

const std::string inFileName = "example.c";
const std::string outFileName = "output.txt";


int main() {

    Parser parser = Parser();
    parser.analyze(inFileName);

    return 0;
}