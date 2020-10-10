import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;


public class WordMapper extends Mapper<LongWritable, Text, Text, IntWritable> {
    public static final String REGEX_NO_SYMBOLS = "[^a-zA-Z0-9а-яА-Я]";
    public static final String SPACE = " ";

    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException,
            InterruptedException {
        String[] words = value.toString().replaceAll(REGEX_NO_SYMBOLS, SPACE)
                .toLowerCase().split(SPACE);
        for (String word : words) {
            context.write(new Text(word), new IntWritable(1));
        }
    }
}