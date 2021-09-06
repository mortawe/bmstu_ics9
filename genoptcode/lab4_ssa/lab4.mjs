import { Compiler } from './compiler.mjs';
import fs from 'fs';
import {exec} from "child_process";


fs.readFile('./input.lang', 'utf8' , (err, data) => {
    if (err) {
        console.error(err);
        return;
    }
    main(data);
});

const main = (program) => {
    const compiler = new Compiler();
    compiler.parse(program);
    console.log("BEFORE: \n")
    compiler.outputCFG();
    fs.writeFile('CFG1.dot', compiler.outputCFGDOT(), (err) => {});
    const child = exec('dot -Tpng CFG1.dot > output1.png',
        (error, stdout, stderr) => {
            console.log(`stdout: ${stdout}`);
            console.log(`stderr: ${stderr}`);
            if (error !== null) {
                console.log(`exec error: ${error}`);
            }
        });
    compiler.toSSA();
    console.log("AFTER: \n")
    compiler.outputCFG();
    fs.writeFile('CFG2.dot', compiler.outputCFGDOT(), (err) => {});


    const child2 = exec('dot -Tpng CFG2.dot > output2.png',
        (error, stdout, stderr) => {
            console.log(`stdout: ${stdout}`);
            console.log(`stderr: ${stderr}`);
            if (error !== null) {
                console.log(`exec error: ${error}`);
            }
        });
}
