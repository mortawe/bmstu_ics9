import org.apache.spark.SparkConf;
import org.apache.spark.api.java.JavaPairRDD;
import org.apache.spark.api.java.JavaRDD;

public class SparkApp {
    SparkConf conf = new SparkConf().setAppName("lab3");
    JavaSparkContext sc = new JavaSparkContext(conf);

    JavaRDD<String> flightsFile = sc.textFile("664600583_T_ONTIME_sample.csv");
    JavaRDD<String> airportsFile = sc.textFile("L_AIRPORT_ID.csv");

}
