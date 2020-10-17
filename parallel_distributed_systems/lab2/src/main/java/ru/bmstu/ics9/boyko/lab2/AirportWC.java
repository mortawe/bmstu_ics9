package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.WritableComparable;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

// :D
public class AirportWC implements WritableComparable {
    public static final int DEST_ROW_NUM = 14;
    public static final int DELAY_ROW_NUM = 17;
    public static final String COMMA = ",";

    int dest_id;
    float delay_time;

    public AirportWC(String value) {
        String[] rows = value.split(COMMA);
        this.dest_id = Integer.parseInt(rows[DEST_ROW_NUM]);
        // negative delay is not delay at all
        // -1 is value for these flights
        this.delay_time = Float.parseFloat(rows[DELAY_ROW_NUM]) > 0 ? Float.parseFloat(rows[DELAY_ROW_NUM]) : -1;
    }

    @Override
    public int compareTo(Object o) {
        AirportWC right = (AirportWC) o;
        if
        return 0;
    }

    @Override
    public void write(DataOutput dataOutput) throws IOException {

    }

    @Override
    public void readFields(DataInput dataInput) throws IOException {

    }
}
