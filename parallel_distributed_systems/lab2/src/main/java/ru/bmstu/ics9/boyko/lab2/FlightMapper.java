package ru.bmstu.ics9.boyko.lab2;

import com.sun.org.apache.bcel.internal.generic.NEW;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class FlightMapper extends Mapper <LongWritable, Text, DelayWritableComparable, IntWritable>{
    private static final String NEW_LINE  = "\n";
    private static final String COMMA = ",";

    private static final int DEST_AIRPORT_ID = 14;
    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
        String[] lines = value.toString().split(NEW_LINE);
        for (String line : lines) {
            String[] parsedFlight = line.split(COMMA);
            Integer destAirportID =  Integer.parseInt(parsedFlight[DEST_AIRPORT_ID]);
            Float delay =  Float.parseFloat(parsedFlight[DEST_AIRPORT_ID]);
            
        }
    }
}
