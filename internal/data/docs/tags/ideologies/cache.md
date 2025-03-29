---
id: cache
aliases: []
tags: []
created_at: 2025-03-28T20:03:41.000-06:00
description: Computer Caches
title: Cache
updated_at: 2025-03-28T20:20:18.000-06:00
---

## FAQ

What is the difference between write-allocate and write-no-allocate caches?

When handling write misses (i.e., when the processor attempts to write to a memory location that is not currently in the cache).

Write-allocate (also known as fetch-on-write):

When a write miss occurs, the cache controller first fetches the corresponding block from the main memory and brings it into the cache.

The write operation is then performed on the cached block.

This policy is based on the assumption that if the processor is writing to a memory location, it is likely to access that location again in the near future.

Write-allocate caches prioritize reducing future read misses by bringing the block into the cache.

Write-no-allocate (also known as write-around or write-through):

When a write miss occurs, the cache controller directly writes the data to the main memory without bringing the block into the cache.

The cache is not updated with the new data.

This policy avoids the overhead of fetching the block from memory and is beneficial when the written data is not expected to be accessed again soon.

Write-no-allocate caches prioritize reducing the penalty of write misses by avoiding the additional memory access to fetch the block.

## What is the motivation in having separate instruction and data caches?

1. Avoiding cache interference: 

Instructions and data compete for the same cache space in a unified cache. By separating them, interference between instruction and data accesses is eliminated. Instructions won't evict data and vice versa. 2. Allowing simultaneous access: With separate caches, the processor can fetch instructions and load/store data in the same cycle. This increases parallelism and performance compared to a unified cache where instructions and data accesses would be serialized. 3. Optimizing each cache independently: Instructions and data have different access patterns and locality characteristics. By using separate caches, each can be optimized based on their usage. For example: Instruction caches are typically read-only, so no write ports are needed. Data caches need to support both reads and writes. Instruction caches may use a larger block size than data caches to exploit spatial locality of code. The replacement policy can be tuned differently for each cache based on access patterns. 4. Simplifying cache coherency: With a unified cache, cache coherency is complicated by the need to keep instruction and data accesses coherent with respect to each other. Having separate caches simplifies the coherence mechanisms. 5. Security: Separating instructions and data provides a layer of protection against some security exploits that rely on the processor executing data or modifying instructions. – Why not just double the L1 cache size

Doubling the size of the L1 cache instead of having separate instruction and data caches is a valid design option, but there are several reasons why most modern processors opt for the split cache architecture:

2. Diminishing returns on cache size: 

Increasing the cache size provides diminishing performance returns due to the principle of locality. Doubling the cache size does not double the hit rate. The performance benefit from splitting the cache into instruction and data is often greater than simply doubling the size. 2. Increased latency: A larger cache typically has a longer access time due to the increased complexity of the addressing logic and the physical distance the signals need to travel. This can negate some of the benefits of the increased hit rate. Keeping each cache smaller by splitting them helps maintain low latency. 3. Power consumption: Larger caches consume more power, both in terms of leakage and dynamic power per access. Splitting the cache allows each one to be smaller and more power-efficient.4. Optimization limitations: As mentioned before, instructions and data have different characteristics. With a unified cache, any optimizations need to be a compromise between the two. Splitting the cache allows each to be optimized independently. 5. Pipeline balance: Fetching instructions and accessing data are typically separate stages in the processor pipeline. Having separate caches allows these stages to proceed in parallel without contention, providing a more balanced pipeline. 6. Cost: While the total cache size might be the same, having two separate caches does introduce some duplication in terms of the addressing logic and the peripheral circuitry. This makes the split cache design slightly more expensive in terms of chip area.
