package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class AirportMapper extends Mapper<LongWritable, Text, DelayWritableComparable, Text> {
    private static final String NEW_LINE = "\n";
    private static final String COMMA = ",";
    private static final String NOT_NUMBERS_REGEX = "[^0-9]+";
    private static final String EMPTY_STRING = "";

    private static final int CODE_POS = 0;
    private static final int DESCRIPTION_POS = 1;

    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
        String[] lines = value.toString().split(NEW_LINE);

        for (String line : lines) {
            String[] parsedAirport = line.split(COMMA, 2);
            String codeStr = parsedAirport[CODE_POS].replaceAll(NOT_NUMBERS_REGEX, EMPTY_STRING);
            if (codeStr.isEmpty()) {
                return;
            }
            int code = Integer.parseInt(codeStr);
            String description = parsedAirport[DESCRIPTION_POS];
            if (code <= 0) {
                continue;
            }
            System.out.println("name " + description + " id " + codeStr + "\n");

            context.write(new DelayWritableComparable(code, DelayWritableComparable.STATE_AIRPORT),
                    new Text(description));
        }

    }
}
