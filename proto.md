# 実験

# 実験1: グラフデータベース使えそうか => Yes

module間の依存関係が与えれたとして(どう与えるかはまだ考えてない)、
graph database neo4jでやりたいこと(影響をうけるmoduleを再帰的に探す)ができるかどうかの実験

```
docker run --publish=7474:7474 --publish=7687:7687 --volume=$HOME/neo4j/data:/data neo4j
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
