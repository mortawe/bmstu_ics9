package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class FlightMapper extends Mapper<LongWritable, Text, DelayWritableComparable, Text> {
    private static final String NEW_LINE = "\n";
    private static final String COMMA = ",";
    private static final String NOT_NUMBERS_REGEX = "[^0-9]+";
    private static final String EMPTY_STRING = "";

    private static final int DEST_AIRPORT_ID_POS = 14;
    private static final int DELAY_POS = 17;

    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
        String[] lines = value.toString().split(NEW_LINE);
        for (String line : lines) {
            String[] parsedFlight = line.split(COMMA);
            int destAirportID = Integer.parseInt(parsedFlight[DEST_AIRPORT_ID_POS].replaceAll(NOT_NUMBERS_REGEX, EMPTY_STRING));
            String delayText = parsedFlight[DELAY_POS];
            if (destAirportID <= 0) {
                // not a valid string
                continue;
            }
            float delayFloat = Float.parseFloat(delayText);
            if (delayFloat <= 0) {
                continue;
            }
            context.write(new DelayWritableComparable(destAirportID, DelayWritableComparable.STATE_FLIGHT), new Text(delayText));
        }
    }
}
