# Go + Alloy + Loki + Grafana

## Run

```bash
docker compose up --build
```

Then generate a log:

```bash
curl http://localhost:8080/hello
```

Open Grafana at [http://localhost:3000](http://localhost:3000) with `admin` / `admin`.

In **Explore**, select the `Loki` datasource and query:

```logql
{app="api"}
```

To see the same logs in Alloy's console:

```bash
docker compose logs -f alloy
```

Useful endpoints:

- API: http://localhost:8080/hello
- API health: http://localhost:8080/healthz
- Alloy UI: http://localhost:12345
- Alloy log ingestion: http://localhost:3101/loki/api/v1/push
- Loki readiness: http://localhost:3100/ready
- Grafana: http://localhost:3000

Stop the stack with:

```bash
docker compose down
```

To also remove persisted Loki and Grafana data, use `docker compose down -v`.
