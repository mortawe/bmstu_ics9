package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class FlightMapper extends Mapper<LongWritable, Text, FlightWC, IntWritable> {
    public static final String NEW_LINE = "\n";
    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException,
            InterruptedException {
        String[] lines = value.toString().split(NEW_LINE);
        for (String line : lines) {
            FlightWC flight = new FlightWC(line);
            if (!flight.cancelled && flight.delay_time > 0) {
                context.write(flight, new IntWritable(1));
            }
        }
    }
}
