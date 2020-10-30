package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

import java.io.IOException;

public class AirportReducer extends Reducer<DelayWritableComparable, Text, Text, Text> {
    private static final float DEFAULT_COUNTERS_VALUE = -1;
    @Override
    protected void reduce(DelayWritableComparable key, Iterable<Text> values, Context context) throws IOException, InterruptedException {
        float max = DEFAULT_COUNTERS_VALUE, min = DEFAULT_COUNTERS_VALUE, avg = 0;
        int flightCount = 0;
        boolean isFirstLine = true;
        Text airportName = null;
        for (Text value: values) {
            if (isFirstLine) {
                airportName = value;
                continue;
            }
            float delay = Float.parseFloat(value.toString());
            if (max == DEFAULT_COUNTERS_VALUE || max < delay) {
                max = delay;
            }
            if (min == DEFAULT_COUNTERS_VALUE || min > delay) {
                min = delay;
            }
            avg += delay;
            flightCount++;
        }
        avg /= flightCount;

        String results = String.format("min : %f; max : %f; avg: %f", min, max, avg);

        context.write(airportName, new Text(results));
    }
}
