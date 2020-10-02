import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class WordMapper extends Mapper<LongWritable, Text, Text, IntWritable> {
    @Override
    protected void map(LongWritable key, Text value, Context context) throws IOException,
            InterruptedException {
        String refactoredString = value.toString().toLowerCase()
                .replaceAll("[^a-zA-Z0-9а-яА-Я]", "");
        String[] words = refactoredString.split(" ");

        for (String word : words) {
            context.write(new Text(word), new IntWritable(1));
        }
        return;
    }
}
