# where条件部分的map说明
一般来说：
1. key为字段名
2. value为字段值
3. value为数组时，表示多个值，会被转换为in语句
4. value为map时，表示多个条件，会被转换为and语句
5. value为null时，表示is null
6. value为string时，表示直接使用该字符串
7. value为其他类型时，表示直接使用该值
8. key不包含空格时，表示直接使用该值
9. key包含空格时，空格分割为两部分，第一部分为字段名，第二部分为操作符，如：name like

## 特殊操作符
1. $or: 表示or语句
2. $and: 表示and语句，默认为and语句，可以省略
3. $limit: 限制条数
4. $offset: 偏移条数
5. $order_by: 排序，可以是字符串或数组，如：$order_by: 'id desc' 或 $order: ['id desc', 'name asc']
6. $group_by: 分组，可以是字符串或数组，如：$group_by: 'id' 或 $group: ['id', 'name']
7. $having: 分组条件，应该是一个对象，如：$having: {"id": 1, "name !=": 'test'}
