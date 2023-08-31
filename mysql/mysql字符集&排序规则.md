#### 字符集
mysql支持的utf8编码最大字符长度位3字节，如果遇到4字节就会插入异常，比如：emoji表情
mysql中的utf8mb4为通常意义上的uft8，支持4个字节，可以存储emoji表情
通常使用uft8mb4即可
#### 排序规则
uft8_bin: 将字符串的每一个字符用二进制数据存储，(bin-->Binary二进制)区分大小写，可以存储二进制内容
uft8_genera_ci: 不区分大小写，仅能够在字符之间进行比较，比较速度很快，正确性较差
uft8_genera_cs: 区分大小写
utf8mb4_unicode_ci: 基于标准的unicode来排序和比较，为了处理特殊字符，实现了略微复杂的排序算法，比较速度较慢
utf8mb4_0900_ai_ci: mysql8.0以上默认的排序规则，属于utf8_unicode_ci的一种，0900是校对算法版本
utf8mb4_0900_ai_cs: 口音不敏感，大小写敏感
utf8mb4_0900_as_cs: 口音敏感，大小写敏感
ai指的是口音不敏感
as口音敏感
ci指的是不区分大小写，p和P没区别，ci是case insentsitive的缩写，即大小写不敏感，
cs指的是区分大小写，p和P有区别，cs是case sentsitive的缩写，即大小写敏感

MySQL8.0.1及更高版本将utf8mb4_0900_ai_ci作为默认排序规则，以前的版本uft8_genera_ci是默认排序规则

varchar字符集类型字段不区分大小写，p=P
SELECT * FROM `test` WHERE name="x";
|id|name|price|num|flag|
---|:--:|---:|---:|---:
|1 |    x|  21|	2|	0|
|9 |    X|    |	 |	0|