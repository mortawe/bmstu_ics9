package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.mapreduce.Reducer;

import java.io.IOException;

public class FlightReducer extends Reducer <FlightWC, IntWritable, FlightWC, IntWritable> {
    @Override
    protected void reduce(FlightWC key, Iterable<IntWritable> values, Context context) throws IOException, InterruptedException {

    }
}
