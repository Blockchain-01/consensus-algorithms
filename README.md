# Raft - A consensus algorithm

## How to use

### Clone the repository
```bash
git clone https://github.com/Blockchain-01/consensus-algorithms.git
cd consensus-algorithms
```   

### Run the test (for example: TestDisconnect2Followers)
Note: See the test list in the [system_test.go](./system_test.go) file.
```bash
go test -v -race -run TestDisconnect2Followers |& tee ./out/raftlog
```

### Convert the test log to a report (html) for visualization
```bash
go run ./tools/raft-testlog-viz/main.go < ./out/raftlog
```
