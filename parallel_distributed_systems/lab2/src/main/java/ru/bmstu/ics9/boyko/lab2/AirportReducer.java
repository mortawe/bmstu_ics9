package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

import java.io.IOException;

public class AirportReducer extends Reducer<DelayWritableComparable, Text, Text, Text> {
    @Override
    protected void reduce(DelayWritableComparable key, Iterable<Text> values, Context context) throws IOException, InterruptedException {
        for (Text value: values) {
            
        }
    }
}
