package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.RawComparator;
import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class DelayGroupingComparator extends WritableComparator {
    @Override
    public int compare(WritableComparable a, WritableComparable b) {
        DelayWritableComparable left = (DelayWritableComparable)a;
        return super.compare(a, b);
    }
}
