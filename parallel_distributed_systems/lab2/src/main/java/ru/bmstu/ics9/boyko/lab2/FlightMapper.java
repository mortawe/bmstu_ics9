package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class FlightMapper extends Mapper <LongWritable, Integer, DelayWritableComparable, Integer>{
    @Override
    protected void map(LongWritable key, Integer value, Context context) throws IOException, InterruptedException {
        super.map(key, value, context);
    }
}
