package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.RawComparator;
import org.apache.hadoop.io.WritableComparator;

public class DelayGroupingComparator extends WritableComparator {
    @Override
    public int compare(byte[] bytes, int i, int i1, byte[] bytes1, int i2, int i3) {
        return 0;
    }
}
