window.BENCHMARK_DATA = {
  "lastUpdate": 1662887605357,
  "repoUrl": "https://github.com/lazyxu/kfs",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "isxuliang@gmail.com",
            "name": "Xu Liang",
            "username": "lazyxu"
          },
          "committer": {
            "email": "isxuliang@gmail.com",
            "name": "Xu Liang",
            "username": "lazyxu"
          },
          "distinct": true,
          "id": "f6707aa7b9667235636a1687cea4fa1f9be3687d",
          "message": "fix all testcases",
          "timestamp": "2022-09-11T16:55:37+08:00",
          "tree_id": "b8f5ac2822e9532115c52c429dccd4f82fdbe578",
          "url": "https://github.com/lazyxu/kfs/commit/f6707aa7b9667235636a1687cea4fa1f9be3687d"
        },
        "date": 1662887604802,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStorage0Upload1000Files1000",
            "value": 1972641604,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkStorage1Upload1000Files1000",
            "value": 1861690440,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkStorage2Upload1000Files1000",
            "value": 1733711253,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkStorage3Upload1000Files1000",
            "value": 1755970843,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkStorage4Upload1000Files1000",
            "value": 1744320178,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkCgoSqliteStorage4Upload1000Files1000",
            "value": 1483006807,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkCgoSqliteStorage4Upload10000Files1000",
            "value": 14703626873,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkMysqlStorage1Upload1000Files1000",
            "value": 2667653383,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkMysqlStorage4Upload1000Files1000",
            "value": 2443438259,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkMysqlStorage4Upload10000Files1000",
            "value": 25037034696,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkMysqlStorage5Upload1000Files1000",
            "value": 2566155962,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkMysqlStorage5Upload10000Files1000",
            "value": 24200982107,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkCgoSqliteStorage5Upload10000Files1000Batch",
            "value": 554209220,
            "unit": "ns/op",
            "extra": "2 times\n2 procs"
          },
          {
            "name": "BenchmarkGoSqliteStorage5Upload10000Files1000Batch",
            "value": 3398203357,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkMysqlStorage5Upload10000Files1000Batch",
            "value": 1044844078,
            "unit": "ns/op",
            "extra": "2 times\n2 procs"
          },
          {
            "name": "BenchmarkCgoSqliteStorage5Upload100000Files1000Batch",
            "value": 37579569188,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkGoSqliteStorage5Upload100000Files1000Batch",
            "value": 68412720354,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkMysqlStorage5Upload100000Files1000Batch",
            "value": 38484249664,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          }
        ]
      }
    ]
  }
}