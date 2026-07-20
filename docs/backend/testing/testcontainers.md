### For integration and e2e testing we use `testcontainers` to spin up real PostgreSQL instance.
1. `Testcontainers` creates docker container for PostgreSQL from docker-compose file.
2. We create admin DB that will manage template DBs
3. We create a template from SQL migrations
4. We create a clone from template for each package that is tested (generally in main_test.go)
  - ### to speed up cloning we mount volume on tmpfs (in-memory filesystem) with fixed upper bound of allocated RAM.
5. We destroy a clone after successful/unsuccessful test
6. After all packages have finished testing `testcontainers` automatically cleans up DB's container/volume/network

#### By using `testcontainers` we are testing on infrastructure that is identical to production, without mocking the DB

This approach enables fully parallel and independent testing at the same time,
while saving resources by reusing one container and admin DB
