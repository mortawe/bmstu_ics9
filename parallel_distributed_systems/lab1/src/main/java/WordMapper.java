import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class WordMapper extends Mapper<LongWritable, Text, Text, IntWritable> {
    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException,
            InterruptedException {
        String line = value.toString();

        line = line.toLowerCase();

        String[] words = line.split("[^a-zA-Z]"); //не забыть русский текст через а-яА-Я

        for (String word : words) {
            context.write(new Text(word), new IntWritable(1));
        }
    }
}