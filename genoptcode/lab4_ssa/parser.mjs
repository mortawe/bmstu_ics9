import { Tag } from './tag.mjs';

const binOpTags = [ Tag.PLUS, Tag.MINUS, Tag.MUL, Tag.DIV, Tag.EQUAL, Tag.LESS, Tag.MORE ];

export class Parser {
    #lexer;
    #sym;
    #tree;
    #instruction;
    #id;
    #declared;
    constructor(lexer) {
        this.#lexer = lexer;
    }
    get id() {
        return ++this.#id;
    }
    init() {
        this.getNextToken();
        this.#id = 1;
        this.#tree = {
            id: this.#id,
            instructions: [],
            parents: [],
            children: [],
        };
        this.#instruction = [];
        this.#declared = [];
    }
    newNode(parent) {
        return {
            id: this.id,
            instructions: [],
            parents: (parent ? [parent] : []),
            children: [],
        };
    }
    getNextToken() {
        this.#sym = this.#lexer.nextToken();
    }
    parse() {
        let currentNode = this.#tree;
        do {
            currentNode = this.parseAxiom(currentNode);
            if (!currentNode) {
                return;
            }
        } while (this.#sym.tag !== Tag.EOF);
    }
    parseAxiom(currentNode) {
        if (this.#sym.tag === Tag.VAR) {
            this.#instruction.push(this.#sym);
            this.getNextToken();
            if (this.#sym.tag !== Tag.ID) {
                return;
            }
            this.#instruction.push(this.#sym);
            this.#declared.push(this.#sym.id);
            this.getNextToken();
            if (this.#sym.tag === Tag.ASSIGN) {
                this.#instruction.push(this.#sym);
                this.getNextToken();
                this.parseExpr();
            }
            if (this.#sym.tag !== Tag.SEMICOLON) {
                return;
            }
            this.#instruction.push(this.#sym);
            this.getNextToken();
            currentNode.instructions.push([...this.#instruction]);
            this.#instruction = [];
            return currentNode;
        }
        if (this.#sym.tag === Tag.ID) {
            this.#instruction.push({tag: Tag.ID, id: this.#sym.id });
            this.getNextToken();
            if (this.#sym.tag !== Tag.ASSIGN) {
                return;
            }
            this.#instruction.push(this.#sym);
            this.getNextToken();
            this.parseExpr();
            if (this.#sym.tag !== Tag.SEMICOLON) {
                return;
            }
            this.#instruction.push(this.#sym);
            this.getNextToken();
            currentNode.instructions.push([...this.#instruction]);
            this.#instruction = [];
            return currentNode;
        }
        if (this.#sym.tag === Tag.IF) {
            let condition;
            if (currentNode.instructions.length) {
                currentNode.children.push(this.newNode(currentNode));
                condition = currentNode.children[0];
            } else {
                condition = currentNode;
            }
            this.getNextToken();
            this.parseExpr();
            condition.instructions.push([...this.#instruction]);
            this.#instruction = [];
            if (this.#sym.tag !== Tag.THEN) {
                return;
            }
            this.getNextToken();
            condition.children.push(this.newNode(condition));
            let exit = condition.children[0];
            while (this.#sym.tag !== Tag.ELSE && this.#sym.tag !== Tag.ENDIF) {
                exit = this.parseAxiom(exit);
            }
            if (exit === condition.children[0]) {
                exit = this.newNode(condition.children[0]);
                condition.children[0].children.push(exit);
            }
            if (this.#sym.tag === Tag.ELSE) {
                this.getNextToken();
                condition.children.push(this.newNode(condition));
                let elseExit = condition.children[1];
                while (this.#sym.tag !== Tag.ENDIF) {
                    elseExit = this.parseAxiom(elseExit);
                }
                if (elseExit !== condition.children[1]) {
                    elseExit.parents.forEach((parent) => {
                        parent.children.splice(parent.children.indexOf(elseExit), 1);
                        parent.children.push(exit);
                        exit.parents.push(parent);
                    });
                } else {
                    condition.children[1].children.push(exit);
                    exit.parents.push(condition.children[1]);
                }
            } else {
                condition.children.push(exit);
                exit.parents.push(condition);
            }
            this.getNextToken();
            return exit;
        }
        if (this.#sym.tag === Tag.WHILE) {
            let condition;
            if (currentNode.instructions.length) {
                currentNode.children.push(this.newNode(currentNode));
                condition = currentNode.children[0];
            } else {
                condition = currentNode;
            }
            this.getNextToken();
            this.parseExpr();
            condition.instructions.push([...this.#instruction]);
            this.#instruction = [];
            if (this.#sym.tag !== Tag.DO) {
                return;
            }
            this.getNextToken();
            condition.children.push(this.newNode(condition));
            let exit = condition.children[0];
            while (this.#sym.tag !== Tag.ENDWHILE) {
                exit = this.parseAxiom(exit);
            }
            if (exit === condition.children[0]) {
                condition.children[0].children.push(condition);
                exit = this.newNode();
            } else {
                exit.parents.forEach((parent) => {
                    parent.children.splice(parent.children.indexOf(exit), 1);
                    parent.children.push(condition);
                    condition.parents.push(parent);
                });
            }
            condition.children.push(exit);
            exit.parents.push(condition);
            this.getNextToken();
            return exit;
        }
    }
    parseExpr() {
        this.parsePrimary();
        this.parseBinOpRHS();
    }
    parsePrimary() {
        if (this.#sym.tag === Tag.ID || this.#sym.tag === Tag.NUMBER) {
            this.#instruction.push(this.#sym);
            this.getNextToken();
        } else if (this.#sym.tag === Tag.LEFTBRACKET) {
            this.#instruction.push(this.#sym);
            this.getNextToken();
            this.parseExpr();
            if (this.#sym.tag !== Tag.RIGHTBRACKET) {
                return;
            }
            this.getNextToken();
        }
    }
    parseBinOpRHS() {
        while (binOpTags.includes(this.#sym.tag)) {
            this.#instruction.push(this.#sym);
            this.getNextToken();
            this.parsePrimary();
        }
    }
    outputCFG() {
        const queue = [ this.#tree ];
        while (queue.length !== 0) {
            const node = queue.shift();
            console.log({
                id: node.id,
                instructions: node.instructions.map((i) => i.map((c) => {
                    return c.tag + (c.id ? ' ' + c.id : '') + (c.number ? ' ' + c.number : '');
                })),
                parents: node.parents.map((p) => p.id),
                children: node.children.map((c) => c.id),
            });
            if (node.children.length) {
                queue.push(...node.children.filter((child) =>
                    child.id > node.id && !queue.includes(child)
                ));
            }
        }
    }
    outputCFGDOT() {
        let ret = 'digraph {\n';
        const queue = [ this.#tree ];
        while (queue.length !== 0) {
            const node = queue.shift();
            ret += `    ${node.id} [label="${node.instructions.map((i) => i.map((c) => {
                return (c.id ? c.id : (c.number !== undefined ? c.number : c.tag));
            }).join(' ')).join('\\n')}"];\n`;
            node.children.forEach((child) => ret += `    ${node.id} -> ${child.id};\n`);
            if (node.children.length) {
                queue.push(...node.children.filter((child) =>
                    child.id > node.id && !queue.includes(child)
                ));
            }
        }
        return ret + '}';
    }
    toSSA() {;
        const varVersions = {};
        this.#declared.forEach((varName) => varVersions[varName] = 0);
        const queue = [ this.#tree ];
        while (queue.length !== 0) {
            const node = queue.shift();
            const parentVarVersionsArray = node.parents.map(({ varVersions }) => varVersions);
            const parentVarVersions = {};
            parentVarVersionsArray.forEach((versionsObj) => {
                for (const [key, value] of Object.entries(versionsObj)) {
                    if (parentVarVersions[key] !== undefined && parentVarVersions[key] !== value) {
                        if (typeof parentVarVersions[key] === 'number') {
                            parentVarVersions[key] = `phi_${parentVarVersions[key]}_${value}`;
                        } else {
                            parentVarVersions[key] = parentVarVersions[key].substring(0, parentVarVersions[key].length - 1) + `, ${value})`;
                        }
                    } else {
                        parentVarVersions[key] = value;
                    }
                }
            });
            if (parentVarVersionsArray.length === 0) {
                Object.assign(parentVarVersions, varVersions);
            }
            for (let i = 0; i < node.instructions.length; ++i) {
                const instruction = node.instructions[i];
                if (instruction[0].tag === Tag.VAR) {
                    instruction[1].id = instruction[1].id + '__' + (++varVersions[instruction[1].id]);
                }
                const assignIndex = instruction.findIndex(({ tag }) => tag === Tag.ASSIGN);
                if (assignIndex !== -1 || (assignIndex === -1 && instruction[0].tag !== Tag.VAR)) {
                    for (let j = (assignIndex === -1 ? 0 : assignIndex + 1); j < instruction.length; ++j) {
                        if (instruction[j].tag === Tag.ID) {
                            instruction[j].id = instruction[j].id + '__' + (parentVarVersions[instruction[j].id]);
                        }
                    }
                }
                if (instruction[1].tag === Tag.ASSIGN) {
                    const varName = instruction[0].id + '__' + (varVersions[instruction[0].id]++);
                    if (varVersions[instruction[0].id] > 2) {
                        node.instructions.splice(i, 0, [
                            { tag: Tag.VAR },
                            { tag: Tag.ID, id: varName },
                            { tag: Tag.SEMICOLON },
                        ]);
                        ++i;
                    }
                    instruction[0].id = varName;
                }
            };
            node.varVersions = Object.assign({}, varVersions);
            if (node.children.length) {
                queue.push(...node.children.filter((child) =>
                    child.id > node.id && !queue.includes(child)
                ));
            }
        }
    }
}
