## 通过钱包生成地址
1、公钥双哈希sha->ripe160
2、添加version（在base58Encode函数中已实现）
3、添加checksum
4、checksum校验