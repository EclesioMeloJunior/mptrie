# mptrie

The *Merkle Patricia Tree* structure is fatest to finding common prefixes and requires small memory.

![example of a merkle patricia trie](https://github.com/EclesioMeloJunior/mptrie/blob/main/assets/mptrie.png?raw=true)

### Ready to use

- [x] Put(key []byte, value []byte) error
- [x] Get(key []byte) ([]byte, bool)
- [] Storage(...)

### Test

Is possible to execute tests runing:

```sh
make test
```

The benchmarking tests are comming soon...
