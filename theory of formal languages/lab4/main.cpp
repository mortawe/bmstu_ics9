#include <iostream>
#include <fstream>
#include <regex>
#include <sstream>

using namespace std;

string read_file(const string &path) {
    ifstream source(path);
    ostringstream oss;
    oss << source.rdbuf();
    source.close();

    return oss.str();
}

vector<string> get_sentences(const string &source) {
    vector<string> sentences;

    regex regex(R"(([[:punct:]])(\s)?(\w+)(\s)?(?=[[:punct:]]))", regex::optimize);
    sregex_iterator regex_it(source.begin(), source.end(), regex), end;
    for (; regex_it != end; regex_it++) {
        smatch match_res = *regex_it;
        string match_str = match_res.str();
        sentences.push_back(match_str);
    }
    return sentences;
}


void write_sentences(const string &path, const vector<string> &sentences) {
    ofstream output_stream(path);
    output_stream << "<sentences>" << endl << endl;
    int c = 0;
    for (const string &s: sentences) {

        if (c != 0) {
            output_stream << "<s>" << s[0] << "</s>" << endl << endl;
            output_stream << "<b>" << s[0] << "</b>" << endl << endl;
        } else {
            output_stream << "<s>" << s[0] << "</s>" << endl << endl;
        }

        output_stream << "\t<w>";
        for (int i=1; i<s.length(); i++) {
            output_stream << s[i];
        }
        output_stream << " </w>"<< endl << endl;
        c++;
        if (c == sentences.size())
            output_stream << "<b>" << "." << "</b>" << endl << endl;

    }
    output_stream << "</sentences>" << endl;
    output_stream.close();
}


int main() {
    string input_path = "./input.txt";
    string output_path = "./output.txt";

    string input = read_file(input_path);
    write_sentences(output_path, get_sentences(input));

    return 0;
}