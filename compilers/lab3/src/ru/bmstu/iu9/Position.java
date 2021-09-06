package ru.bmstu.iu9;


class Position {
    private int line, pos;

    Position(int line, int pos) {
        this.line = line;
        this.pos = pos;
    }

    void next_line() {
        this.line += 1;
        this.pos = 1;
    }

    void next_pos(int pos) {
        this.pos += pos;
    }

    public String toString() {
        return "(" + this.line + "," + this.pos + ")";
    }
}