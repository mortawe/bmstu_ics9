

public class SparkApp {
    SparkConf conf = new SparkConf().setAppName("lab5");
    JavaSparkContext sc = new JavaSparkContext(conf);

    JavaRDD<String> flightsFile = sc.textFile("664600583_T_ONTIME_sample.csv");
    JavaRDD<String> airportsFile = sc.textFile("L_AIRPORT_ID.csv");

}
