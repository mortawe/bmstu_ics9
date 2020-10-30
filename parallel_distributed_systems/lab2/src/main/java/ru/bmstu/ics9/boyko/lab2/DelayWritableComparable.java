package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.WritableComparable;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;



public class DelayWritableComparable implements WritableComparable {
    public static final Integer STATE_AIRPORT = 0;
    public static final Integer STATE_FLIGHT = 1;

    int airportID;
    int state; // 0 for airports; 1 for flights

    public DelayWritableComparable(int airportID, int state) {
        this.airportID = airportID;
        this.state = state;
    }

    @Override
    public int compareTo(Object o) {
        return 0;
    }

    @Override
    public void write(DataOutput dataOutput) throws IOException {

    }

    @Override
    public void readFields(DataInput dataInput) throws IOException {

    }
}
