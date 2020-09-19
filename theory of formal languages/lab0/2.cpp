#include <iostream>
#include "vector"
#include "string"

using namespace std;

vector<int> perm = {0, 3, 2, 4, 1, 5, 6};


vector<vector<int>> decomposeCycles(vector<int> perm) {
    vector<vector<int>> cycles;
    vector<bool> isDecomposed(perm.size(), false);
    for (int i = 0; i < perm.size(); i++) {
        if (isDecomposed[i]) {
            continue;
        }
        isDecomposed[i] = true;
        vector<int> curCycle = {i};
        int cur = i;
        while (perm[cur] != i) {
            cur = perm[cur];
            curCycle.push_back(cur);
            isDecomposed[cur] = true;
        }
        cycles.push_back(curCycle);
    }
    return cycles;
}

vector<vector<int>> generatePerms(vector<int> perm) {
    vector<vector<int>> allPerms = {perm};
    vector<int> permResult = perm;
    for (;;) {
        for (int i = 0; i < perm.size(); i++) {
            permResult[i] = perm[permResult[i]];
        }
        if (find(allPerms.begin(), allPerms.end(), permResult) != allPerms.end()) {
            return allPerms;
        }
        allPerms.push_back(permResult);
    }
}


int main() {
    auto perms = generatePerms(perm);
    cout << "generated permutations : ";
    for (auto i: perms) {
        cout << "( ";
        for (auto j : i) {
            cout << j << " ";
        }
        cout << ") ";
    }

    auto cycles = decomposeCycles(perm);
    cout << "\ncycles : ";
    for (auto i: cycles) {
        if (i.size() == 1) {
            continue;
        }
        cout << "( ";
        for (auto j : i) {
            cout << j << " ";
        }
        cout << ")";
    }

    cout << "\norbits : ";
    string result = "";
    for (auto i: cycles) {
        result += "{";
        for (auto j : i) {
            result += to_string(j) + ",";
        }
        result[result.length() - 1] = '}';
        result += ", ";
    }
    result[result.length() - 1] = 0;
    result[result.length() - 2] = 0;
    cout << result;
    return 0;
}