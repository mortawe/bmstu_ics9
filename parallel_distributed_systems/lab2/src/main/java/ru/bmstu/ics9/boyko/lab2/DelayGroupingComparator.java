package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class DelayGroupingComparator extends WritableComparator {
    @Override
    public int compare(Object a, Object b) {
        DelayWritableComparable left = (DelayWritableComparable) a;
        DelayWritableComparable right = (DelayWritableComparable) b;

        return left.airportID - right.airportID;
    }
}
