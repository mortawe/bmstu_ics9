package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Partitioner;

public class AirportPartitioner extends Partitioner<DelayWritableComparable, Text> {
    @Override
    public int getPartition(DelayWritableComparable delayWritableComparable, Text text, int numPartitions) {
        return (delayWritableComparable.airportID & Integer.MAX_VALUE) % numPartitions;
    }
}
