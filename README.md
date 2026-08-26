# 多租户全文检索索引与查询服务

纯Go实现的多租户全文检索服务，支持集合字段映射、版本化文档写入、幂等键、全文/过滤查询、相关度排序、高亮、游标分页和索引重建。内存搜索适配器用于本地运行，搜索引擎、PostgreSQL变更日志和消息队列通过接口隔离。

## 运行

需要Go 1.22或更高版本：

```bash
go test ./...
go run ./cmd/search-index-service
```

默认监听 `:8083`。可通过 `SEARCH_INDEX_HTTP_ADDR`、`SEARCH_INDEX_ENVIRONMENT` 和 `SEARCH_INDEX_SHUTDOWN_SECONDS` 覆盖配置。`configs/config.yaml`提供配置示例；SQL迁移位于`migrations`，Docker文件位于`deploy`。

## API流程

```bash
curl -X POST localhost:8083/v1/search/collections -H 'Content-Type: application/json' -d '{"id":"articles","tenant_id":"tenant-a","name":"Articles","mappings":[{"name":"title","type":"text","searchable":true}]}'
curl -X POST localhost:8083/v1/search/documents -H 'Content-Type: application/json' -H 'Idempotency-Key: write-1' -d '{"id":"doc-1","tenant_id":"tenant-a","collection_id":"articles","version":1,"fields":{"title":"Go search service","status":"published"}}'
curl -X POST localhost:8083/v1/search/query -H 'Content-Type: application/json' -d '{"tenant_id":"tenant-a","collection_id":"articles","text":"search","filters":{"status":"published"}}'
curl -X POST localhost:8083/v1/search/rebuild -H 'Content-Type: application/json' -d '{"tenant_id":"tenant-a","collection_id":"articles"}'
```

集合按租户严格隔离，文档变更通过幂等键和版本号防止重复或乱序写入。服务收到SIGINT/SIGTERM后停止消费新任务并优雅关闭。

## 目录

`cmd`为启动入口，`internal/domain`为纯领域模型和规则求值，`internal/application`为用例，`internal/adapter`为HTTP和缓存适配器，`internal/infrastructure`为配置、日志、指标和仓储实现，`api`为OpenAPI，`migrations`为数据库迁移，`deploy`和`scripts`提供部署与校验辅助。
