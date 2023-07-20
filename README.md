参考 bitcask 存储模型的建议kv数据库存储引擎。

# Bitcask

## 特性

* k-v首先存储在内存中，然后写入到磁盘;

* 在内存中可以选取哈希表、跳表、ART等数据结构;

* 读取快，数据是追加写入，但会有冗余信息，需要对磁盘文件进行合并;

* 读取首先查找内存中key的索引信息，取出value的文件位置信息(fileid, offset, value_size)，然后到磁盘中读出value;

* 不同数据类型间可以并发操作.



## 内容

* active file 活跃文件

  存于内存中，k-v写入的地方

* archived file 已封存文件

  当活跃文件的大小达到阈值后，将其封存为只读不可写，并刷盘完成持久化

* entry 条目

  由于数据库是日志类型，数据文件中存储的是每次操作的条目而非数据实体。条目包括k-v值与大小，操作数据类型，标志位，时间戳和CRC校验和等


## 问题

* 对内存数据结构的存取效率要求较高。可使用哈希表、跳表、B+Tree、ART、哈希数等。
* 空间放大，数据文件的容量可能大于数据实际量。数据合并可以一定程度缓解。


## 改进

* 批量操作Batch
* key-value类型([]byte)
* 热点数据缓存以及缓存淘汰(LRU, LFU)
* 数据压缩(snappy, zstd, zlib)
* 跨语言调用(http, grpc)


# 数据结构

## Hash

```go
type (
	Hash struct {
		record Record
	}
	Record map[string]map[string][]byte
)
```

* key -- field -- value
    * key：哈希表的名称
    * field：哈希表中的键值
    * value：哈希表中键值对应的value值

## Set

```go
type (
	Set struct {
		record Record
	}

	Record map[string]map[string]struct{}
)
```

key -- member

* key：集合的名称
* member：集合的成员
* struct{}：空结构体用于占位



## List

```go
type (
	// List idx.
	List struct {
		// record saves the List of a specified key.
		record Record

		// values saves the values of a List, help checking if a value exists in List.
		values map[string]map[string]int
	}

	// Record list record to save.
	Record map[string]*list.List
)
```

* record：存储不同列表，包含列表名key以及实体**双端链表** *list.List

* key -- value -- index，kv都是``[]byte``类型，

* values：存储每个列表的index值，用于判断元素是否存在在列表中



## ZSet

```go
type (
	// SortedSet sorted set struct
	SortedSet struct {
		record map[string]*SortedSetNode
	}

	// SortedSetNode node of sorted set
	SortedSetNode struct {
		dict map[string]*sklNode
		skl  *skipList
	}

	sklLevel struct {
		forward *sklNode
		span    uint64
	}

	sklNode struct {
		member   string
		score    float64
		backward *sklNode
		level    []*sklLevel
	}

	skipList struct {
		head   *sklNode
		tail   *sklNode
		length int64
		level  int16
	}
)
```

* SortedSet：record存放所有zset，其中key为有序集合的名称
* SortedSetNode：每个有序集合的实体
    * dict：将key与跳表节点sklNode对应的map，
    * skl：跳表数据类型
* sklNode：跳表中每个节点的信息
    * member：表示存储对象的名称或者key
    * score：有序集合中存储对象的分数
    * backward：跳表中上一个节点的地址
    * level：跳表中每层对应的下一个节点的值
        * forward：第i层下一个节点的地址
        * span：本层中前面的节点数
* sklList：跳表结构
    * head：最底层的头节点
    * tail：最底层的尾节点
    * length：数据的总个数，级最底层链表的长度
    * level：跳表目前的高度

