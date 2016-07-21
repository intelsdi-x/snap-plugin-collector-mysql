# snap collector plugin - mysql

## Collected Metrics
This plugin has the ability to gather the following metrics:

Counter data type is used for values which increment continuously, it's always positive.
Derive data type is used to represent value changed in time.
Gauge data simply returns stored value.

Namespace | Type |Description
----------|------|------------------
/intel/mysql/bytes/buffer_pool_size |gauge|The number of row locks currently being waited for (innodb_row_lock_current_waits).
/intel/mysql/bytes/ibuf_size |gauge| The Number of row locks currently being waited for (innodb_row_lock_current_waits).
/intel/mysql/bytes/metadata_mem_pool_size |gauge| The Size of a memory pool InnoDB uses to store data dictionary and internal data structures.
/intel/mysql/cache_result/qcache-hits |derive|  The number of query cache hits.
/intel/mysql/cache_result/qcache-inserts |derive| The number of queries added to the query cache.
/intel/mysql/cache_result/qcache-not_cached |derive| The number of noncached queries (not cacheable, or not cached due to the query_cache_type setting).
/intel/mysql/cache_result/qcache-prunes |derive| The number of queries that were deleted from the query cache because of low memory.
/intel/mysql/cache_size/qcache |gauge| The number of queries registered in the query cache.
/intel/mysql/gauge/buffer_pool_bytes_data |gauge| Buffer bytes containing data (innodb_buffer_pool_bytes_data).
/intel/mysql/gauge/buffer_pool_bytes_dirty |gauge| Buffer bytes currently dirty (innodb_buffer_pool_bytes_dirty).
/intel/mysql/gauge/buffer_pool_pages_data |gauge| Buffer pages containing data (innodb_buffer_pool_pages_data).
/intel/mysql/gauge/buffer_pool_pages_dirty |gauge| Buffer pages containing data (innodb_buffer_pool_pages_data)
/intel/mysql/gauge/buffer_pool_pages_free |gauge|  Buffer pages currently free (innodb_buffer_pool_pages_free).
/intel/mysql/gauge/buffer_pool_pages_misc |gauge| Buffer pages for misc use such as row locks or the adaptive hash index (innodb_buffer_pool_pages_misc).
/intel/mysql/gauge/buffer_pool_pages_total |gauge| Total buffer pool size in pages (innodb_buffer_pool_pages_total).
/intel/mysql/gauge/file_num_open_files |gauge|  The number of files currently open (innodb_num_open_files).
/intel/mysql/gauge/innodb_activity_count |gauge| Th number of files currently open (innodb_num_open_files).
/intel/mysql/gauge/innodb_dblwr_page_size |gauge| InnoDB page size in bytes (innodb_page_size).
/intel/mysql/gauge/trx_rseg_history_len |gauge| The length of the TRX_RSEG_HISTORY list.
/intel/mysql/mysql_bpool_bytes/data |gauge| The total number of bytes in the InnoDB buffer pool containing data. The number includes both dirty and clean pages.
/intel/mysql/mysql_bpool_bytes/dirty |gauge| The total current number of bytes held in dirty pages in the InnoDB buffer pool.
/intel/mysql/mysql_bpool_counters/pages_flushed |counter| The number of requests to flush pages from the InnoDB buffer pool.
/intel/mysql/mysql_bpool_counters/read_ahead |counter| The number of pages read into the InnoDB buffer pool by the read-ahead background thread.
/intel/mysql/mysql_bpool_counters/read_ahead_evicted |counter| The number of pages read into the InnoDB buffer pool by the read-ahead background thread that were subsequently evicted without having been accessed by queries.
/intel/mysql/mysql_bpool_counters/read_ahead_rnd |counter| The number of “random” read-aheads initiated by InnoDB. This happens when a query scans a large portion of a table but in random order.
/intel/mysql/mysql_bpool_counters/read_requests |counter| The number of logical read requests. 
/intel/mysql/mysql_bpool_counters/reads |counter| he number of logical reads that InnoDB could not satisfy from the buffer pool, and had to read directly from disk.
/intel/mysql/mysql_bpool_counters/write_requests |counter| The number of writes done to the InnoDB buffer pool.
/intel/mysql/mysql_bpool_pages/data |gauge| The number of pages in the InnoDB buffer pool containing data. The number includes both dirty and clean pages.
/intel/mysql/mysql_bpool_pages/dirty |gauge| The total current number of bytes held in dirty pages in the InnoDB buffer pool.
/intel/mysql/mysql_bpool_pages/free |gauge| The number of free pages in the InnoDB buffer pool.
/intel/mysql/mysql_bpool_pages/misc |gauge| The number of pages in the InnoDB buffer pool that are busy because they have been allocated for administrative overhead, such as row locks or the adaptive hash index.
/intel/mysql/mysql_bpool_pages/total |gauge| The total size of the InnoDB buffer pool, in pages.
/intel/mysql/mysql_innodb_data/fsyncs |counter| The number of fsync() operations so far.
/intel/mysql/mysql_innodb_data/read |counter| The amount of data read since the server was started.
/intel/mysql/mysql_innodb_data/reads |counter| The total number of data reads.
/intel/mysql/mysql_innodb_data/writes |counter| The total number of data writes.
/intel/mysql/mysql_innodb_data/written |counter| The amount of data written so far, in bytes.
/intel/mysql/mysql_innodb_dblwr/writes |counter| The number of doublewrite operations that have been performed.
/intel/mysql/mysql_innodb_dblwr/written |counter| The number of pages that have been written to the doublewrite buffer.
/intel/mysql/mysql_innodb_log/fsyncs |counter| The number of fsync() writes done to the InnoDB redo log files.
/intel/mysql/mysql_innodb_log/waits |counter| The number of times that the log buffer was too small and a wait was required for it to be flushed before continuing.
/intel/mysql/mysql_innodb_log/write_requests |counter| The number of write requests for the InnoDB redo log.
/intel/mysql/mysql_innodb_log/writes |counter| The number of physical writes to the InnoDB redo log file.
/intel/mysql/mysql_innodb_log/written |counter| The number of bytes written to the InnoDB redo log files.
/intel/mysql/mysql_innodb_pages/created |counter|  The number of pages created by operations on InnoDB tables.
/intel/mysql/mysql_innodb_pages/read |counter| The number of pages read by operations on InnoDB tables.
/intel/mysql/mysql_innodb_pages/written |counter| The number of pages written by operations on InnoDB tables.
/intel/mysql/mysql_innodb_row_lock/time |counter| The total time spent in acquiring row locks for InnoDB tables, in milliseconds.
/intel/mysql/mysql_innodb_row_lock/waits |counter| The number of times operations on InnoDB tables had to wait for a row lock.
/intel/mysql/mysql_innodb_rows/deleted |counter| The number of rows deleted from InnoDB tables.
/intel/mysql/mysql_innodb_rows/inserted |counter| The number of rows inserted into InnoDB tables.
/intel/mysql/mysql_innodb_rows/read |counter| The number of rows read from InnoDB tables.
/intel/mysql/mysql_innodb_rows/updated |counter| The number of rows updated in InnoDB tables.
/intel/mysql/mysql_locks/lock_deadlocks |counter| The number of deadlocks.
/intel/mysql/mysql_locks/lock_row_lock_current_waits |derive| The number of row locks currently being waited for (innodb_row_lock_current_waits).
/intel/mysql/mysql_locks/lock_timeouts |derive| The number of row locks currently being waited for (innodb_row_lock_current_waits).
/intel/mysql/mysql_log_position/master-bin |counter| The position of  the binary log file of the master.
/intel/mysql/mysql_log_position/slave-exec |counter|  The position in the current master binary log file to which the SQL thread has read and executed, marking the start of the next transaction or event to be processed. 
/intel/mysql/mysql_log_position/slave-read |counter| The position in the current master binary log file up to which the I/O thread has read. 
/intel/mysql/mysql_log_position/time_offset |gauge|  This field is an indication of how “late” the slave is when the slave is actively processing updates, this field shows the difference between the current timestamp on the slave and the original timestamp logged on the master for the event currently being processed on the slave or when no event is currently being processed on the slave, this value is 0. 
/intel/mysql/mysql_octets/rx |gauge| The number of bytes received from all clients.
/intel/mysql/mysql_octets/tx |gauge| The number of bytes sent to all clients.
/intel/mysql/operations/adaptive_hash_searches |derive| The number of successful searches using Adaptive Hash Index.
/intel/mysql/operations/buffer_data_reads |derive| The amount of data read in bytes (innodb_data_reads).
/intel/mysql/operations/buffer_data_written |derive| The amount of data written in bytes (innodb_data_written).
/intel/mysql/operations/buffer_pages_created |derive| Number of pages created (innodb_pages_created)
/intel/mysql/operations/buffer_pages_read |derive|  The number of pages read (innodb_pages_read).
/intel/mysql/operations/buffer_pages_written |derive| The amount of data written in bytes (innodb_data_written).
/intel/mysql/operations/buffer_pool_read_ahead |counter| The number of pages read as read ahead (innodb_buffer_pool_read_ahead).
/intel/mysql/operations/buffer_pool_read_ahead_evicted |counter|  Read-ahead pages evicted without being accessed (innodb_buffer_pool_read_ahead_evicted).
/intel/mysql/operations/buffer_pool_read_requests |counter| The number of logical read requests (innodb_buffer_pool_read_requests).
/intel/mysql/operations/buffer_pool_reads |counter|  The number of reads directly from disk (innodb_buffer_pool_reads).
/intel/mysql/operations/buffer_pool_wait_free |derive| The number of times waited for free buffer (innodb_buffer_pool_wait_free).
/intel/mysql/operations/buffer_pool_write_requests |counter| The number of write requests (innodb_buffer_pool_write_requests).
/intel/mysql/operations/dml_deletes |derive| The number of rows deleted.
/intel/mysql/operations/dml_inserts |derive| The number of rows inserted.
/intel/mysql/operations/dml_reads |derive| The number of rows read.
/intel/mysql/operations/dml_updates |serive| The number of rows updated.
/intel/mysql/operations/ibuf_merges_delete |derive| The number of purge records merged by change buffering.
/intel/mysql/operations/ibuf_merges_delete_mark |derive|  The number of deleted records merged by change buffering.      
/intel/mysql/operations/ibuf_merges_discard_delete |derive| The number of purge merged  operations discarded.
/intel/mysql/operations/ibuf_merges_discard_delete_mark |derive| The number of deleted merged operations discarded. 
/intel/mysql/operations/ibuf_merges_discard_insert |derive| The number of insert merged operations discarded.
/intel/mysql/operations/ibuf_merges_discard_merges |derive|
/intel/mysql/operations/ibuf_merges_insert |derive|  The number of inserted records merged by change buffering.  
/intel/mysql/operations/innodb_dblwr_pages_written |derive|  The number of pages that have been written for doublewrite operations (innodb_dblwr_pages_written).
/intel/mysql/operations/innodb_dblwr_writes |counter| The number of doublewrite operations that have been performed (innodb_dblwr_writes).
/intel/mysql/operations/innodb_rwlock_s_os_waits |derive| The number of OS waits due to shared latch request.
/intel/mysql/operations/innodb_rwlock_s_spin_rounds |derive|  The number of rwlock spin loop rounds due to shared latch request.
/intel/mysql/operations/innodb_rwlock_s_spin_waits |derive| The number of rwlock spin waits due to shared latch request.
/intel/mysql/operations/innodb_rwlock_x_os_waits |derive|  The number of OS waits due to exclusive latch request.
/intel/mysql/operations/innodb_rwlock_x_spin_rounds |derive| The number of rwlock spin loop rounds due to exclusive latch request.
/intel/mysql/operations/innodb_rwlock_x_spin_waits |derive| The number of rwlock spin waits due to exclusive latch request.
/intel/mysql/operations/log_waits |counter| The number of log waits due to small log buffer (innodb_log_waits).
/intel/mysql/operations/log_write_requests |counter| The number of log write requests (innodb_log_write_requests).   
/intel/mysql/operations/log_writes |counter| The number of log writes (innodb_log_writes).
/intel/mysql/operations/os_data_fsyncs |derive| The number of fsync() calls (innodb_data_fsyncs).
/intel/mysql/operations/os_data_reads |derive|The number of reads initiated (innodb_data_reads).  
/intel/mysql/operations/os_data_writes |derive| The number of writes initiated (innodb_data_writes).
/intel/mysql/operations/os_log_bytes_written |derive| Bytes of log written (innodb_os_log_written).
/intel/mysql/operations/os_log_fsyncs |derive| The number of fsync log writes (innodb_os_log_fsyncs). 
/intel/mysql/operations/os_log_pending_fsyncs |derive| The number of pending fsync write (innodb_os_log_pending_fsyncs).
/intel/mysql/operations/os_log_pending_writes |derive| The number of pending log file writes (innodb_os_log_pending_writes).
/intel/mysql/threads/cached |gauge| The number of threads in the thread cache. 
/intel/mysql/threads/connected |gauge| The number of currently open connections. 
/intel/mysql/threads/running |gauge|  The number of threads that are not sleeping.
/intel/mysql/total_threads/created |derive| The number of threads created to handle connections.
/intel/mysql/mysql_locks/immediate |counter| The number of times that a request for a table lock could be granted immediately.
/intel/mysql/mysql_locks/waited |counter| The number of times that a request for a table lock could not be granted immediately and a wait was needed.
/intel/mysql/mysql_select/full_join |counter| The number of joins that perform table scans because they do not use indexes. If this value is not 0, you should carefully check the indexes of your tables.
/intel/mysql/mysql_select/full_range_join |counter| The number of joins that used a range search on a reference table.
/intel/mysql/mysql_select/range |counter| The number of joins that used ranges on the first table. 
/intel/mysql/mysql_select/range_check |counter| The number of joins without keys that check for key usage after each row. If this is not 0, you should carefully check the indexes of your tables.
/intel/mysql/mysql_select/scan |counter| The number of joins that did a full scan of the first table.
/intel/mysql/mysql_sort/merge_passes |counter| The number of merge passes that the sort algorithm has had to do.
/intel/mysql/mysql_sort/range |counter| The number of sorts that were done using ranges.
/intel/mysql/mysql_sort/rows |counter| The number of sorted rows.
/intel/mysql/mysql_sort/scan |counter| The number of sorts that were done by scanning the table.
/intel/mysql/mysql_commands/[subnamespace] |counter| Available namespaces are evaluated in runtime, metrics indicate the number of times each statement has been executed.  The variable [subnamespace] means the command name.
/intel/mysql/mysql_handler/[subnamespace] |counter| Available namespaces are evaluated in runtime, metrics indicate the number of internal operations. The variable [subnamespace] means the operation name.

Notice, that the list of available metrics might vary depending on the MySQL version or the system configuration.