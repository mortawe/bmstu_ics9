package ru.bmstu.ics9.boyko.lab2;

import org.apache.hadoop.io.Writable;
import org.apache.hadoop.util.IdentityHashStore;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

public class AirportWritable implements Writable {
    private int code;
    private String decription;
    @Override
    public void write(DataOutput dataOutput) throws IOException {
        dataOutput.writeInt(code);
        dataOutput.writeChars(decription);
    }

    @Override
    public void readFields(DataInput dataInput) throws IOException {
        code = dataInput.readInt();
        decription = dataInput.readLine();
    }
}
