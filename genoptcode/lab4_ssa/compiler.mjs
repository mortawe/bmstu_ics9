import { Lexer } from './lexer.mjs';
import { Parser } from './parser.mjs';

export class Compiler {
    #lexer;
    #parser;
    constructor() {
        this.#lexer = new Lexer();
        this.#parser = new Parser(this.#lexer);
    }
    outputCFG() {
        this.#parser.outputCFG();
    }
    outputCFGDOT() {
        return this.#parser.outputCFGDOT();
    }
    toSSA() {
        this.#parser.toSSA();
    }
    parse(program) {
        this.#lexer.program = program;
        this.#parser.init();
        this.#parser.parse();
    }
}
