import { Tag } from './tag.mjs';

const idRegex = /^[a-zA-Z][a-zA-Z0-9]*/;
const numberRegex = /^[0-9]+/;
const skipWhitespacesRegex = new RegExp(`^[ \r\t\f\v\n]*`);

export class Lexer {
    #program;
    #current;
    get suffix() {
        return this.#program.substring(this.#current);
    }
    set program(value) {
        this.#program = value;
        this.#current = 0;
    }
    nextToken() {
        let match = skipWhitespacesRegex.exec(this.suffix);
        this.#current += match[0].length;
        if (this.#current >= this.#program.length) {
            return { tag: Tag.EOF };
        }
        match = idRegex.exec(this.suffix);
        if (match && match[0].length) {
            this.#current += match[0].length;
            if (match[0] === 'var') {
                return { tag: Tag.VAR };
            }
            if (match[0] === 'if') {
                return { tag: Tag.IF };
            }
            if (match[0] === 'then') {
                return { tag: Tag.THEN };
            }
            if (match[0] === 'else') {
                return { tag: Tag.ELSE };
            }
            if (match[0] === 'endif') {
                return { tag: Tag.ENDIF };
            }
            if (match[0] === 'while') {
                return { tag: Tag.WHILE };
            }
            if (match[0] === 'do') {
                return { tag: Tag.DO };
            }
            if (match[0] === 'endwhile') {
                return { tag: Tag.ENDWHILE };
            }
            return { tag: Tag.ID, id: match[0] };
        }
        match = numberRegex.exec(this.suffix);
        if (match && match[0].length) { // Number: [0-9]+
            this.#current += match[0].length;
            return { tag: Tag.NUMBER, number: +match[0]};
        }
        switch (this.suffix[0]) {
            case '=':
                ++this.#current;
                if (this.suffix[0] === '=') {
                    ++this.#current;
                    return { tag: Tag.EQUAL };
                }
                return { tag: Tag.ASSIGN };
            case '+':
                ++this.#current;
                return { tag: Tag.PLUS };
            case '-':
                ++this.#current;
                return { tag: Tag.MINUS };
            case '*':
                ++this.#current;
                return { tag: Tag.MUL };
            case '/':
                ++this.#current;
                return { tag: Tag.DIV };
            case '<':
                ++this.#current;
                return { tag: Tag.LESS };
            case '>':
                ++this.#current;
                return { tag: Tag.MORE };
            case '(':
                ++this.#current;
                return { tag: Tag.LEFTBRACKET };
            case ')':
                ++this.#current;
                return { tag: Tag.RIGHTBRACKET };
            case ';':
                ++this.#current;
                return { tag: Tag.SEMICOLON };
            case ',':
                ++this.#current;
                return { tag: Tag.COMMA };
        }
        return { tag: Tag.EOF };
    }
}
