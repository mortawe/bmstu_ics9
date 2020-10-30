package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

import java.io.IOException;

public class AirportReducer extends Reducer<DelayWritableComparable, Text, Text, Text> {
    private static final int DEFAULT_COUNTERS_VALUE = -1;
    @Override
    protected void reduce(DelayWritableComparable key, Iterable<Text> values, Context context) throws IOException, InterruptedException {
        float max = DEFAULT_COUNTERS_VALUE, min = DEFAULT_COUNTERS_VALUE;
        boolean isFirstLine = true;
        String airportName = "";
        for (Text value: values) {
            if (isFirstLine) {
                airportName = value.toString();
                continue;
            }
            float delay = Float.parseFloat(value.toString());
            if (max == DEFAULT_COUNTERS_VALUE || max < delay) {
                max = delay;
            }
            if (min == DEFAULT_COUNTERS_VALUE || min > delay) {
                min = delay;
            }
        }
    }
}
