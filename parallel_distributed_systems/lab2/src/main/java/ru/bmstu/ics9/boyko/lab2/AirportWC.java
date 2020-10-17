package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.WritableComparable;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;
// :D
public class AirportWC implements WritableComparable {
    public static final DEST_ROW_NUM =
    int dest_id;
    double delay_time;

    public AirportWC(String value){
        String[] rows = value.split(",");
        this.dest_id = rows[]
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
