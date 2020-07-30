# partialkey

Implementation of the Partial Key Grouping Load-Balanced Partitioning of Distributed Streams algorithm as described on

https://arxiv.org/abs/1510.07623

https://melmeric.files.wordpress.com/2014/11/the-power-of-both-choices-practical-load-balancing-for-distributed-stream-processing-engines.pdf


Optionally it includes a small variation to make the balancing take into account the in-flight processing of tasks instead of the total allocated in order to prevent contention on slots due to processing time of tasks.
