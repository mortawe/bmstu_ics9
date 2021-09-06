package ru.bmstu.iu9;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;


public class IdentVsNumber {

    public static void main(String[] args) {
        if (args.length != 3) {
            System.exit(1);
        } else {
            Position pos = new Position(1, 1);
            List<String> lines;
            String file_name = args[2];
            try {
                lines = Files.readAllLines(Paths.get(file_name), StandardCharsets.UTF_8);
                for (String line : lines) {
                    test_match(line, pos);
                    pos.next_line();
                }
            } catch (IOException ex) {
                ex.printStackTrace();
            }
        }
    }

    private static void test_match(String line, Position pos) {
        String ident = "[A-z][A-z]+|\\([0-9]+\\)";
        String number = "[1-9][0-9]*|0";
        String operators = "(:=|\\(\\)|:)";
        String pattern = "(?<ident>^" + ident + ")|(?<number>^" + number + ")|(?<operator>^" + operators + ")";

        Pattern p = Pattern.compile(pattern);
        Matcher m;

        while (!line.equals("")) {
            m = p.matcher(line);
            if (m.find()) {
                if (m.group("ident") != null) {
                    String item = m.group("ident");
                    System.out.println("IDENT " + pos.toString() + ": " + item);
                    pos.next_pos(item.length());
                    line = line.substring(line.indexOf(item) + item.length());
                } else if (m.group("number") != null) {
                    String item = m.group("number");
                    System.out.println("NUMBER " + pos.toString() + ": " + item);
                    pos.next_pos(item.length());
                    line = line.substring(line.indexOf(item) + item.length());
                } else {
                    String item = m.group("operator");
                    System.out.println("OPERATOR " + pos.toString() + ": " + item);
                    pos.next_pos(item.length());
                    line = line.substring(line.indexOf(item) + item.length());
                }
            } else {
                if (Character.isWhitespace(line.charAt(0))) {
                    while (Character.isWhitespace(line.charAt(0))) {
                        line = line.substring(1);
                        pos.next_pos(1);
                    }
                } else {
                    System.out.println("syntax error " + pos.toString());
                    while (!m.find() && !line.equals("")) {
                        line = line.substring(1);
                        pos.next_pos(1);
                        m = p.matcher(line);
                    }
                }
            }
        }
    }
}