package ru.bmstu.ics9.boyko.lab2;


import org.apache.hadoop.mapreduce.Partitioner;

public class FlightPartitioner <K, V> extends Partitioner<K, V> {
    public int getPartition(K key, V value, int numReduceTasks) {
        return (key.hashCode() & Integer.MAX_VALUE) % numReduceTasks;
    }
}