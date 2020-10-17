package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.WritableComparable;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;
import java.util.Objects;

// :D
public class AirportWC implements WritableComparable {
    public static final int DEST_ROW_NUM = 14;
    public static final int DELAY_ROW_NUM = 17;
    public static final int CANCELLED_ROW = 19;
    public static final String COMMA = ",";

    int dest_id;
    float delay_time;
    boolean cancelled;

    public AirportWC(String value) {
        String[] rows = value.split(COMMA);
        this.dest_id = Integer.parseInt(rows[DEST_ROW_NUM]);
        // negative delay is not delay at all
        // -1 is value for these flights
        this.delay_time = Float.parseFloat(rows[DELAY_ROW_NUM]) > 0 ? Float.parseFloat(rows[DELAY_ROW_NUM]) : -1;
        this.cancelled = Boolean.parseBoolean(rows[CANCELLED_ROW]);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        AirportWC airportWC = (AirportWC) o;
        return dest_id == airportWC.dest_id &&
                Float.compare(airportWC.delay_time, delay_time) == 0 &&
                cancelled == airportWC.cancelled;
    }

    @Override
    public int hashCode() {
        return Objects.hash(dest_id, delay_time, cancelled);
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
