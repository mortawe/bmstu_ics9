mvn package

export HADOOP_CLASSPATH=target/hadoop-examples-1.0-SNAPSHOT.jar

hadoop WordCountApp warandpeace1.txt output

hadoop fs -copyToLocal output

hadoop fs -rm -r output