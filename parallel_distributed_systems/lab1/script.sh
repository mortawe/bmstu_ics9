#bin/bash
rm -r output

hadoop fs -rm -r output

mvn package

#hadoop fs -copyFromLocal warandpeace1.txt

export HADOOP_CLASSPATH=target/hadoop-examples-1.0-SNAPSHOT.jar

hadoop WordCountApp warandpeace1.txt output

hadoop fs -copyToLocal output

