#include <iostream>

using namespace std;
const int n = 7;

int main() {
    for (int i = 1; i < n; i++) {
        for (int j = 1; j < n; j++) {
            cout << i * j % n << " ";
        }
        cout << "\n";
    }
    return 0;
}
