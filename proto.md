# 実験

# 実験1: グラフデータベース使えそうか => Yes

module間の依存関係が与えれたとして(どう与えるかはまだ考えてない)、
graph database neo4jでやりたいこと(影響をうけるmoduleを再帰的に探す)ができるかどうかの実験

```
docker run \
  --publish=7474:7474 \
  --publish=7687:7687 \
  --volume=$PWD/neo4j/data:/data \
  --volume=$PWD/neo4j/import:/import \
  neo4j
```

```
CREATE (a:Module { name: 'User', path: 'user.rb' })
CREATE (b:Module { name: 'Order', path: 'order.rb' })
CREATE (c:Module { name: 'OrderItem', path: 'order_item.rb' })
return a, b, c

CREATE (i: Module {name: 'Item', path: 'item.rb'}), (v: Module {name: 'Variant', path: 'variant.rb'})

MATCH (a:Module),(b:Module)
WHERE a.name = 'OrderItem' AND b.name = 'Order'
CREATE (a)-[r:KNOWS]->(b)
RETURN type(r)

MATCH (a:Module),(b:Module)
WHERE a.name = 'Order' AND b.name = 'User'
CREATE (a)-[r:KNOWS]->(b)
RETURN type(r)


MATCH (a:Module),(b:Module)
WHERE a.name = 'Variant' AND b.name = 'Item'
CREATE (a)-[r:KNOWS]->(b)
RETURN type(r)

MATCH (a:Module),(b:Module)
WHERE a.name = 'OrderItem' AND b.name = 'Variant'
CREATE (a)-[r:KNOWS]->(b)
RETURN type(r)


# https://neo4j.com/docs/cypher-manual/3.5/clauses/match/#varlength-rels
# 何回までhopするのか指定できる*で無限
# つまり下記で依存の方向性ありで検索できる
MATCH (m:Module {name: 'Item'}) <-[r:KNOWS*]- (n:Module)
RETURN n

# 依存の方向性なしで検索
MATCH (m:Module {name: 'Item'}) -[r:KNOWS*]- (n:Module)
RETURN n


# たとえば、ItemとUserが変更されたときに影響を受けるnodeを検索するには..(UNIONはdefaultで重複を除いてくれる)
# https://neo4j.com/docs/cypher-manual/3.5/clauses/union/
MATCH (m:Module {name: 'Item'}) <-[r:KNOWS*]- (di:Module)
return di.name as name, di.path as path
UNION
MATCH (m:Module {name: 'User'}) <-[r:KNOWS*]- (du:Module)
return du.name as name, du.path as path
```

# 実験2: rubyからneo4jにお願いして、data import

https://neo4j.com/developer/data-import/

rubyじゃなくてもできるから、csv形式にしてあげるのが良さそう.

https://neo4j.com/docs/getting-started/current/cypher-intro/load-csv/

上記を読むかぎり2つのcsvに分けてあげるのがいいかなnodeとrelationで

node csv sample

``` csv
id,name,path
1,User,user.rb
2,Order,order.rb
3,OrderItem,order_item.rb
4,Item,item.rb
5,Variant,variant.rb
```

relation csv sample
``` csv
from,to
3,2
2,1
3,5
5,4
```

このとき、cypherで

``` cypher
# node
LOAD CSV WITH HEADERS FROM "file:///nodes.csv" AS n_row
CREATE (m:Module {id: toInteger(n_row.id), name: n_row.name, path: n_row.path})

# relation
LOAD CSV WITH HEADERS FROM "file:///relations.csv" AS r_row
MATCH (from:Module {id: toInteger(r_row.from)}), (to:Module {id: toInteger(r_row.to)})
CREATE (from) -[:KNOWS]-> (to)
```

ここまでできたので、あとはどうやってrubyでこのcsvをつくるか
